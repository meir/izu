package izu

import (
	"fmt"
	"io"
	"os"
	"path"
	"strings"

	"github.com/meir/izu/internal/parser"
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
			izu.ASTBinding.String(),
			izu.ASTSingle.String(),
			izu.ASTMultiple.String(),
			izu.ASTString.String(),
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
func (luaf *LuaFormatter) call(method izu.AST, value lua.LValue, section int8) ([]string, error) {
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
func (luaf *LuaFormatter) recursive_format(hotkeys []*izu.Hotkey, section int8) ([]string, error) {
	inputs := [][]string{}
	for _, hotkey := range hotkeys {
		// get the state and parts of the part
		state, parts := hotkey.Info()

		switch state {
		case izu.ASTBinding, izu.ASTSingle, izu.ASTMultiple:
			// initialize the first line
			inputs = append(inputs, []string{})

			// loop through all parts and format them
			err := parts.Iterate(func(part izu.Part) error {
				strs, err := luaf.recursive_format(part, section)
				if err != nil {
					return err
				}

				// mesh the lines and the formatted strings together like a matrix
				newInputs := make([][]string, len(inputs)*len(strs))
				for i, str := range strs {
					for j, input := range inputs {
						newInputs[i*len(inputs)+j] = append(input, str)
					}
				}
				inputs = newInputs
				return nil
			})
			if err != nil {
				return nil, err
			}
			break

		case izu.ASTString:
			return luaf.call(state, lua.LString(part.String()), section)
		}
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

// getSingles is used to iterate through the partlist tree and find any single parts
// sadly you cant find single people using this function
func getSingles(partlist izu.PartList) []izu.Part {
	singles := []izu.Part{}
	partlist.Iterate(func(part izu.Part) error {
		if kind, subpartlist := part.Info(); kind == izu.ASTSingle {
			singles = append(singles, part)
		} else {
			singles = append(singles, getSingles(subpartlist)...)
		}
		return nil
	})
	return singles
}

// validateKeys is used to validate the keys in the binding
func validateKeys(binding izu.PartList) error {
	// go through all single parts
	for _, part := range getSingles(binding) {
		bindings := [][]izu.Part{{}}
		_, parts := part.Info()

		// add the parts accordingly to the array as either string or get the string parts from the singlepart
		parts.Iterate(func(sub_part izu.Part) error {
			state, subparts := sub_part.Info()

			switch state {
			case izu.ASTString:
				// add the string part to all the bindings
				for bind := range bindings {
					bindings[bind] = append(bindings[bind], sub_part)
				}
			case izu.ASTMultiple:
				// multiply the bindings by the amount of string parts in multiple
				newBindings := [][]izu.Part{}
				for _, bind := range bindings {
					subparts.Iterate(func(subbinding izu.Part) error {
						newBindings = append(newBindings, append(bind, subbinding))
						return nil
					})
				}
				bindings = newBindings
			}
			return nil
		})

		// go through all binding arrays
		for _, bind := range bindings {
			binding := ""
			// build up the actual keycode
			for _, part := range bind {
				binding += part.String()
			}

			// check if the keycode is actually valid and get the actual casing
			str, ok := izu.Validate(binding)
			if !ok {
				return fmt.Errorf("invalid keybind: %s", binding)
			}

			// apply the casing to the individual parts
			j := 0
			for i := 0; i < len(bind); i++ {
				if kind, _ := bind[i].Info(); kind == izu.ASTString {
					// this will break multi parts for custom keys but i have no idea how to map those
					// would have to have some magic to understand how to map stuff like `XF86Audio{Play,Pause}` to a custom key like `MediaStart`
					// realistically this wouldnt happen but for sway we need to map `Super` to `Mod4`, so you cant break it into `S{_,uper} + w` for example, because it would break in sway
					end := j + len(bind[i].String())
					if end > len(str) {
						end = len(str)
					}
					bind[i].Append(parser.NewPartString(str[j:end]))
					j += len(bind[i].String())
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
	hotkeys, err := parser.Parse(s)
	if err != nil {
		return nil, err
	}

	// validate the keybind
	for _, hotkey := range hotkeys {
		if err := validateKeys(hotkey.Binding); err != nil {
			return nil, err
		}
	}

	// format the keybind using the lua formatter
	str, err := luaf.recursive_format(hotkeys, -1)
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

	return luaf.ParseString(content)
}
