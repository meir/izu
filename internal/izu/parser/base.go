package parser

import (
	"strings"

	"github.com/meir/izu/pkg/izu"
)

// Base will parse the entire shortcut part such as
// Super + { _, Shift +} XF68Media{Play,Pause}
type Base struct {
	formatter izu.Formatter

	parts []izu.Part
}

// NewBase creates a new empty base parser
func NewBase(formatter izu.Formatter) *Base {
	return &Base{formatter: formatter}
}

// Info returns StateBase and the parts that are parsed by it
func (base *Base) Info() (izu.State, []izu.Part) {
	return izu.StateBase, base.parts
}

// Parse will parse the data into the base
func (base *Base) Parse(data []byte) (int, error) {
	for i := 0; i < len(data); i++ {
		char := data[i]
		switch char {
		case '{':
			multiple := NewMultiple(base.formatter)
			read, err := multiple.Parse(data[i+1:])
			if err != nil {
				return 0, err
			}
			base.parts = append(base.parts, multiple)
			i += read
		case '}':
			return i - 1, nil
		case ',':
			return i, nil
		case ' ', '+':
			continue
		case '\n':
			return i, nil
		default:
			single := NewSingle(base.formatter)
			read, err := single.Parse(data[i:])
			if err != nil {
				return 0, err
			}
			if !(single.parts == nil || len(single.parts) == 0) {
				base.parts = append(base.parts, single)
				i += read
			}
		}
	}
	return len(data), nil
}

// String will return the string representation of the base part that has been parsed
func (base *Base) String() string {
	out := make([]string, len(base.parts))
	for i, part := range base.parts {
		out[i] = part.String()
	}
	return strings.Join(out, " + ")
}
