package parser

import "strings"
import "github.com/meir/izu/pkg/izu"

type Path struct {
	parts []izu.Part
}

func NewPath(parts []izu.Part) *Path {
	return &Path{parts: parts}
}

func (p *Path) String() string {
	if len(p.parts) == 1 {
		return p.parts[0].String()
	}

	out := []string{}
	for _, part := range p.parts {
		out = append(out, part.String())
	}
	return strings.Join(out, " + ")
}
