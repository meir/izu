package izu

import (
	"strings"

	"github.com/meir/izu/pkg/izu"
	lua "github.com/yuin/gopher-lua"
)

// registerHelpers will register the helper functions to the lua
func registerHelpers(state *lua.LState) {
	table := state.NewTable()
	table.RawSetString("lowercase", state.NewFunction(lowercase))
	table.RawSetString("uppercase", state.NewFunction(uppercase))
	table.RawSetString("hasKey", state.NewFunction(hasKey))
	table.RawSetString("registerKeycode", state.NewFunction(registerKeycode))

	state.SetGlobal("izu", table)
}

// lowercase will convert a string to lowercase
func lowercase(state *lua.LState) int {
	str := state.CheckString(1)
	state.Push(lua.LString(strings.ToLower(str)))
	return 1
}

// uppercase will convert a string to uppercase
func uppercase(state *lua.LState) int {
	str := state.CheckString(1)
	state.Push(lua.LString(strings.ToUpper(str)))
	return 1
}

// hasKey will check if a key exists in a table
func hasKey(state *lua.LState) int {
	table := state.CheckTable(1)
	key := state.CheckString(2)
	result := false
	table.ForEach(func(k, v lua.LValue) {
		if v.String() == key {
			result = true
		}
	})
	state.Push(lua.LBool(result))
	return 1
}

// addKey will add a key to the validation map
func registerKeycode(state *lua.LState) int {
	keyOrValue := state.Get(1)
	value := state.Get(2)

	switch keyOrValue.Type() {
	case lua.LTString:
		key := keyOrValue.String()
		value := value.String()

		// if there has been no value given, we use the key as the value
		if value == "" {
			value = key
			key = strings.ToLower(key)
		}

		izu.AddValidationKey(key, value)
	case lua.LTTable:
		keyOrValue.(*lua.LTable).ForEach(func(k, v lua.LValue) {
			key := k.String()
			value := v.String()

			// if the key is a number, we use the value as the key
			if strings.Trim(key, "0123456789") == "" {
				key = strings.ToLower(value)
			}

			izu.AddValidationKey(key, value)
		})
	}
	return 0
}
