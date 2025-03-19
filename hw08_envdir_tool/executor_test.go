package main

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunCmd(t *testing.T) {
	t.Run("command executes successfully", func(t *testing.T) {
		cmd := []string{"ls", "-la"}

		returnCode := RunCmd(cmd, make(Environment))
		require.Equal(t, 0, returnCode)
	})

	t.Run("command fails with exit code", func(t *testing.T) {
		env := Environment{}
		cmd := []string{"sh", "-c", "exit 2"}

		returnCode := RunCmd(cmd, env)
		require.Equal(t, 2, returnCode)
	})

	t.Run("environment variable is set", func(t *testing.T) {
		env := Environment{
			"FOO": EnvValue{Value: "bar", NeedRemove: false},
		}
		cmd := []string{"sh", "-c", "echo $FOO"}

		r, w, _ := os.Pipe()
		oldStdout := os.Stdout
		os.Stdout = w

		returnCode := RunCmd(cmd, env)
		require.Equal(t, 0, returnCode)

		w.Close()
		var buf bytes.Buffer
		io.Copy(&buf, r)
		os.Stdout = oldStdout

		require.Contains(t, buf.String(), "bar")
	})

	t.Run("environment variable is unset", func(t *testing.T) {
		os.Setenv("FOO", "bar")
		env := Environment{
			"FOO": EnvValue{NeedRemove: true},
		}
		cmd := []string{"sh", "-c", "echo $FOO"}

		r, w, _ := os.Pipe()
		oldStdout := os.Stdout
		os.Stdout = w

		returnCode := RunCmd(cmd, env)
		require.Equal(t, 0, returnCode)

		w.Close()
		var buf bytes.Buffer
		io.Copy(&buf, r)
		os.Stdout = oldStdout
		require.Equal(t, buf.String(), "\n")
	})

	t.Run("empty command", func(t *testing.T) {
		cmd := []string{}

		returnCode := RunCmd(cmd, make(Environment))
		require.Equal(t, 1, returnCode)
	})

	t.Run("invalid command", func(t *testing.T) {
		cmd := []string{"        "}

		returnCode := RunCmd(cmd, make(Environment))
		require.Equal(t, 1, returnCode)
	})
}
