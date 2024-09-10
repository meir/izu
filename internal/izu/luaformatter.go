package izu

import (
	"fmt"
	"io"
	"os"
	"path"
	"strings"

	"github.com/meir/izu/internal/izu/parser"
	"github.com/meir/izu/pkg/izu"
	lua "github.com/yuin/gopher-lua"
)

type LuaFormatter struct {
	name  string
	state *lua.LState

	methods map[string]lua.LValue
}

func NewLuaFormatter(name string) (*LuaFormatter, error) {
	luaf := &LuaFormatter{
		name:  name,
		state: lua.NewState(),
	}

	registerHelpers(luaf.state)

	if err := luaf.load(); err != nil {
		return nil, err
	}
	return luaf, nil
}

func (luaf *LuaFormatter) load() error {
	data, err := luaf.getFormatter()
	if err != nil {
		return err
	}

	// use lua to run the formatter code and receive the return object from it
	if err := luaf.state.DoString(string(data)); err != nil {
		return fmt.Errorf("failed to load formatter: %w", err)
	}

	// get the return object from the lua code
	ret := luaf.state.Get(-1)

	// expect a table with the following functions:
	// base (parts)
	// multiple (parts)
	// single (parts)
	// single_part (parts)
	// string (parts)
	luaf.methods = make(map[string]lua.LValue)
	if tbl, ok := ret.(*lua.LTable); ok {
		// get the base function
		methods := []string{
			"keybind",
			"command",
			"base",
			"multiple",
			"single",
			"single_part",
			"string",
		}

		for _, method := range methods {
			if fn := tbl.RawGetString(method); fn.Type() == lua.LTFunction {
				luaf.methods[method] = fn
			} else {
				return fmt.Errorf("expected a function '%s' to be returned", method)
			}
		}
	}

	return nil
}

func (luaf *LuaFormatter) getFormatter() (data []byte, err error) {
	if !strings.HasSuffix(luaf.name, ".lua") {
		luaf.name += ".lua"
	}

	if file, err := izu.Formatters.Open(path.Join("formatters", luaf.name)); err == nil {
		// load file from embedded resources
		data, err = io.ReadAll(file)
	} else if os.IsNotExist(err) {
		// load file from disk
		if _, err = os.Stat(luaf.name); err != nil {
			return nil, os.ErrNotExist
		}

		data, err = os.ReadFile(luaf.name)
	}

	return
}

// methods

func (luaf *LuaFormatter) call_method(method lua.LValue, methodName string, value lua.LValue) ([]string, error) {
	err := luaf.state.CallByParam(lua.P{
		Fn:      method,
		NRet:    1,
		Protect: true,
	}, value)
	if err != nil {
		panic(err)
	}

	ret := luaf.state.Get(-1)

	switch ret.Type() {
	case lua.LTString:
		return []string{ret.String()}, nil
	case lua.LTTable:
		out := []string{}
		ret.(*lua.LTable).ForEach(func(key, value lua.LValue) {
			out = append(out, value.String())
		})

		return out, nil
	default:
		return nil, fmt.Errorf("expected a string or string slice to be returned from formatter for method " + methodName)
	}
}

func (luaf *LuaFormatter) call(state izu.State, input []string) ([]string, error) {
	var method lua.LValue
	var methodName string
	switch state {
	case izu.StateKeybind:
		method = luaf.methods["keybind"]
		methodName = "keybind"
	case izu.StateCommand:
		method = luaf.methods["command"]
		methodName = "command"
	case izu.StateBase:
		method = luaf.methods["base"]
		methodName = "base"
	case izu.StateMultiple:
		method = luaf.methods["multiple"]
		methodName = "multiple"
	case izu.StateSingle:
		method = luaf.methods["single"]
		methodName = "single"
	case izu.StateSinglePart:
		method = luaf.methods["single_part"]
		methodName = "single_part"
	case izu.StateString:
		method = luaf.methods["string"]
		methodName = "string"
	}

	table := luaf.state.NewTable()
	for i, str := range input {
		table.RawSetInt(i+1, lua.LString(str))
	}

	return luaf.call_method(method, methodName, table)
}

// formatter calls

func (luaf *LuaFormatter) recursive_format(part izu.Part) ([]string, error) {
	state, parts := part.Info()
	switch state {
	case izu.StateKeybind:
		bind := parts[0]
		command := parts[1]

		options, err := luaf.recursive_format(bind)
		if err != nil {
			return nil, err
		}

		commands, err := luaf.recursive_format(command)
		if err != nil {
			return nil, err
		}

		inputs := [][]string{}
		for i := range options {
			inputs = append(inputs, []string{options[i], commands[i%len(commands)]})
		}

		outputs := []string{}
		for _, input := range inputs {
			strs, err := luaf.call(state, input)
			if err != nil {
				return nil, err
			}
			outputs = append(outputs, strs...)
		}
		return outputs, nil
	case izu.StateCommand, izu.StateBase, izu.StateMultiple, izu.StateSingle, izu.StateSinglePart:
		{
			inputs := [][]string{{}}
			for _, part := range parts {
				strs, err := luaf.recursive_format(part)
				if err != nil {
					return nil, err
				}

				switch len(strs) {
				case 0:
					continue
				case 1:
					for i := range inputs {
						inputs[i] = append(inputs[i], strs[0])
					}
				default:
					newInputs := make([][]string, len(inputs)*len(strs))
					for i, str := range strs {
						for j, input := range inputs {
							newInputs[i*len(inputs)+j] = append(input, str)
						}
					}
					inputs = newInputs
				}
			}

			outputs := []string{}
			for _, input := range inputs {
				strs, err := luaf.call(state, input)
				if err != nil {
					return nil, err
				}
				outputs = append(outputs, strs...)
			}
			return outputs, nil
		}

	case izu.StateString:
		{
			if str, ok := part.(*parser.String); ok {
				return luaf.call(state, []string{str.Key()})
			}
		}
	}

	return []string{}, nil
}

// Public methods

func (luaf *LuaFormatter) ParseString(s string) ([]string, error) {
	keybind := parser.NewKeybind(luaf)
	if _, err := keybind.Parse([]byte(s)); err != nil {
		return nil, err
	}
	str, err := luaf.recursive_format(keybind)
	return str, err
}
