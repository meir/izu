package parser

import (
	"fmt"

	"github.com/meir/izu/pkg/izu"
)

// ParserState is the type that defines the state of the parser
// This is used to keep track of what stage the parser is in
type ParserState uint8

const (
	StateRoot ParserState = iota
	StateBinding
	StateFlags
	StateCommand
)

// unexpectedToken is a helper function that returns an error
func unexpectedToken(token Token, state ParserState) error {
	stateMap := map[ParserState]string{
		StateRoot:    "root",
		StateBinding: "binding",
		StateFlags:   "flags",
		StateCommand: "command",
	}
	return fmt.Errorf("unexpected token '%s' at %s (state %v)\n", token.String(), token.Position(), stateMap[state])
}

// filter is a helper function that filters a slice based on a test function
func filter[T any](ss []T, test func(T) bool) (ret []T) {
	for _, s := range ss {
		if test(s) {
			ret = append(ret, s)
		}
	}
	return
}

// parseBinding is a helper function that is used to parse the binding part of a hotkey
func parseBinding(parent izu.Part, tokenizer *Tokenizer) error {
	// loop through the given tokenizer
	for tokenizer.Next() {
		token := tokenizer.Current()
		switch token.Kind() {
		case TokenString:
			// if the token is a string, create a new single part with the string token
			single := &PartSingle{
				parts: izu.NewDefaultPartList(" + ", &PartString{token.String()}),
			}

			// if the next token is a multiple open token, then it should be part of the single
			// this means its something like "XF86Audio{Play,Pause}" otherwise there would be an empty token in between
			if tokenizer.Peek().Kind() == TokenMultiOpen {
				// continue the parser with the single as its parent
				err := parseBinding(single, tokenizer)
				if err != nil {
					return err
				}
			}

			// add the single to the parent
			parent.Append(single)

		case TokenMultiOpen:
			// create a new multiple part
			multiple := &PartMultiple{izu.NewDefaultPartListWithNfixes("{", ",", "}")}

			// skip the { so that the next parser wont be stuck on it and create an infinite loop
			tokenizer.Next()
			// get all the tokens till the closing part
			tokens, _ := tokenizer.Until(TokenMultiClose)
			subtokenizer := NewTokenizerFromTokens(tokens)
		Loop:
			for {
				// create a new binding and start parsing using that as the parent
				// this binding will be one of the paths in the multiple, such as {binding,binding}
				binding := &PartBinding{izu.NewDefaultPartList(" + ")}
				err := parseBinding(binding, subtokenizer)
				if err != nil {
					return err
				}

				// append the binding
				multiple.parts = multiple.parts.Append(binding)
				// if the subtokenizer ends, break out of this loop
				switch subtokenizer.Peek().Kind() {
				case TokenMultiClose, TokenEOF:
					break Loop
				}
			}

			// add the multiple to the parent
			parent.Append(multiple)

		case TokenPlus, TokenEmpty:
			// skip empty tokens and plus signs
			// we just ignore those since they wont give any additional context to the final binding
			continue

		case TokenSemicolon, TokenNewLine, TokenMultiDivide:
			// if the token is a semicolon, newline or multi divide, return
			return nil

		default:
			// if we find any other tokens that arent handled by our cases, we should error on this
			return unexpectedToken(token, StateBinding)
		}
	}
	return nil
}

// parseCommand is a helper function that is used to parse the command part of a hotkey
func parseCommand(parent izu.Part, tokenizer *Tokenizer) error {
	// loop through the given tokenizer
	for tokenizer.Next() {
		token := tokenizer.Current()
		switch token.Kind() {
		case TokenMultiOpen:
			// create a new multiple part
			multiple := &PartMultiple{izu.NewDefaultPartListWithNfixes("{", ",", "}")}

			// skip the { so that the next parser wont be stuck on it and create an infinite loop
			tokenizer.Next()
			// get all the tokens till the closing part
			tokens, _ := tokenizer.Until(TokenMultiClose)
			subtokenizer := NewTokenizerFromTokens(tokens)

			// create a new binding and start parsing using that as the parent
			// this binding will be one of the paths in the multiple, such as {binding,binding}
			binding := &PartBinding{izu.NewDefaultPartList("")}
			multiple.Append(binding)
		Loop:
			for subtokenizer.Next() {
				token := subtokenizer.Current()

				switch token.Kind() {
				case TokenMultiDivide:
					binding = &PartBinding{izu.NewDefaultPartList("")}
					multiple.Append(binding)
				case TokenMultiClose, TokenEOF:
					break Loop
				default:
					binding.Append(&PartString{
						value: token.String(),
					})
				}
			}

			// add the multiple to the parent
			parent.Append(multiple)

		default:
			// dump everything into a string and into the parent
			parent.Append(&PartString{
				value: token.String(),
			})
		}
	}
	return nil
}

// --- parser states ---

// stateRoot is the parser state for the root of the parser
func stateRoot(tokenizer *Tokenizer, hotkeys *[]*izu.Hotkey, state *ParserState) error {
	token := tokenizer.Current()

	switch token.Kind() {
	case TokenEmpty, TokenNewLine:
		// skip any empty lines
		return nil
	case TokenComment:
		// if theres a comment, skip to the next line and ignore anything in between
		tokenizer.Until(TokenNewLine)
	case TokenString, TokenMultiOpen:
		// if we get a string or multi open, we should start parsing
		// go to the previous token so that the parser can get it
		tokenizer.Previous()
		*state = StateBinding
	default:
		// if we get any other token, we should error
		return unexpectedToken(token, *state)
	}
	return nil
}

