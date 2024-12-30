package luaformatter

import (
	"github.com/meir/izu/pkg/izu"
	lua "github.com/yuin/gopher-lua"
)

type Option struct {
	name  string
	value lua.LValue
}

func OptionString(value string) Option {
	return Option{
		name:  "value",
		value: lua.LString(value),
	}
}

func OptionStringArray(values []string) Option {
	array := &lua.LTable{}
	for i, value := range values {
		// +1 because lua is 1 indexed
		array.RawSetInt(i+1, lua.LString(value))
	}

	return Option{
		name:  "value",
		value: array,
	}
}

func OptionFlags(flags []string) Option {
	array := &lua.LTable{}
	for i, flag := range flags {
		array.RawSetInt(i, lua.LString(flag))
	}

	return Option{
		name:  "flags",
		value: array,
	}
}

func OptionAST(ast izu.AST) Option {
	return Option{
		name:  "ast",
		value: lua.LString(ast.String()),
	}
}

func OptionStateHotkey() Option {
	return Option{
		name:  "state",
		value: lua.LNumber(0),
	}
}

func OptionStateBinding() Option {
	return Option{
		name:  "state",
		value: lua.LNumber(1),
	}
}

func OptionStateMultiBinding() Option {
	return Option{
		name:  "state",
		value: lua.LNumber(3),
	}
}

func OptionStateCommand() Option {
	return Option{
		name:  "state",
		value: lua.LNumber(2),
	}
}
