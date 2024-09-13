package parser

import (
	"strings"

	"github.com/meir/izu/pkg/izu"
)

// Base will parse the entire shortcut part such as
// Super + { _, Shift +} XF68Media{Play,Pause}
type Command struct {
	parts []izu.Part
}

// NewBase creates a new empty base parser
func NewCommand() *Command {
	return &Command{}
}

// Info returns StateBase and the parts that are parsed by it
func (cmd *Command) Info() (izu.State, []izu.Part) {
	return izu.StateCommand, cmd.parts
}

// Parse will parse the data into the base
func (cmd *Command) Parse(data []byte) (int, error) {
	buffer := ""
	add_buffer := func() {
		str := NewString()
		str.Key = buffer
		cmd.parts = append(cmd.parts, str)
		buffer = ""
	}

	defer add_buffer()

	for i := 0; i < len(data); i++ {
		char := data[i]
		switch char {
		case '{':
			add_buffer()

			sub := NewSingleSub()
			read, err := sub.Parse(data[i+1:])
			if err != nil {
				return 0, err
			}
			cmd.parts = append(cmd.parts, sub)
			i += read
		case '}':
			buffer = ""
		case '\n':
			return i, nil
		default:
			buffer += string(char)
		}
	}
	return len(data), nil
}

// String will return the string representation of the base part that has been parsed
func (cmd *Command) String() string {
	out := make([]string, len(cmd.parts))
	for i, part := range cmd.parts {
		out[i] = part.String()
	}
	return strings.TrimSpace(strings.Join(out, " "))
}
