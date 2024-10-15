package izu

type Parser interface {
	Parse([]byte) ([]Hotkey, error)
}

// stateMap is a map that maps the AST type to a readable name for it
var stateMap = map[AST]string{
	ASTHotkey:   "hotkey",
	ASTBinding:  "binding",
	ASTSingle:   "single",
	ASTMultiple: "multiple",
	ASTString:   "string",
}

// String will return the string representation of the state
// This should also be used to define the names of formatter functions
func (state AST) String() string {
	if str, ok := stateMap[state]; ok {
		return str
	}
	panic("invalid state")
}
