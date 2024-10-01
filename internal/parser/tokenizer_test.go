package parser

import (
	"testing"
)

func TestTokenizer(t *testing.T) {
	type Output struct {
		kind TokenKind
		text string
	}

	testcases := []struct {
		input  string
		output []Output
	}{
		{
			"this is a test",
			[]Output{
				{TokenString, "this"},
				{TokenEmpty, " "},
				{TokenString, "is"},
				{TokenEmpty, " "},
				{TokenString, "a"},
				{TokenEmpty, " "},
				{TokenString, "test"},
			},
		},
		{
			"Super w\necho \"Hello there\"",
			[]Output{
				{TokenString, "Super"},
				{TokenEmpty, " "},
				{TokenString, "w"},
				{TokenNewLine, "\n"},
				{TokenString, "echo"},
				{TokenEmpty, " "},
				{TokenOther, "\""},
				{TokenString, "Hello"},
				{TokenEmpty, " "},
				{TokenString, "there"},
				{TokenOther, "\""},
			},
		},
		{
			"Super + XF86Audio{Play,Pause}\nplayerctl {play,pause}",
			[]Output{
				{TokenString, "Super"},
				{TokenEmpty, " "},
				{TokenPlus, "+"},
				{TokenEmpty, " "},
				{TokenString, "XF86Audio"},
				{TokenMultiOpen, "{"},
				{TokenString, "Play"},
				{TokenMultiDivide, ","},
				{TokenString, "Pause"},
				{TokenMultiClose, "}"},
				{TokenNewLine, "\n"},
				{TokenString, "playerctl"},
				{TokenEmpty, " "},
				{TokenMultiOpen, "{"},
				{TokenString, "play"},
				{TokenMultiDivide, ","},
				{TokenString, "pause"},
				{TokenMultiClose, "}"},
			},
		},
		{
			"Super + w | hyprland[r]\nhyprland|echo \"Hello there\"",
			[]Output{
				{TokenString, "Super"},
				{TokenEmpty, " "},
				{TokenPlus, "+"},
				{TokenEmpty, " "},
				{TokenString, "w"},
				{TokenEmpty, " "},
				{TokenSystem, "|"},
				{TokenEmpty, " "},
				{TokenString, "hyprland"},
				{TokenFlagOpen, "["},
				{TokenString, "r"},
				{TokenFlagClose, "]"},
				{TokenNewLine, "\n"},
				{TokenString, "hyprland"},
				{TokenSystem, "|"},
				{TokenString, "echo"},
				{TokenEmpty, " "},
				{TokenOther, "\""},
				{TokenString, "Hello"},
				{TokenEmpty, " "},
				{TokenString, "there"},
				{TokenOther, "\""},
			},
		},
	}

	for _, tc := range testcases {
		tk := NewTokenizer([]byte(tc.input))

		i := 0
		for tk.Next() {
			token := tk.Current()
			if token.kind != tc.output[i].kind {
				t.Errorf("expected token kind %d, got %d", tc.output[i].kind, token.kind)
			}

			if string(token.value) != tc.output[i].text {
				t.Errorf("expected token text %s, got %s", tc.output[i].text, string(token.value))
			}
			i++
		}

		if i != len(tc.output) {
			t.Errorf("expected %d tokens, got %d", len(tc.output), i)
		}
	}
}
