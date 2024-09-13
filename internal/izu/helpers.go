package izu

import (
	"strings"

	"github.com/iancoleman/strcase"
	"github.com/meir/izu/pkg/izu"
	lua "github.com/yuin/gopher-lua"
)

// registerHelpers will register the helper functions to the lua
func registerHelpers(state *lua.LState) {
	state.SetGlobal("lowercase", state.NewFunction(lowercase))
	state.SetGlobal("uppercase", state.NewFunction(uppercase))
	state.SetGlobal("pascalcase", state.NewFunction(pascalcase))
	state.SetGlobal("has_key", state.NewFunction(hasKey))
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

// pascalcase will convert a string to PascalCase
func pascalcase(state *lua.LState) int {
	str := state.CheckString(1)
	state.Push(lua.LString(strcase.ToCamel(str)))
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
func addKey(state *lua.LState) int {
	value := state.Get(1)

	switch value.Type() {
	case lua.LTString:
		key := value.String()
		izu.AddValidationKey(key)
	case lua.LTTable:
		value.(*lua.LTable).ForEach(func(k, v lua.LValue) {
			key := v.String()
			izu.AddValidationKey(key)
		})
	}
	return 0
}
