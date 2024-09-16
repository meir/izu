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
func (luaf *LuaFormatter) call(method izu.State, value lua.LValue, section int8) ([]string, error) {
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
	}, value, lua.LNumber(section))
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
func (luaf *LuaFormatter) recursive_format(part izu.Part, section int8) ([]string, error) {
	// get the state and parts of the part
	state, parts := part.Info()
	inputs := [][]string{}

	switch state {
	case izu.StateKeybind:
		// prepare the binding part
		binds, err := luaf.recursive_format(parts[0], 0)
		if err != nil {
			return nil, err
		}

		// prepare the command part
		commands, err := luaf.recursive_format(parts[1], 1)
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
			strs, err := luaf.recursive_format(part, section)
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
			return luaf.call(state, lua.LString(str.Key), section)
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
		strs, err := luaf.call(state, table, section)
		if err != nil {
			return nil, err
		}
		outputs = append(outputs, strs...)
	}
	return outputs, nil
}

// getSingles is used to harvest all the single parts from the children of parts
// sadly you cant find single people using this function
func getSingles(binding izu.Part) []izu.Part {
	state, parts := binding.Info()
	// if this is a single state we assume it doesnt have any single children,
	// otherwise the parser did something wrong
	if state == izu.StateSingle {
		return []izu.Part{binding}
	}

	// loop through all the parts and get the singles from them
	out := []izu.Part{}
	for _, part := range parts {
		out = append(out, getSingles(part)...)
	}
	return out
}

// validateKeys is used to validate the keys in the binding
func validateKeys(binding izu.Part) error {
	// go through all single parts
	for _, part := range getSingles(binding) {
		bindings := [][]izu.Part{{}}
		_, parts := part.Info()

		// add the parts accordingly to the array as either string or get the string parts from the singlepart
		for _, sub_part := range parts {
			state, subparts := sub_part.Info()

			switch state {
			case izu.StateString:
				// add the string part to all the bindings
				for bind := range bindings {
					bindings[bind] = append(bindings[bind], sub_part)
				}
			case izu.StateSinglePart:
				// multiply the bindings by the amount of string parts in singleparts
				newBindings := [][]izu.Part{}
				for _, bind := range bindings {
					for _, str_part := range subparts {
						newBindings = append(newBindings, append(bind, str_part))
					}
				}
				bindings = newBindings
			}
		}

		// go through all binding arrays
		for _, bind := range bindings {
			binding := ""
			// build up the actual keycode
			for _, part := range bind {
				binding += part.(*parser.String).Key
			}

			// check if the keycode is actually valid and get the actual casing
			str, ok := izu.Validate(binding)
			if !ok {
				return fmt.Errorf("invalid keybind: %s", binding)
			}

			// apply the casing to the individual parts
			j := 0
			for i := 0; i < len(bind); i++ {
				if strbind, ok := bind[i].(*parser.String); ok {
					// this will break multi parts for custom keys but i have no idea how to map those
					// would have to have some magic to understand how to map stuff like `XF86Audio{Play,Pause}` to a custom key like `MediaStart`
					// realistically this wouldnt happen but for sway we need to map `Super` to `Mod4`, so you cant break it into `S{_,uper} + w` for example, because it would break in sway
					end := j + len(strbind.Key)
					if end > len(str) {
						end = len(str)
					}
					strbind.Key = str[j:end]
					j += len(strbind.Key)
				}
			}
		}
	}

	return nil
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
	//TODO: the parser can actually be changed into a list instead of a tree, that would make stuff
	//  so much easier and the parser would be incredibly small compared to what we have now
	//  Also, we wont have to have a process of parsing single lines for izu.ParseFile anymore if we do
	//  it would also mean we wont be able to do hacky stuff like "XF68{Audio{Play,Pause},MediaStop}" anymore but that would only be better,
	//  im not even sure how to make a command for this even.
	if _, err := keybind.Parse(s); err != nil {
		return nil, err
	}

	// validate the keybind
	if err := validateKeys(keybind); err != nil {
		return nil, err
	}

	// format the keybind using the lua formatter
	str, err := luaf.recursive_format(keybind, -1)
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
	// split the binding from the command by newline or semicolon
	lines := strings.FieldsFunc(string(content), func(r rune) bool {
		return r == '\n' || r == ';'
	})
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
			if command == "" || command[0] == '#' {
				// if the line is empty or a comment, continue to the next line
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
