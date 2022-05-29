package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadDir(t *testing.T) {
	t.Run("success dir read", func(t *testing.T) {
		environments, err := ReadDir("testdata/env")

		expected := Environment{
			"BAR":   EnvValue{Value: "bar", NeedRemove: false},
			"EMPTY": EnvValue{Value: "", NeedRemove: true},
			"FOO":   EnvValue{Value: "   foo\nwith new line", NeedRemove: false},
			"HELLO": EnvValue{Value: "\"hello\"", NeedRemove: false},
			"UNSET": EnvValue{Value: "", NeedRemove: true},
		}

		require.Equal(t, expected, environments)
		require.Nil(t, err)
	})

	t.Run("fail dir read", func(t *testing.T) {
		environments, err := ReadDir("")

		expected := Environment{}

		require.Equal(t, expected, environments)
		require.NotNil(t, err)
	})
}
