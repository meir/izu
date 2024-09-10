package parser

import (
	"strings"

	"github.com/meir/izu/pkg/izu"
)

// Base will parse the entire shortcut part such as
// Super + { _, Shift +} XF68Media{Play,Pause}
type Keybind struct {
	formatter izu.Formatter

	parts []izu.Part
}

// NewKeybind creates a new empty base parser
func NewKeybind(formatter izu.Formatter) *Keybind {
	return &Keybind{formatter: formatter}
}

// Info returns StateBase and the parts that are parsed by it
func (keybind *Keybind) Info() (izu.State, []izu.Part) {
	return izu.StateKeybind, keybind.parts
}

// Parse will parse the data into the base
func (keybind *Keybind) Parse(data []byte) (int, error) {
	bind := NewBase(keybind.formatter)
	var read int
	var err error
	if read, err = bind.Parse(data); err != nil {
		return 0, err
	}

	command := NewCommand(keybind.formatter)
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
