package parser

import (
	"testing"
)

func Test_get_parts(t *testing.T) {
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

		if tc.err && err == nil {
			t.Errorf("case %d: unexpected error: %v", i, err)
		}

		out := path.String()
		if out != tc.out {
			t.Errorf("case %d: expected %q, got %q", i, tc.out, out)
		}
	}
}
