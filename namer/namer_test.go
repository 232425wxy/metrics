package namer

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFullyQualifiedName(t *testing.T) {
	var tests = []struct {
		name   Namer
		fqname string
	}{
		{
			name: Namer{
				namespace: "namespace",
				subsystem: "subsystem",
				name:      "name",
			},
			fqname: "namespace.subsystem.name",
		},
		{
			name: Namer{
				namespace: "namespace",
				name:      "name",
			},
			fqname: "namespace.name",
		},
		{
			name: Namer{
				subsystem: "subsystem",
				name:      "name",
			},
			fqname: "subsystem.name",
		},
		{
			name: Namer{
				name: "name",
			},
			fqname: "name",
		},
		{
			name:   Namer{},
			fqname: "",
		},
	}

	for _, test := range tests {
		t.Run(test.fqname, func(t *testing.T) {
			require.Equal(t, test.fqname, test.name.FullyQualifiedName())
		})
	}
}

func TestFormat(t *testing.T) {
	var tests = []struct {
		namer       *Namer
		labelValues []string
		expected    string
	}{
		{
			namer: &Namer{
				namespace:  "namespace",
				subsystem:  "subsystem",
				name:       "name",
				nameFormat: "%{#namespace}%{#name}",
				labelNames: map[string]struct{}{},
			},
			labelValues: []string{},
			expected:    "namespacename",
		},
		{
			namer: &Namer{
				namespace:  "namespace",
				subsystem:  "subsystem",
				name:       "name",
				nameFormat: "%{#namespace}%{name}",
				labelNames: map[string]struct{}{"name": {}},
			},
			labelValues: []string{"name"},
			expected:    "namespaceunknown",
		},
		{
			namer: &Namer{
				namespace:  "namespace",
				subsystem:  "subsystem",
				name:       "name",
				nameFormat: "%{#namespace}%{p2p}",
				labelNames: map[string]struct{}{"p2p": {}},
			},
			labelValues: []string{"p2p", "gossip"},
			expected:    "namespacegossip",
		},
	}

	for _, test := range tests {
		res := test.namer.Format(test.labelValues...)
		require.Equal(t, test.expected, res)
	}
}
