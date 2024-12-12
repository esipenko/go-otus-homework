package hw02unpackstring

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUnpack(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{input: "a4bc2d5e", expected: "aaaabccddddde"},
		{input: "abccd", expected: "abccd"},
		{input: "", expected: ""},
		{input: "aaa0b", expected: "aab"},
		{input: "aaa0bbb0", expected: "aabb"},

		{input: "d\n5abc", expected: "d\n\n\n\n\nabc"},
		{input: "🇷🇺3🏴‍☠️🏴‍☠️0", expected: "🇷🇺🇷🇺🇷🇺🏴‍☠️"},

		// uncomment if task with asterisk completed
		{input: `qwe\4\5`, expected: `qwe45`},
		{input: `qwe\45`, expected: `qwe44444`},
		{input: `qwe\\5`, expected: `qwe\\\\\`},
		{input: `qwe\\\3`, expected: `qwe\3`},
		{input: `\\3`, expected: `\\\`},
		{input: `\\0`, expected: ``},
		{input: `\\1`, expected: `\`},
		{input: `гоо3`, expected: `гоооо`},
		{input: "a\\\\4bc2d5e", expected: "a\\\\\\\\bccddddde"},
		{input: "a\\4bc2d5e", expected: "a4bccddddde"},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.input, func(t *testing.T) {
			result, err := Unpack(tc.input)
			require.NoError(t, err)
			require.Equal(t, tc.expected, result)
		})
	}
}

func TestUnpackInvalidString(t *testing.T) {
	invalidStrings := []string{"3abc", "45", "aaa10b", "d\\ж5abc", `qw\ne`, `\🇷🇺`}
	for _, tc := range invalidStrings {
		tc := tc
		t.Run(tc, func(t *testing.T) {
			_, err := Unpack(tc)
			require.Truef(t, errors.Is(err, ErrInvalidString), "actual error %q", err)
		})
	}
}
