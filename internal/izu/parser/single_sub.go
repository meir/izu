package parser

import (
	"errors"
	"fmt"
	"strings"

	"github.com/meir/izu/pkg/izu"
)

type SingleSubPath struct {
	parts []string
}

func NewSingleSubPath() *SingleSubPath {
	return &SingleSubPath{}
}

func (b *SingleSubPath) Info() (izu.State, []izu.Part) {
	return izu.StateSinglePart, []izu.Part{}
}

func (b *SingleSubPath) Parse(s []byte) (int, error) {
	b.parts = append(b.parts, "")
	for i := 0; i < len(s); i++ {
		char := s[i]
		switch char {
		case '{':
			return i, errors.New("unexpected '{'")
		case '}':
			if b.parts[len(b.parts)-1] == "" {
				b.parts = b.parts[:len(b.parts)-1]
			}
			return i + 1, nil
		case ',':
			b.parts = append(b.parts, "")
		case ' ', '+':
			return i, errors.New("unexpected '" + string(char) + "'")
		default:
			b.parts[len(b.parts)-1] += string(char)
		}
	}
	return len(s), nil
}

func (b *SingleSubPath) String() string {
	return fmt.Sprintf("{%v}", strings.Join(b.parts, ","))
}
