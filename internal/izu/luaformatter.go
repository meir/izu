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

	// expect a table to be returned with the required methods
	luaf.methods = make(map[string]lua.LValue)
	if tbl, ok := ret.(*lua.LTable); ok {
		methods := []string{
			izu.StateKeybind.String(),
			izu.StateCommand.String(),
			izu.StateBase.String(),
			izu.StateMultiple.String(),
			izu.StateSingle.String(),
			izu.StateSinglePart.String(),
			izu.StateString.String(),
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

// call is used to call a method in the lua formatter
func (luaf *LuaFormatter) call(method izu.State, value lua.LValue) ([]string, error) {
	var luamethod lua.LValue
	var ok bool
	if luamethod, ok = luaf.methods[method.String()]; !ok {
		return nil, fmt.Errorf("cannot find method for %s", method.String())
	}

	// call the method
	err := luaf.state.CallByParam(lua.P{
		Fn:      luamethod,
		NRet:    1,
		Protect: true,
	}, value)
	if err != nil {
		return nil, fmt.Errorf("failed to call lua formatting method %s: %w", method.String(), err)
	}

	// get the response
	response := luaf.state.Get(-1)

	// check the type of the response and return the string array
	switch response.Type() {
	case lua.LTString:
		return []string{response.String()}, nil
	case lua.LTTable:
		out := []string{}
		response.(*lua.LTable).ForEach(func(key, value lua.LValue) {
			if value.Type() != lua.LTString {
				panic("expected a string to be returned from formatter")
			}
			out = append(out, value.String())
		})

		return out, nil

	default:
		// if method doesnt return a string or string array, error
		return nil, fmt.Errorf("expected a string or string slice to be returned from formatter for method %s", method.String())
	}
}

// recursive_format is used to recursively format the keybind by going through all the parts and formatting them using lua methods
func (luaf *LuaFormatter) recursive_format(part izu.Part) ([]string, error) {
	// get the state and parts of the part
	state, parts := part.Info()
	inputs := [][]string{}

	switch state {
	case izu.StateKeybind:
		// prepare the binding part
		binds, err := luaf.recursive_format(parts[0])
		if err != nil {
			return nil, err
		}

		// prepare the command part
		commands, err := luaf.recursive_format(parts[1])
		if err != nil {
			return nil, err
		}

		// mesh the commands onto the bindings
		for i := range binds {
			inputs = append(inputs, []string{binds[i], commands[i%len(commands)]})
		}
		break

	case izu.StateCommand, izu.StateBase, izu.StateMultiple, izu.StateSingle, izu.StateSinglePart:
		// initialize the first line
		inputs = append(inputs, []string{})

		// loop through all parts and format them
		for _, part := range parts {
			strs, err := luaf.recursive_format(part)
			if err != nil {
				return nil, err
			}

			// mesh the lines and the formatted strings together like a matrix
			newInputs := make([][]string, len(inputs)*len(strs))
			for i, str := range strs {
				for j, input := range inputs {
					newInputs[i*len(inputs)+j] = append(input, str)
				}
			}
			inputs = newInputs
		}
		break

	case izu.StateString:
		// strings should be handled differently, they should directly be called as
		//this is the lowest AST type there is, only thing it can do is transform text
		if str, ok := part.(*parser.String); ok {
			return luaf.call(state, lua.LString(str.Key()))
		}
		return nil, fmt.Errorf("formatter part is returning the wrong state")
	}

	// call the lua method with the inputs
	outputs := []string{}
	for _, input := range inputs {
		// prepare input table
		table := luaf.state.NewTable()
		for i, str := range input {
			table.RawSetInt(i+1, lua.LString(str))
		}

		// call the lua method
		strs, err := luaf.call(state, table)
		if err != nil {
			return nil, err
		}
		outputs = append(outputs, strs...)
	}
	return outputs, nil
}

// ParseString is used to parse a string using the lua formatter
// this expects 2 lines, in the format of "[keybind]\n[command]"
// formatting using spaces is allowed but will be omitted in the final format
// example:
// XF86Audio{Play,Pause}
// playerctl --{play,pause}
func (luaf *LuaFormatter) ParseString(s []byte) ([]byte, error) {
	// create a new keybind parser
	keybind := parser.NewKeybind()
	if _, err := keybind.Parse(s); err != nil {
		return nil, err
	}

	// format the keybind using the lua formatter
	str, err := luaf.recursive_format(keybind)
	return []byte(strings.Join(str, "\n")), err
}

// ParseFile is used to parse a file using the lua formatter
// this expects the file to be in the same format as ParseString
// example:
// XF86Audio{Play,Pause}
// playerctl --{play,pause}
func (luaf *LuaFormatter) ParseFile(f string) ([]byte, error) {
	content, err := os.ReadFile(f)
	if err != nil {
		return nil, err
	}

	hotkeys := []string{}
	lines := strings.Split(string(content), "\n")
	// iterate through the config lines to bundle the binding and the command together
	for i := 0; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])
		// skip if the line is empty or commented
		if line == "" || line[0] == '#' {
			continue
		}

		// binding part
		bind := line

		// iterate through the config lines to find the command thats part of the binding
		// if the config file ends or has only comments left, it will return an error
		j := 1
		for {
			// check if file hasnt ended yet
			if len(lines) < i+j {
				return nil, fmt.Errorf("expected a command after keybind")
			}

			// get the command part and check if its not empty or a comment
			command := strings.TrimSpace(lines[i+j])
			if command == "" {
				return nil, fmt.Errorf("expected a command after keybind")
			} else if command[0] == '#' {
				// if the line is a comment, continue to the next line
				j++
				continue
			}

			// combine the binding and the command
			hotkeys = append(hotkeys, bind+"\n"+command)
			i += j
			break
		}
	}

	// process all the hotkeys found
	output := []string{}
	for _, hotkey := range hotkeys {
		binding, err := luaf.ParseString([]byte(hotkey))
		if err != nil {
			return nil, err
		}

		output = append(output, string(binding))
	}

	return []byte(strings.Join(output, "\n\n")), nil
}
