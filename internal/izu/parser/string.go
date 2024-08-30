package parser

import "github.com/meir/izu/pkg/izu"

type StringPath struct {
	key string
}

func NewStringPath() *StringPath {
	return &StringPath{}
}

func (b *StringPath) Info() (izu.State, []izu.Part) {
	return izu.StateString, []izu.Part{}
}

func (b *StringPath) Parse(s []byte) (int, error) {
	for i := 0; i < len(s); i++ {
		char := s[i]
		switch {
		case char >= 'a' && char <= 'z':
			fallthrough
		case char >= 'A' && char <= 'Z':
			fallthrough
		case char >= '0' && char <= '9':
			b.key += string(char)
		default:
			return i - 1, nil
		}
	}
	return len(s), nil
}

func (b *StringPath) String() string {
	return b.key
}
