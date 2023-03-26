package namer

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestInvalidLabelValueRegexp(t *testing.T) {
	var tests = []struct {
		desc   string
		format string
		cmp    string
	}{
		{desc: "empty string", format: "", cmp: ""},
		{desc: "single white-space", format: " ", cmp: "_"},
		{desc: "double white-space", format: "  ", cmp: "__"},
		{desc: "trible white-space", format: "   ", cmp: "___"},
		{desc: "single dot", format: ".", cmp: "_"},
		{desc: "double dot", format: "..", cmp: "__"},
		{desc: "trible dot", format: "...", cmp: "___"},
		{desc: "letter 's'", format: "s", cmp: "s"},
		{desc: "word 'apple'", format: "apple", cmp: "apple"},
		{desc: "sentence 'You are handsome.'", format: "You are handsome.", cmp: "You_are_handsome_"},
		{desc: "sentence 'My name: Satoshi Nakamoto.'", format: "My name: Satoshi Nakamoto.", cmp: "My_name__Satoshi_Nakamoto_"},
		{desc: "single vertical bar", format: "|", cmp: "_"},
		{desc: "double vertical bar", format: "||", cmp: "__"},
		{desc: "vertical bar and white-space", format: "| |", cmp: "___"},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			value := invalidLabelValueRegexp.ReplaceAllString(test.format, "_")
			require.Equal(t, value, test.cmp)
		})
	}
}

func TestFormatRegexp(t *testing.T) {
	var tests = []struct {
		desc   string
		format string
		cmp    []string
	}{
		{
			desc:   "empty string",
			format: "",
			cmp:    []string{},
		},
		{
			desc:   "%{name}",
			format: "%{name}",
			cmp: []string{
				"%{name}",
				"name",
			},
		},
		{
			desc: "%{#name}",
			format: "%{#name}",
			cmp: []string{
				"%{#name}",
				"#name",
			},
		},
		{
			desc: "%{##name}",
			format: "%{##name}",
			cmp: []string{
				"%{##name}",
				"##name",
			},
		},
		{
			desc: "1234%{###}",
			format: "1234%{###}",
			cmp: []string{
				"%{###}",
				"###",
			},
		},
		{
			desc: "1234%{123}",
			format: "1234%{123}",
			cmp: []string{
				"%{123}",
				"123",
			},
		},
		{
			desc: "%{123abc",
			format: "%{123abc",
			cmp: []string{},
		},
		{
			desc: "name_%{}",
			format: "name_%{}",
			cmp: []string{
				"%{}",
				"",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			matches := formatRegexp.FindAllStringSubmatchIndex(test.format, -1)
			for i, m := range matches {
				require.Equal(t, test.cmp[i*2], test.format[m[0]:m[1]])
				require.Equal(t, test.cmp[i*2+1], test.format[m[2]:m[3]])
			}
		})
	}
}
