package luaformatter

import (
	"strings"

	lua "github.com/yuin/gopher-lua"
)

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

// contains will check if a key exists in a table
func contains(state *lua.LState) int {
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
