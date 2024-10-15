package luaformatter

import (
	"fmt"

	"github.com/meir/izu/pkg/izu"
	lua "github.com/yuin/gopher-lua"
)

type Formatter struct {
	system  string
	state   *lua.LState
	methods map[string]lua.LValue
}

func NewFormatter(system string) (*Formatter, error) {
	// initialize helper methods
	state := lua.NewState()
	table := state.NewTable()
	table.RawSetString("lowercase", state.NewFunction(lowercase))
	table.RawSetString("uppercase", state.NewFunction(uppercase))
	table.RawSetString("hasKey", state.NewFunction(hasKey))
	table.RawSetString("registerKeycode", state.NewFunction(registerKeycode))

	state.SetGlobal("izu", table)

	// load in lua formatter file
	content, err := izu.GetFormatterFile("lua", system)
	if err != nil {
		return nil, err
	}

	if err := state.DoString(string(content)); err != nil {
		return nil, err
	}

	module := state.Get(-1)

	methods := map[string]lua.LValue{}
	if module, ok := module.(*lua.LTable); ok {
		asts := []string{
			izu.ASTHotkey.String(),
			izu.ASTBinding.String(),
			izu.ASTSingle.String(),
			izu.ASTMultiple.String(),
			izu.ASTString.String(),
		}

		for _, method := range asts {
			if function := module.RawGetString(method); function.Type() == lua.LTFunction {
				methods[method] = function
			} else {
				return nil, fmt.Errorf("expected a function '%s' to be returned within the lua formatter module", method)
			}
		}
	} else {
		return nil, fmt.Errorf("expected a table to be returned in the lua formatter file")
	}

	return &Formatter{
		system:  system,
		state:   state,
		methods: methods,
	}, nil
}

func (formatter *Formatter) Call(method izu.AST, options ...Option) ([]string, error) {
	function, ok := formatter.methods[method.String()]
	if !ok {
		return nil, fmt.Errorf("cannot find method for %s", method.String())
	}

	value := &lua.LTable{}
	for _, option := range options {
		value.RawSetString(option.name, option.value)
	}

	err := formatter.state.CallByParam(lua.P{
		Fn:      function,
		NRet:    1,
		Protect: true,
	}, value)
	if err != nil {
		return nil, fmt.Errorf("failed to call lua formatting method %s: %w", method.String(), err)
	}

	response := formatter.state.Get(-1)

	switch response.Type() {
	case lua.LTString:
		return []string{response.String()}, nil
	case lua.LTTable:
		output := []string{}
		var err error
		response.(*lua.LTable).ForEach(func(key, value lua.LValue) {
			if err != nil {
				// skip to the end if theres been an error
				return
			}

			if value.Type() != lua.LTString {
				err = fmt.Errorf("expected only strings to be returned in an array from formatter method '%s', received a '%s'", method.String(), value.Type().String())
			}

			output = append(output, value.String())
		})

		return output, nil
	default:
		return nil, fmt.Errorf("expected a string or string array to be returned from formatter method '%s' instead got '%s'", method.String(), response.Type().String())
	}
}

func (formatter *Formatter) Format(hotkeys []*izu.Hotkey) ([]string, error) {
	output := []string{}
	for _, hotkey := range hotkeys {
		flags := []string{}
		if sflags, ok := hotkey.Flags[formatter.system]; ok {
			flags = sflags
		}

		bindings, err := formatter.format(hotkey.Binding, OptionFlags(flags), OptionStateBinding())
		if err != nil {
			return nil, err
		}

		var command izu.Part
		if scommand, ok := hotkey.Command[formatter.system]; ok {
			command = scommand
		} else if scommand, ok := hotkey.Command["default"]; ok {
			command = scommand
		} else {
			return nil, fmt.Errorf("no command found for hotkey '%s' on system '%s'", hotkey.String(), formatter.system)
		}

		commands, err := formatter.format(command, OptionFlags(flags), OptionStateCommand())

		for i, binding := range bindings {
			command := commands[i%len(commands)]
			response, err := formatter.Call(izu.ASTHotkey,
				OptionStringArray([]string{
					binding,
					command,
				}),
				OptionStateHotkey(),
				OptionAST(izu.ASTHotkey),
				OptionFlags(flags),
			)
			if err != nil {
				return nil, err
			}

			output = append(output, response...)
		}

	}
	return output, nil
}

func (formatter *Formatter) format(root izu.Part, opts ...Option) (output []string, err error) {
	kind, partlist := root.Info()

	if kind == izu.ASTString {
		opts = append(opts, OptionString(root.String()))
		opts = append(opts, OptionAST(kind))
		output, err = formatter.Call(izu.ASTString, opts...)
		return
	}

	inputs := [][]string{{}}
	err = partlist.Iterate(func(part izu.Part) error {
		kind, _ := part.Info()

		opts = append(opts, OptionAST(kind))
		bindings, err := formatter.format(part, opts...)
		if err != nil {
			return err
		}

		ninputs := make([][]string, len(inputs)*len(bindings))
		for x, binding := range bindings {
			for y, input := range inputs {
				entry := append([]string{}, input...)
				entry = append(entry, binding)
				ninputs[(x*len(inputs))+y] = entry
			}
		}

		inputs = ninputs

		return nil
	})
	if err != nil {
		return
	}

	for _, input := range inputs {
		opts = append(opts, OptionStringArray(input))
		response, err := formatter.Call(kind, opts...)
		if err != nil {
			return nil, err
		}
		output = append(output, response...)
	}

	return
}
