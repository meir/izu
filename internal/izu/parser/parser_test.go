package parser

import (
	"testing"
)

func Test_base(t *testing.T) {
	cases := []struct {
		in  string
		out string
		err bool
	}{
		{"", "", false},
		{"a", "a", false},
		{"a + b", "a + b", false},
		{"a + b + c", "a + b + c", false},
		{"a + {a, b} + d", "a + { a, b } + d", false},
		{"a + {_, b + e + } d", "a + { , b + e } + d", false},
		{"a + {_, b +} d", "a + { , b } + d", false},
		{"XF68Media{Play,Pause}", "XF68Media{Play,Pause}", false},
		{"a + }", "a", false},
	}

	for i, tc := range cases {
		path := NewBase()
		_, err := path.Parse([]byte(tc.in))

		if tc.err == (err == nil) {
			t.Errorf("case %d: unexpected error: %v", i, err)
		}

		out := path.String()
		if out != tc.out {
			t.Errorf("case %d: expected %q, got %q", i, tc.out, out)
		}
	}
}

func Test_keybind(t *testing.T) {
	cases := []struct {
		in  string
		out string
		err bool
	}{
		{"", "", true},
		{"a\necho", "a\n\techo", false},
		{"a + b\necho", "a + b\n\techo", false},
		{"a + b + c\necho", "a + b + c\n\techo", false},
		{"a + {a, b} + d\necho {a, b}", "a + { a, b } + d\n\techo {a, b}", false},
		{"a + {_, b + e + } d\n echo {a, b}", "a + { , b + e } + d\n\techo {a, b}", false},
		{"a + {_, b +} d\n echo {a, b}", "a + { , b } + d\n\techo {a, b}", false},
		{"XF68Media{Play,Pause}\n echo {a, b}", "XF68Media{Play,Pause}\n\techo {a, b}", false},
		{"a + }\necho", "", true},
		{"a + }", "", true},
		{"a", "", true},
	}

	for i, tc := range cases {
		path := NewKeybind()
		_, err := path.Parse([]byte(tc.in))

		if tc.err == (err == nil) {
			t.Errorf("case %d: unexpected error: %v", i, err)
		}

		out := path.String()
		if out != tc.out {
			t.Errorf("case %d: expected %q, got %q", i, tc.out, out)
		}
	}
}
