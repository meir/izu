package parser

import (
	"testing"

	"github.com/meir/izu/pkg/izu"
)

func TestParser(t *testing.T) {
	cases := []struct {
		input string

		hotkeys []izu.Hotkey
	}{
		{
			input: `a + b + c; echo hello`,
			hotkeys: []izu.Hotkey{
				{
					Binding: &PartBinding{
						izu.NewDefaultPartList(
							" + ",
							&PartSingle{izu.NewDefaultPartList(" + ", &PartString{"a"})},
							&PartSingle{izu.NewDefaultPartList(" + ", &PartString{"b"})},
							&PartSingle{izu.NewDefaultPartList(" + ", &PartString{"c"})},
						),
					},
					Command: map[string]izu.Part{
						"default": &PartBinding{
							parts: izu.NewDefaultPartList(
								"",
								// commands use the same system but are parsed differently to store as much of the original command as possible
								&PartSingle{izu.NewDefaultPartList(
									"",
									&PartString{"echo"},
									&PartString{" "},
									&PartString{"hello"},
								)},
							),
						},
					},
					Flags: map[string][]string{},
				},
			},
		},
		{
			input: `a + b + c | test[left]; echo hello`,
			hotkeys: []izu.Hotkey{
				{
					Binding: &PartBinding{
						izu.NewDefaultPartList(
							" + ",
							&PartSingle{izu.NewDefaultPartList(" + ", &PartString{"a"})},
							&PartSingle{izu.NewDefaultPartList(" + ", &PartString{"b"})},
							&PartSingle{izu.NewDefaultPartList(" + ", &PartString{"c"})},
						),
					},
					Command: map[string]izu.Part{
						"default": &PartBinding{
							izu.NewDefaultPartList(
								"",
								&PartSingle{izu.NewDefaultPartList("", &PartString{"echo hello"})},
							),
						},
					},
					Flags: map[string][]string{
						"test": {"left"},
					},
				},
			},
		},
		{
			input: `a + b + c | test[right]; abc | echo hello`,
			hotkeys: []izu.Hotkey{
				{
					Binding: &PartBinding{
						izu.NewDefaultPartList(
							" + ",
							&PartSingle{izu.NewDefaultPartList(" + ", &PartString{"a"})},
							&PartSingle{izu.NewDefaultPartList(" + ", &PartString{"b"})},
							&PartSingle{izu.NewDefaultPartList(" + ", &PartString{"c"})},
						),
					},
					Command: map[string]izu.Part{
						"abc": &PartBinding{
							izu.NewDefaultPartList(
								"",
								&PartSingle{izu.NewDefaultPartList("", &PartString{"echo hello"})},
							),
						},
					},
					Flags: map[string][]string{
						"test": {"right"},
					},
				},
			},
		},
		{
			input: `super + XF86Audio{Play,Pause} | test[right]; abc | playerctl {play,pause}`,
			hotkeys: []izu.Hotkey{
				{
					Binding: &PartBinding{
						izu.NewDefaultPartList(
							" + ",
							&PartSingle{izu.NewDefaultPartList(" + ", &PartString{"super"})},
							&PartSingle{izu.NewDefaultPartList(
								" + ",
								&PartString{"XF86Audio"},
								&PartMultiple{
									izu.NewDefaultPartListWithNfixes(
										"{",
										",",
										"}",
										&PartBinding{
											izu.NewDefaultPartList(
												" + ",
												&PartSingle{izu.NewDefaultPartList(" + ", &PartString{"Play"})},
											),
										},
										&PartBinding{
											izu.NewDefaultPartList(
												" + ",
												&PartSingle{izu.NewDefaultPartList(" + ", &PartString{"Pause"})},
											),
										},
									),
								},
							)},
						),
					},
					Command: map[string]izu.Part{
						"abc": &PartBinding{
							izu.NewDefaultPartList(
								"",
								&PartSingle{
									izu.NewDefaultPartList("", &PartString{"playerctl "},
										&PartMultiple{
											izu.NewDefaultPartListWithNfixes(
												"{",
												",",
												"}",
												&PartString{"play"},
												&PartString{"pause"},
											),
										},
									),
								},
							),
						},
					},
					Flags: map[string][]string{
						"test": {"right"},
					},
				},
			},
		},
	}

	for _, c := range cases {
		hotkeys, err := Parse([]byte(c.input))
		if err != nil {
			t.Errorf("'%s' returned error: %v", c.input, err)
			continue
		}

		if len(hotkeys) != len(c.hotkeys) {
			t.Errorf("'%s' returned %d hotkeys, want %d", c.input, len(hotkeys), len(c.hotkeys))
			continue
		}

		for i, hotkey := range hotkeys {
			if hotkey.String() != c.hotkeys[i].String() {
				t.Errorf("'%s' returned hotkey %d: got '%s', want '%s'", c.input, i, hotkey.String(), c.hotkeys[i].String())
			}
		}
	}
}
