package parser

import (
	"errors"
	"fmt"
	"strings"

	"github.com/meir/izu/pkg/izu"
)

// Multiple is used to parse a branching path in the hotkeys such as
// { x, y }
// where theres a choice between 1 or more keys
type Multiple struct {
	formatter izu.Formatter

	parts []izu.Part
}

// NewMultiple returns a new empty multiple part
func NewMultiple(formatter izu.Formatter) *Multiple {
	return &Multiple{formatter: formatter}
}

// Info returns StateMultiple and the parts that are parsed by it
func (multiple *Multiple) Info() (izu.State, []izu.Part) {
	return izu.StateMultiple, multiple.parts
}

// Parse will parse the data into the multiple part
func (multiple *Multiple) Parse(data []byte) (int, error) {
	for i := 0; i < len(data); i++ {
		char := data[i]
		switch char {
		case '{':
			return len(data), errors.New("unexpected '}'")
		case '}':
			return i + 1, nil
		default:
			base := NewBase(multiple.formatter)
			read, err := base.Parse(data[i:])
			if err != nil {
				return 0, err
			}
			multiple.parts = append(multiple.parts, base)
			i += read
		}
	}
	return len(data), nil
}

// String will return the string representation of the multiple part that has been parsed
func (multiple *Multiple) String() string {
	out := make([]string, len(multiple.parts))
	for i, part := range multiple.parts {
		out[i] = part.String()
	}
	return fmt.Sprintf("{ %v }", strings.Join(out, ", "))
}
