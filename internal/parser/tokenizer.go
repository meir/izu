package parser

import (
	"fmt"
	"log/slog"
	"slices"
)

// TokenKind is the type that defines the kind of token that was found
type TokenKind uint8

const (
	TokenEOF TokenKind = iota
	TokenString
	TokenEmpty
	TokenPlus
	TokenComment
	TokenNewLine
	TokenSemicolon
	TokenMultiOpen
	TokenMultiClose
	TokenMultiDivide
	TokenSystem
	TokenFlagOpen
	TokenFlagClose

	TokenOther
)

// Token is the type that defines the token
// This includes the kind, position and the value
type Token struct {
	kind TokenKind
	line int
	col  int

	value []byte
}

// NewToken creates a new token
func NewToken(data []byte, kind TokenKind, line, col int) Token {
	return Token{
		kind,
		line,
		col,
		data,
	}
}

// Kind returns the kind of the token
func (t Token) Kind() TokenKind {
	return t.kind
}

// Match checks if the token matches the given string
func (t Token) Match(s string) bool {
	return t.String() == s
}

// Position returns the position of the token formatted in "line:column"
func (t Token) Position() string {
	return fmt.Sprintf("%d:%d", t.line, t.col)
}

// String returns the value of the token
func (t Token) String() string {
	return string(t.value)
}

// ---

// Tokenizer is a type that stores an array of tokens and loops through it using methods
type Tokenizer struct {
	tokens []Token

	index int
}

// NewTokenizer creates a new tokenizer from the data and tokenizes everything
func NewTokenizer(data []byte) *Tokenizer {
	return tokenize(data)
}

// NewTokenizerFromTokens creates a tokenizer without parsing the tokens, it directly uses the tokens given
func NewTokenizerFromTokens(tokens []Token) *Tokenizer {
	return &Tokenizer{
		tokens: tokens,
		index:  -1,
	}
}

// tokenize is a helper function that reads through the data and turns everything into a Token
func tokenize(data []byte) *Tokenizer {
	slog.Debug("Tokenizing data", "data", string(data))
	tokens := []Token{}
	// keep track of the line and column
	line := 1
	col := 0
	// map of all the tokens that are not strings
	tokenMap := map[byte]TokenKind{
		'+':  TokenPlus,
		'#':  TokenComment,
		'\n': TokenNewLine,
		';':  TokenSemicolon,
		'{':  TokenMultiOpen,
		'}':  TokenMultiClose,
		',':  TokenMultiDivide,
		'|':  TokenSystem,
		'[':  TokenFlagOpen,
		']':  TokenFlagClose,
	}

	// accumulate_token is a helper function that tries to string tokens together into the same type
	accumulate_token := func(char byte, kind TokenKind) {
		// check if theres a token yet
		if len(tokens) == 0 {
			// if not, create a new token
			tokens = append(tokens, NewToken([]byte{char}, kind, line, col))
			return
		}

		// get the last token
		last := tokens[len(tokens)-1]
		// if the token is of the same type, append the current byte to its value
		if last.kind == kind {
			tokens[len(tokens)-1].value = append(last.value, char)
			return
		}
		// otherwise start a new token of this type
		tokens = append(tokens, NewToken([]byte{char}, kind, line, col))
	}

	// loop through all the data
	for i := 0; i < len(data); i++ {
		char := data[i]

		// update the line and column
		col++
		if char == '\n' {
			line++
			col = 0
		}

		switch {
		case char >= 'A' && char <= 'Z': // check for A-Z
			fallthrough // fallthrough to make it more readable
		case char >= 'a' && char <= 'z': // check for a-z
			fallthrough
		case char >= '0' && char <= '9': // check for 0-9
			fallthrough
		case char == '_' || char == '-': // check for - or _
			// using the fallthroughs we basically have [A-Za-z0-9_-]+ without using regex
			accumulate_token(char, TokenString)
		case char == ' ' || char == '\t':
			// if its a space or a tab, accumulate it as an empty token
			// these cant be ignores because we need them for commands
			accumulate_token(char, TokenEmpty)
		default:
			// find any other tokens from the tokenMap and otherwise add them as TokenOther,
			// because we cant crash on them as they might be part of the command
			kind := TokenOther
			if k, ok := tokenMap[char]; ok {
				kind = k
			}
			tokens = append(tokens, NewToken([]byte{char}, kind, line, col))
		}
	}

	slog.Debug("Tokenized data", "tokens", len(tokens))

	return &Tokenizer{
		tokens: tokens,
		index:  -1,
	}
}

// Next moves to the next index and returns a boolean if index is still within range
func (t *Tokenizer) Next() bool {
	t.index++
	return t.index < len(t.tokens)
}

// Current returns the current token and otherwise EOF
func (t *Tokenizer) Current() Token {
	if t.index < 0 || t.index >= len(t.tokens) {
		return Token{
			kind: TokenEOF,
		}
	}
	return t.tokens[t.index]
}

// Previous moves to the previous index
func (t *Tokenizer) Previous() {
	t.index--
	// check for -1 because otherwise tokenizer.Next() will automatically get item 1 instead of 0
	if t.index < -1 {
		t.index = -1
	}
}

// Peek returns the next token without moving the index or EOF
func (t *Tokenizer) Peek(ignore ...TokenKind) Token {
	for i := t.index + 1; i < len(t.tokens); i++ {
		token := t.tokens[i]
		if !slices.Contains(ignore, token.kind) {
			return token
		}
	}
	return Token{
		kind: TokenEOF,
	}
}

// SkipTo moves to the next token that matches the kind + line + column + value
func (t *Tokenizer) SkipTo(token Token) {
	for {
		current := t.Current()

		kind := current.Kind() == token.Kind()
		line := current.line == token.line
		column := current.col == token.col
		value := current.String() == token.String()
		if kind && line && column && value {
			break
		}

		if !t.Next() {
			break
		}
	}
}

// Until moves through the tokens until it finds the given kind and returns the [current token + in between tokens], final token
func (t *Tokenizer) Until(kind ...TokenKind) ([]Token, Token) {
	tokens := []Token{t.Current()}
	for t.Next() {
		token := t.Current()

		if slices.Contains(kind, token.kind) {
			break
		}

		tokens = append(tokens, token)
	}
	return tokens, t.Current()
}

// UntilNot moves through the tokens until it finds a token that is not in the given kind and returns the [current token + in between tokens], final
// This can be used to skip empty tokens
func (t *Tokenizer) UntilNot(kind ...TokenKind) ([]Token, Token) {
	tokens := []Token{t.Current()}
	for t.Next() {
		token := t.Current()

		if !slices.Contains(kind, token.kind) {
			break
		}

		tokens = append(tokens, token)
	}
	return tokens, t.Current()
}
