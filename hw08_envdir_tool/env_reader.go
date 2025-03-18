package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	env := make(Environment)

	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		entryName := entry.Name()

		if strings.Contains(entryName, "=") {
			continue
		}

		file, err := os.Open(filepath.Join(dir, entryName))
		if err != nil {
			fmt.Println(err)
			return nil, err
		}

		fileState, err := file.Stat()
		if err != nil {
			return nil, err
		}

		if fileState.Size() == 0 {
			env[entryName] = EnvValue{"", true}
			continue
		}

		fileScanner := bufio.NewScanner(file)
		fileScanner.Split(bufio.ScanLines)
		fileScanner.Scan()
		envVal := fileScanner.Text()
		envVal = strings.TrimRight(envVal, " \t")
		envVal = strings.ReplaceAll(envVal, "\x00", "\n")

		env[entryName] = EnvValue{envVal, false}
		file.Close()
	}

	return env, nil
}
