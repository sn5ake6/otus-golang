package main

import (
	"errors"
	"os"
	"os/exec"
)

var exitError *exec.ExitError

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	command := cmd[0]
	args := cmd[1:]

	cmdForExec := exec.Command(command, args...)
	for name, value := range env {
		if value.NeedRemove {
			os.Unsetenv(name)
			continue
		}
		os.Setenv(name, value.Value)
	}

	cmdForExec.Env = os.Environ()
	cmdForExec.Stdin = os.Stdin
	cmdForExec.Stdout = os.Stdout
	cmdForExec.Stderr = os.Stderr

	if err := cmdForExec.Run(); err != nil {
		if errors.As(err, &exitError) {
			return exitError.ExitCode()
		}
	}

	return 0
}
