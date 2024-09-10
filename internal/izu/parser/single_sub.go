package parser

import (
	"errors"
	"fmt"
	"strings"

	"github.com/meir/izu/pkg/izu"
)

// SingleSub will parse a branching part of a single part such as
// XF68Media{Play,Pause}
type SingleSub struct {
	parts []string
}

// NewSingleSub creates a new empty single sub part
func NewSingleSub() *SingleSub {
	return &SingleSub{}
}

// Info returns StateSinglePart and the parts that are parsed by it
func (sub *SingleSub) Info() (izu.State, []izu.Part) {
	parts := make([]izu.Part, len(sub.parts))
	for i, part := range sub.parts {
		str := NewString()
		str.key = part
		parts[i] = str
	}
	return izu.StateSinglePart, parts
}

// Parse will parse the data into the single sub part
func (sub *SingleSub) Parse(data []byte) (int, error) {
	sub.parts = append(sub.parts, "")
	for i := 0; i < len(data); i++ {
		char := data[i]
		switch char {
		case '{':
			return i, errors.New("unexpected '{'")
		case '}', '\n':
			if sub.parts[len(sub.parts)-1] == "" {
				sub.parts = sub.parts[:len(sub.parts)-1]
			}
			return i + 1, nil
		case ',':
			sub.parts = append(sub.parts, "")
		case '+':
			return i, errors.New("unexpected '" + string(char) + "'")
		default:
			sub.parts[len(sub.parts)-1] += string(char)
		}
	}
	return len(data), nil
}

// String will return the string representation of the single sub part that has been parsed
func (sub *SingleSub) String() string {
	return fmt.Sprintf("{%v}", strings.Join(sub.parts, ","))
}
