package parser

import (
	"strings"

	"github.com/meir/izu/pkg/izu"
)

// Single will parse a single part such as
// A
// XF68AudioPlay
// XF86Audio{Play,Pause}
type Single struct {
	formatter izu.Formatter

	parts []izu.Part
}

// NewSingle creates a new empty single part
func NewSingle(formatter izu.Formatter) *Single {
	return &Single{formatter: formatter}
}

// Info returns StateSingle and the parts parsed by it
func (single *Single) Info() (izu.State, []izu.Part) {
	return izu.StateSingle, single.parts
}

// Parse will parse the data into the single part
func (single *Single) Parse(data []byte) (int, error) {
	for i := 0; i < len(data); i++ {
		char := data[i]
		switch char {
		case '{':
			single_sub := NewSingleSub(single.formatter)
			read, err := single_sub.Parse(data[i+1:])
			if err != nil {
				return 0, err
			}
			single.parts = append(single.parts, single_sub)
			i += read
		case ',', '}', '\n':
			return i - 1, nil
		case ' ', '+':
			return i, nil
		default:
			str := NewString(single.formatter)
			read, err := str.Parse(data[i:])
			if err != nil {
				return 0, err
			}
			if strings.TrimSpace(str.key) != "" {
				single.parts = append(single.parts, str)
				i += read
			}
		}
	}
	return len(data), nil
}

// String will return the string representation of the single part that has been parsed
func (single *Single) String() string {
	out := make([]string, len(single.parts))
	for i, part := range single.parts {
		out[i] = part.String()
	}
	return strings.Join(out, "")
}
