package izu

import (
	"strings"

	"github.com/iancoleman/strcase"
	lua "github.com/yuin/gopher-lua"
)

func registerHelpers(state *lua.LState) {
	state.SetGlobal("lowercase", state.NewFunction(lowercase))
	state.SetGlobal("uppercase", state.NewFunction(uppercase))
	state.SetGlobal("pascalcase", state.NewFunction(pascalcase))
}

func lowercase(state *lua.LState) int {
	str := state.CheckString(1)
	state.Push(lua.LString(strings.ToLower(str)))
	return 1
}

func uppercase(state *lua.LState) int {
	str := state.CheckString(1)
	state.Push(lua.LString(strings.ToUpper(str)))
	return 1
}

func pascalcase(state *lua.LState) int {
	str := state.CheckString(1)
	state.Push(lua.LString(strcase.ToCamel(str)))
	return 1
}
