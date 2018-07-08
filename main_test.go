package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestArgsWithoutKnownFlags(t *testing.T) {
	a := assert.New(t)
	for _, c := range []struct {
		name     string
		args     []string
		expected []string
	}{
		{name: "if with no args, return no args"},
		{
			name:     "if with single `-` prefix, omit validly",
			args:     []string{"-nocolor", "hoge"},
			expected: []string{"hoge"},
		},
		{
			name:     "if with double `-` prefix, omit validly",
			args:     []string{"--nocolor", "hoge"},
			expected: []string{"hoge"},
		},
		{
			name: "if with many options, omit validly",
			args: []string{
				"--hoge",
				"-nocolor",
				"--selfupdate",
				"--version",
				"-v",
				"--fuga",
				"hogefuga",
			},
			expected: []string{
				"--hoge",
				"--fuga",
				"hogefuga",
			},
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			defer setOsArgs(c.args)()
			a.Equal(c.expected, argsWithoutKnownFlags())
		})
	}
}

func setOsArgs(args []string) func() {
	original := osArgs
	osArgs = os.Args[0:1]
	osArgs = append(osArgs, args...)
	return func() { osArgs = original }
}
