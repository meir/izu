package parser

import (
	"fmt"
	"strings"

	"github.com/meir/izu/pkg/izu"
)

// Base will parse the entire shortcut part such as
// Super + { _, Shift +} XF68Media{Play,Pause}
type Keybind struct {
	parts []izu.Part
}

// NewKeybind creates a new empty base parser
func NewKeybind() *Keybind {
	return &Keybind{}
}

// Info returns StateBase and the parts that are parsed by it
func (keybind *Keybind) Info() (izu.State, []izu.Part) {
	return izu.StateKeybind, keybind.parts
}

// Parse will parse the data into the base
func (keybind *Keybind) Parse(data []byte) (int, error) {
	if len(data) == 0 {
		return 0, fmt.Errorf("empty keybind")
	}

	bind := NewBase()
	var read int
	var err error
	if read, err = bind.Parse(data); err != nil {
		return 0, err
	}

	if read >= len(data) {
		return 0, fmt.Errorf("unexpected end of keybind")
	}

	if data[read] != '\n' {
		// base part probably ended with } instead of \n so this is unexpected and should error
		return 0, fmt.Errorf("unexpected '%c'", data[read])
	}

	command := NewCommand()
	if _, err := command.Parse(data[read+1:]); err != nil {
		return 0, err
	}

	keybind.parts = []izu.Part{bind, command}

	return len(data), nil
}

// String will return the string representation of the base part that has been parsed
func (keybind *Keybind) String() string {
	out := make([]string, len(keybind.parts))
	for i, part := range keybind.parts {
		out[i] = part.String()
	}
	return strings.Join(out, "\n\t")
}
