package parser

import (
	"github.com/meir/izu/pkg/izu"
	"strings"
)

type SinglePath struct {
	parts []izu.Part
}

func NewSinglePath() *SinglePath {
	return &SinglePath{}
}

func (b *SinglePath) Info() (izu.State, []izu.Part) {
	return izu.StateSingle, b.parts
}

func (b *SinglePath) Parse(s []byte) (int, error) {
	for i := 0; i < len(s); i++ {
		char := s[i]
		switch char {
		case '{':
			np := NewSingleSubPath()
			rs, err := np.Parse(s[i+1:])
			if err != nil {
				return 0, err
			}
			b.parts = append(b.parts, np)
			i += rs
		case ',', '}':
			return i - 1, nil
		case ' ', '+':
			return i, nil
		default:
			sp := NewStringPath()
			rs, err := sp.Parse(s[i:])
			if err != nil {
				return 0, err
			}
			if sp.key != "" {
				b.parts = append(b.parts, sp)
				i += rs
			}
		}
	}
	return len(s), nil
}

func (b *SinglePath) isEmpty() bool {
	return b.parts == nil || len(b.parts) == 0
}

func (b *SinglePath) String() string {
	out := make([]string, len(b.parts))
	for i, p := range b.parts {
		out[i] = p.String()
	}
	return strings.Join(out, "")
}
