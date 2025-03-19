package main

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadDir(t *testing.T) {
	t.Run("= in file name not allowed, this file will not add to Env", func(t *testing.T) {
		err := os.Mkdir("testenv", os.ModeDir|os.ModePerm)
		if err != nil {
			fmt.Println(err)
			return
		}

		file, err := os.Create("testenv/FOO=BAR")
		if err != nil {
			fmt.Println(err)
			return
		}

		env, _ := ReadDir("testenv")
		require.Empty(t, env)
		file.Close()
		os.RemoveAll("testenv")
	})

	t.Run("empty file will be added to Env with ShouldRemove true", func(t *testing.T) {
		err := os.Mkdir("testenv", os.ModeDir|os.ModePerm)
		if err != nil {
			fmt.Println(err)
			return
		}

		file, err := os.Create("testenv/FOO")
		if err != nil {
			fmt.Println(err)
			return
		}

		env, _ := ReadDir("testenv")

		require.NotEmpty(t, env)
		require.True(t, env["FOO"].NeedRemove)

		defer file.Close()
		defer os.RemoveAll("testenv")
	})

	t.Run("terminal zeroes shold be replaces with \\n", func(t *testing.T) {
		err := os.Mkdir("testenv", os.ModeDir|os.ModePerm)
		if err != nil {
			return
		}

		file, err := os.Create("testenv/FOO")
		if err != nil {
			fmt.Println(err)
			return
		}

		_, _ = file.WriteString("FOO\x00BAR")

		env, _ := ReadDir("testenv")

		require.NotEmpty(t, env)
		require.False(t, env["FOO"].NeedRemove)
		require.Equal(t, env["FOO"].Value, "FOO\nBAR")

		defer file.Close()
		defer os.RemoveAll("testenv")
	})

	t.Run("spaces and tabulations in the end of string should be trimmed", func(t *testing.T) {
		err := os.Mkdir("testenv", os.ModeDir|os.ModePerm)
		if err != nil {
			fmt.Println(err)
			return
		}

		file, err := os.Create("testenv/FOO")
		if err != nil {
			fmt.Println(err)
			return
		}

		file.WriteString("\t  \tFOO\t      \t")

		env, _ := ReadDir("testenv")
		require.NotEmpty(t, env)
		require.False(t, env["FOO"].NeedRemove)
		require.Equal(t, env["FOO"].Value, "\t  \tFOO")

		defer file.Close()
		defer os.RemoveAll("testenv")
	})

	t.Run("default case", func(t *testing.T) {
		env, _ := ReadDir("testdata/env/")
		testEnv := make(Environment)
		testEnv["BAR"] = EnvValue{"bar", false}
		testEnv["HELLO"] = EnvValue{"\"hello\"", false}
		testEnv["UNSET"] = EnvValue{"", true}
		testEnv["EMPTY"] = EnvValue{"", false}
		testEnv["FOO"] = EnvValue{"   foo\nwith new line", false}
		require.Exactly(t, testEnv, env)
	})
}
