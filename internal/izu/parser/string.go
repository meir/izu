package parser

import (
	"strings"

	"github.com/meir/izu/pkg/izu"
)

// String will parse a string part such as
// XF68AudioPlay
// A
// Super
type String struct {
	key string
}

// NewString creates a new empty string part
func NewString() *String {
	return &String{}
}

// Info returns StateString and the parts that are parsed by it
func (str *String) Info() (izu.State, []izu.Part) {
	return izu.StateString, []izu.Part{}
}

// Key will return the key that has been parsed
func (str *String) Key() string {
	return strings.TrimSpace(str.key)
}

// Parse will parse the data into the string part
func (str *String) Parse(data []byte) (int, error) {
	for i := 0; i < len(data); i++ {
		char := data[i]
		switch {
		case char == '_', char == '-':
			fallthrough
		case char >= 'a' && char <= 'z':
			fallthrough
		case char >= 'A' && char <= 'Z':
			fallthrough
		case char >= '0' && char <= '9':
			str.key += string(char)
		default:
			return i - 1, nil
		}
	}
	return len(data), nil
}

// String will return the string representation of the string part that has been parsed
func (str *String) String() string {
	return strings.TrimSpace(str.key)
}
