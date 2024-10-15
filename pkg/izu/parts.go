package izu

import (
	"fmt"
	"strings"
)

// State is the type that defines the AST type
type AST uint8

const (
	ASTHotkey AST = iota
	ASTBinding
	ASTSingle
	ASTMultiple
	ASTString
)

// Part is the interface that should be implemented for single AST parts
type Part interface {
	Info() (AST, PartList)
	Append(...Part)
	String() string
}

type PartList interface {
	Iterate(func(Part) error) error
	Append(...Part) PartList
	String() string
}

type DefaultPartList struct {
	pre, seperator, suf string
	parts               []Part
}

func NewDefaultPartList(seperator string, parts ...Part) PartList {
	return DefaultPartList{
		"",
		seperator,
		"",
		parts,
	}
}

func NewDefaultPartListWithNfixes(pre, sep, suf string, parts ...Part) PartList {
	return DefaultPartList{
		pre,
		sep,
		suf,
		parts,
	}
}

func (p DefaultPartList) Iterate(fn func(Part) error) error {
	for _, part := range p.parts {
		err := fn(part)
		if err != nil {
			return err
		}
	}
	return nil
}

func (p DefaultPartList) Append(part ...Part) PartList {
	p.parts = append(p.parts, part...)
	return p
}

func (p DefaultPartList) String() string {
	output := []string{}
	for _, part := range p.parts {
		AST, parts := part.Info()
		switch AST {
		case ASTString:
			output = append(output, part.String())
		default:
			output = append(output, parts.String())
		}
	}
	return fmt.Sprintf("%s%s%s", p.pre, strings.Join(output, p.seperator), p.suf)
}
