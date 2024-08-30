package parser

import (
	"strings"
)
import "github.com/meir/izu/pkg/izu"

type BasePath struct {
	parts []izu.Part
}

func NewBasePath() *BasePath {
	return &BasePath{}
}

func (b *BasePath) Info() (izu.State, []izu.Part) {
	return izu.StateBase, b.parts
}

func (b *BasePath) Parse(s []byte) (int, error) {
	for i := 0; i < len(s); i++ {
		char := s[i]
		switch char {
		case '{':
			np := NewMultiPath()
			rs, err := np.Parse(s[i+1:])
			if err != nil {
				return 0, err
			}
			b.parts = append(b.parts, np)
			i += rs
		case '}':
			return i - 1, nil
		case ',':
			return i, nil
		case ' ', '+':
			continue
		default:
			sp := NewSinglePath()
			rs, err := sp.Parse(s[i:])
			if err != nil {
				return 0, err
			}
			if !sp.isEmpty() {
				b.parts = append(b.parts, sp)
				i += rs
			}
		}
	}
	return len(s), nil
}

func (b *BasePath) String() string {
	out := make([]string, len(b.parts))
	for i, p := range b.parts {
		out[i] = p.String()
	}
	return strings.Join(out, " + ")
}
