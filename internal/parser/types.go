package parser

import (
	"github.com/meir/izu/pkg/izu"
)

// ---

// PartBinding is a type that represents a binding of multiple parts,
// this is basically the highest part of a hotkey binding there is
type PartBinding struct {
	parts izu.PartList
}

// Info returns ASTBinding and the binding partlist
func (p *PartBinding) Info() (izu.AST, izu.PartList) {
	return izu.ASTBinding, p.parts
}

// Append appends a part to the binding
func (p *PartBinding) Append(part ...izu.Part) {
	p.parts = p.parts.Append(part...)
}

// String returns the string representation
func (p *PartBinding) String() string {
	return p.parts.String()
}

// ---

// PartSingle is a type that represents a single part,
// This can contain a String or a Multiple
type PartSingle struct {
	parts izu.PartList
}

// Info returns ASTSingle and the partlist
func (p *PartSingle) Info() (izu.AST, izu.PartList) {
	return izu.ASTSingle, p.parts
}

// Append appends a part to the single
func (p *PartSingle) Append(part ...izu.Part) {
	p.parts = p.parts.Append(part...)
}

// String returns the string representation
func (p *PartSingle) String() string {
	return p.parts.String()
}

// ---

// PartMultiple is a type that represents a multiple part,
// This can contain multiple bindings
type PartMultiple struct {
	parts izu.PartList
}

// Info returns ASTMultiple and the partlist
func (p *PartMultiple) Info() (izu.AST, izu.PartList) {
	return izu.ASTMultiple, p.parts
}

// Append appends a part to the multiple
func (p *PartMultiple) Append(part ...izu.Part) {
	p.parts = p.parts.Append(part...)
}

// String returns the string representation
func (p *PartMultiple) String() string {
	return p.parts.String()
}

// ---

// PartString is a type that represents a single string
// This is the lowest part of a binding there is
type PartString struct {
	value string
}

// NewPartString creates a new PartString
func NewPartString(v string) *PartString {
	return &PartString{value: v}
}

// Info returns ASTString and nil
// The value has to be retrieved using the String method
func (p *PartString) Info() (izu.AST, izu.PartList) {
	return izu.ASTString, nil
}

// Append sets the string value to the first item given
func (p *PartString) Append(part ...izu.Part) {
	if len(part) > 0 {
		p.value = part[0].String()
	}
}

// String returns the string value
func (p *PartString) String() string {
	return p.value
}