// stateBinding is the parser state for the binding of the parser
func stateBinding(tokenizer *Tokenizer, hotkeys *[]*izu.Hotkey, state *ParserState) error {
	// get the content of the binding, bindings will always end with either a semicolon, a newline of a pipe to specify the flags
	binding, token := tokenizer.Until(TokenSemicolon, TokenNewLine, TokenSystem)

	// start parsing the binding
	bindingPart := &PartBinding{izu.NewDefaultPartList(" + ")}
	err := parseBinding(bindingPart, NewTokenizerFromTokens(binding))
	if err != nil {
		return err
	}

	// add the hotkey, subsequent states will fill this hotkey further
	*hotkeys = append(*hotkeys, &izu.Hotkey{
		Binding: bindingPart,
		Command: map[string]izu.Part{},
		Flags:   map[string][]string{},
	})

	switch token.Kind() {
	case TokenSemicolon, TokenNewLine:
		// if we find a semicolon or newline, get to the command state
		*state = StateCommand
	case TokenSystem:
		// if we find a pipe, go to the flag state
		*state = StateFlags
	default:
		return unexpectedToken(token, *state)
	}
	return nil
}

// stateFlags is the parser state for the flags of the parser
// we parse the states directly, since its fairly easy
// the format is always | ([A-Za-z0-9_-]+\[([A-Za-z0-9_-]]+)+\])+
func stateFlags(tokenizer *Tokenizer, hotkeys *[]*izu.Hotkey, state *ParserState) error {
	flags := map[string][]string{}
	name := ""
	values := []string{}

FlagLoop:
	for tokenizer.Next() {
		token := tokenizer.Current()

		switch token.Kind() {
		case TokenEmpty:
			// just ignore any empty tokens
			// we dont want any yucky spaces in our flags
			continue
		case TokenString:
			// if theres a string, save it as the name, and check if the next token is a [
			// if its not, its not properly formatted
			// flags should always be system[flag] without spaces in between
			name = token.String()
			next := tokenizer.Peek()
			if next.Kind() != TokenFlagOpen {
				return unexpectedToken(next, *state)
			}
		case TokenFlagOpen:
			values = []string{}

			// call next so that TokenFlagOpen token wont be included in parsing the flag items
			tokenizer.Next()
			flaglist, _ := tokenizer.Until(TokenFlagClose)
			for _, flag := range flaglist {
				// if theres an empty token, just skip it
				if flag.Kind() == TokenEmpty {
					continue
				}

				// if theres a token thats not an id, error
				if flag.Kind() != TokenString {
					return unexpectedToken(flag, *state)
				}
				// add the flag to the current system
				values = append(values, flag.String())
			}

			// if the flag already exists, that means the user specified the system twice, error on this
			if _, ok := flags[name]; ok {
				return fmt.Errorf("flag '%s' already exists", name)
			}

			// flag ends here, add the values to the system's name
			// and reset the name for the next iteration
			flags[name] = values
			name = ""
			values = []string{}

		case TokenNewLine, TokenSemicolon:
			// if theres a new line or a semicolon, skip to the command state
			*state = StateCommand
			break FlagLoop
		}
	}

	(*hotkeys)[len(*hotkeys)-1].Flags = flags
	return nil
}

// stateCommand is the parser state for the command of the parser
func stateCommand(tokenizer *Tokenizer, hotkeys *[]*izu.Hotkey, state *ParserState) error {
	// check up to the next newline or for a pipe
	// if theres a pipe, it means a system has been specified but it might also mean its a system identifier
	command, token := tokenizer.Until(TokenNewLine, TokenSystem)
	system := "default" // default system, this will be omitted when printed

	switch token.Kind() {
	case TokenSystem:
		// if the token is a system token, check if theres only 1 token in between thats not empty
		// filter for all the tokens that arent empty
		if pre := filter(command, func(t Token) bool {
			return t.Kind() != TokenEmpty
		}); len(pre) != 1 {
			// multiple components before system token
			// so its probably part of a command
			rest, _ := tokenizer.Until(TokenNewLine)
			command = append(command, rest...)
		} else {
			system = pre[0].String()
			//skip until first non empty
			tokenizer.UntilNot(TokenEmpty)
			command, _ = tokenizer.Until(TokenNewLine)
		}
	}

	commandBinding := &PartBinding{izu.NewDefaultPartList("")}
	commandTokenizer := NewTokenizerFromTokens(command)
	// skip prefix empty spaces
	commandTokenizer.UntilNot(TokenEmpty)
	// because we want the command parser to start with index-1 so that the first Next() will be at the start
	commandTokenizer.Previous()
	err := parseCommand(commandBinding, commandTokenizer)
	if err != nil {
		return err
	}
	// because we want the command parser to start with index-1 so that the first Next() will be at the start
	(*hotkeys)[len(*hotkeys)-1].Command[system] = commandBinding

	_, token = tokenizer.UntilNot(TokenEmpty)
	if token.Kind() == TokenNewLine {
		*state = StateRoot
	}
	return nil
}

// Parse will parse the given data into a list of hotkeys or an error
// Check the README.md or the example folder to see what the syntax is
func Parse(data []byte) ([]*izu.Hotkey, error) {
	tokenizer := NewTokenizer(data)
	hotkeys := []*izu.Hotkey{}
	state := StateRoot

	// stateMap is a map that contains the state functions
	stateMap := map[ParserState]func(*Tokenizer, *[]*izu.Hotkey, *ParserState) error{
		StateRoot:    stateRoot,
		StateBinding: stateBinding,
		StateFlags:   stateFlags,
		StateCommand: stateCommand,
	}

	// loop through the tokenizer
	for tokenizer.Next() {
		if stateFunc, ok := stateMap[state]; ok {
			err := stateFunc(tokenizer, &hotkeys, &state)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, fmt.Errorf("unknown state %v", state)
		}
	}

	return hotkeys, nil
}
