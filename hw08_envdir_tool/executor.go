package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	if len(cmd) == 0 {
		fmt.Println("error: command is empty")
		return 1
	}

	command := cmd[0]
	if !isValidCommand(command) {
		fmt.Printf("error: invalid command: %s\n", cmd[0])
		return 1
	}

	for key, envVal := range env {
		if envVal.NeedRemove {
			err := os.Unsetenv(key)
			if err != nil {
				fmt.Println("remvoe", key, err)
			}
			continue
		}

		err := os.Setenv(key, envVal.Value)
		if err != nil {
			fmt.Println("add", key, envVal.Value, err)
		}
	}

	var args []string
	if len(cmd) > 1 {
		args = cmd[1:]
	}

	c := exec.Command(command, args...)

	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	c.Stdin = os.Stdin
	err := c.Run()
	if err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			return exitErr.ExitCode()
		}

		return 1
	}

	return 0
}

func isValidCommand(command string) bool {
	if strings.TrimSpace(command) == "" || strings.ContainsAny(command, "&|;<>") {
		return false
	}
	return true
}
