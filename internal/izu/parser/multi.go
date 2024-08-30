package parser

import (
	"errors"
	"fmt"
	"strings"
)
import "github.com/meir/izu/pkg/izu"

type MultiPath struct {
	parts []izu.Part
}

func NewMultiPath() *MultiPath {
	return &MultiPath{}
}

func (b *MultiPath) Info() (izu.State, []izu.Part) {
	return izu.StateMulti, b.parts
}

func (b *MultiPath) Parse(s []byte) (int, error) {
	for i := 0; i < len(s); i++ {
		char := s[i]
		switch char {
		case '{':
			return len(s), errors.New("unexpected '}'")
		case '}':
			return i + 1, nil
		default:
			bp := NewBasePath()
			rs, err := bp.Parse(s[i:])
			if err != nil {
				return 0, err
			}
			b.parts = append(b.parts, bp)
			i += rs
		}
	}
	return len(s), nil
}

func (b *MultiPath) String() string {
	out := make([]string, len(b.parts))
	for i, p := range b.parts {
		out[i] = p.String()
	}
	return fmt.Sprintf("{ %v }", strings.Join(out, ", "))
}
