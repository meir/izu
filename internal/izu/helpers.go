package izu

import (
	"strings"

	"github.com/iancoleman/strcase"
	lua "github.com/yuin/gopher-lua"
)

// registerHelpers will register the helper functions to the lua
func registerHelpers(state *lua.LState) {
	state.SetGlobal("lowercase", state.NewFunction(lowercase))
	state.SetGlobal("uppercase", state.NewFunction(uppercase))
	state.SetGlobal("pascalcase", state.NewFunction(pascalcase))
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
