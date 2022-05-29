package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunCmd(t *testing.T) {
	t.Run("run with exit code 0", func(t *testing.T) {
		cmd := []string{
			"/bin/bash",
			"testdata/echo.sh",
		}
		returnCode := RunCmd(cmd, Environment{})

		require.Equal(t, 0, returnCode)
	})

	t.Run("run with exit code 2", func(t *testing.T) {
		cmd := []string{
			"/bin/bash",
			"testdata/exit2.sh",
		}
		returnCode := RunCmd(cmd, Environment{})

		require.Equal(t, 2, returnCode)
	})

	t.Run("run with not existed file", func(t *testing.T) {
		cmd := []string{
			"/bin/bash",
			"testdata/echo_not_existed.sh",
		}
		returnCode := RunCmd(cmd, Environment{})

		require.Equal(t, 127, returnCode)
	})

	t.Run("run without args", func(t *testing.T) {
		cmd := []string{
			"/bin/bash",
		}
		returnCode := RunCmd(cmd, Environment{})

		require.Equal(t, 0, returnCode)
	})
}
