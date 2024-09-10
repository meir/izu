package izu

import "embed"

// State is the type that defines the AST type
type State uint8

const (
	StateKeybind State = iota
	StateCommand
	StateBase
	StateMultiple
	StateSingle
	StateSinglePart
	StateString
)

// Part is the interface that should be implemented for single AST parts
type Part interface {
	Info() (State, []Part)
	Parse([]byte) (int, error)
	String() string
}

// Formatter is the interface that should be implemented for all hotkey formatters
// These are used to interface with a specific formatter language
type Formatter interface {
	ParseString([]byte) ([]byte, error)
	ParseFile(string) ([]byte, error)
}

//go:embed formatters/*
var Formatters embed.FS

// String will return the string representation of the state
// This should also be used to define the names of formatter functions
func (state State) String() string {
	switch state {
	case StateKeybind:
		return "keybind"
	case StateCommand:
		return "command"
	case StateBase:
		return "base"
	case StateMultiple:
		return "multiple"
	case StateSingle:
		return "single"
	case StateSinglePart:
		return "single_part"
	case StateString:
		return "string"
	default:
		panic("invalid state")
	}
}
