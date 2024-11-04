package luaformatter

import (
	"fmt"
	"log/slog"

	"github.com/meir/izu/pkg/izu"
	lua "github.com/yuin/gopher-lua"
)

// Formatter is a lua based formatter to generate hotkey configurations
type Formatter struct {
	system  string
	state   *lua.LState
	methods map[string]lua.LValue
}

// NewFormatter creates a new lua formatter for the given system
func NewFormatter(system string) (*Formatter, error) {
	// initialize helper methods
	slog.Debug("Initializing lua formatter", "system", system)
	state := lua.NewState()
	table := state.NewTable()
	table.RawSetString("lowercase", state.NewFunction(lowercase))
	table.RawSetString("uppercase", state.NewFunction(uppercase))
	table.RawSetString("hasKey", state.NewFunction(hasKey))

	state.SetGlobal("izu", table)

	if system == "" {
		return nil, fmt.Errorf("formatter/system cannot be empty")
	}

	// load in lua formatter file
	slog.Debug("Loading lua formatter file", "system", system)
	content, err := izu.GetFormatterFile("lua", system)
	if err != nil {
		return nil, fmt.Errorf("failed to load lua formatter file for %s: %w", system, err)
	}

	// run the lua file in order to retrieve the AST methods
	slog.Debug("Running lua formatter file", "system", system)
	if err := state.DoString(string(content)); err != nil {
		return nil, err
	}

	slog.Debug("Checking for the module returned by the lua formatter file")
	module := state.Get(-1)

	// check if the response is an object
	methods := map[string]lua.LValue{}
	if module, ok := module.(*lua.LTable); ok {
		// add all the AST names in a list, these will be used as the required method names
		asts := []string{
			izu.ASTHotkey.String(),
			izu.ASTBinding.String(),
			izu.ASTSingle.String(),
			izu.ASTMultiple.String(),
			izu.ASTString.String(),
		}

		// go through each AST name and see if the returned object has a method with the same name
		// if not, return an error
		for _, method := range asts {
			if function := module.RawGetString(method); function.Type() == lua.LTFunction {
				slog.Debug("Found method in lua formatter module", "method", method)
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

// Call will run the lua method for the given AST type using the options given
func (formatter *Formatter) Call(method izu.AST, options ...Option) ([]string, error) {
	// get the method based on the AST type
	function, ok := formatter.methods[method.String()]
	if !ok {
		return nil, fmt.Errorf("cannot find method for %s", method.String())
	}

	// add all the options in a table using the key-value
	// later ones will override the key, this just means that its higher in the tree
	value := &lua.LTable{}
	for _, option := range options {
		value.RawSetString(option.name, option.value)
	}

	// call the lua method
	err := formatter.state.CallByParam(lua.P{
		Fn:      function,
		NRet:    1,
		Protect: true,
	}, value)
	if err != nil {
		return nil, fmt.Errorf("failed to call lua formatting method %s: %w", method.String(), err)
	}

	response := formatter.state.Get(-1)

	// check if the response is either a string or a string array
	// if its anything else, return an error
	switch response.Type() {
	case lua.LTString:
		return []string{response.String()}, nil
	case lua.LTTable:
		output := []string{}
		var err error
		// loop through all the values in the table to check if any of them are not a string
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

// Format will take a list of hotkeys and format them into strings that can be used in the config file of the hotkey system
func (formatter *Formatter) Format(hotkeys []*izu.Hotkey) ([]string, error) {
	slog.Debug("Formatting hotkeys", "system", formatter.system)
	output := []string{}
	for _, hotkey := range hotkeys {
		slog.Debug("Formatting hotkey", "hotkey", hotkey.String())
		flags := []string{}
		// check if there are any flags assigned for this system
		if sflags, ok := hotkey.Flags[formatter.system]; ok {
			flags = sflags
		}

		// format the binding of this hotkey
		bindings, err := formatter.format(hotkey.Binding, OptionFlags(flags), OptionStateBinding())
		if err != nil {
			return nil, err
		}

		// check if theres a specific command for this system, otherwise use the default
		// if theres no default and this system is not specified, return an error
		var command izu.Part
		if scommand, ok := hotkey.Command[formatter.system]; ok {
			command = scommand
		} else if scommand, ok := hotkey.Command["default"]; ok {
			command = scommand
		} else {
			return nil, fmt.Errorf("no command found for hotkey '%s' on system '%s'", hotkey.String(), formatter.system)
		}

		// format the command part of this hotkey
		commands, err := formatter.format(command, OptionFlags(flags), OptionStateCommand())

		for i, binding := range bindings {
			// each hotkey might turn into several bindings and several commands (due to multiples)
			// if this is the case, we need to find the command thats part of the current binding
			command := commands[i%len(commands)]

			// format the final hotkey
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

			// and add it to the output
			output = append(output, response...)
			slog.Debug("Formatted hotkey", "binding", binding, "command", command)
		}
	}
	return output, nil
}

// format will take a part and format it into one or multiple bindings/commands and call the lua methods in order to properly format it
func (formatter *Formatter) format(root izu.Part, opts ...Option) (output []string, err error) {
	kind, partlist := root.Info()

	// if the part is a string, call the lua method and return its output, we dont need any other processing on this part
	if kind == izu.ASTString {
		opts = append(opts, OptionString(root.String()))
		opts = append(opts, OptionAST(kind))
		output, err = formatter.Call(izu.ASTString, opts...)
		return
	}

	inputs := [][]string{{}}
	// iterate through all the subparts in this part
	err = partlist.Iterate(func(part izu.Part) error {
		kind, _ := part.Info()

		// format the subpart recursively
		opts = append(opts, OptionAST(kind))
		bindings, err := formatter.format(part, opts...)
		if err != nil {
			return err
		}

		// multiple the inputs by the bindings given from the subpart
		// this assures we get a full binding/command for every multiple
		ninputs := make([][]string, len(inputs)*len(bindings))
		for x, binding := range bindings {
			for y, input := range inputs {
				// create a new slice like this
				// apparently if you just use append(input, binding) and assign that
				// it might sometimes point to the original slice and put the last binding in each input
				entry := append([]string{}, input...)
				entry = append(entry, binding)
				ninputs[(x*len(inputs))+y] = entry
			}
		}

		inputs = ninputs
		if kind == izu.ASTSingle {
			inputs = izu.CapitalizeKey(ninputs)
		}

		return nil
	})
	if err != nil {
		return
	}

	// call the lua method using the inputs and return the output
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
