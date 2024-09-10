package parser

import (
	"strings"

	"github.com/meir/izu/pkg/izu"
)

// Base will parse the entire shortcut part such as
// Super + { _, Shift +} XF68Media{Play,Pause}
type Command struct {
	formatter izu.Formatter

	parts []izu.Part
}

// NewBase creates a new empty base parser
func NewCommand(formatter izu.Formatter) *Command {
	return &Command{formatter: formatter}
}

// Info returns StateBase and the parts that are parsed by it
func (cmd *Command) Info() (izu.State, []izu.Part) {
	return izu.StateCommand, cmd.parts
}

// Parse will parse the data into the base
func (cmd *Command) Parse(data []byte) (int, error) {
	str := NewString(cmd.formatter)
	for i := 0; i < len(data); i++ {
		char := data[i]
		switch char {
		case '{':
			cmd.parts = append(cmd.parts, str)
			sub := NewSingleSub(cmd.formatter)
			read, err := sub.Parse(data[i+1:])
			if err != nil {
				return 0, err
			}
			cmd.parts = append(cmd.parts, sub)
			i += read
		case '}':
			str = NewString(cmd.formatter)
		case '\n':
			return i, nil
		default:
			str.key += string(char)
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
	return strings.Join(out, " + ")
}
