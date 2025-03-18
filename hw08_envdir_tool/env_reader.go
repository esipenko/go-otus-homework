package main

import (
	"bufio"
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

func handleEnvFile(filePath string, env Environment) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	fileState, err := file.Stat()
	if err != nil {
		return err
	}

	if fileState.Size() == 0 {
		env[filepath.Base(filePath)] = EnvValue{"", true}
		return nil
	}

	fileScanner := bufio.NewScanner(file)
	fileScanner.Split(bufio.ScanLines)
	fileScanner.Scan()
	envVal := fileScanner.Text()
	envVal = strings.TrimRight(envVal, " \t")
	envVal = strings.ReplaceAll(envVal, "\x00", "\n")

	env[filepath.Base(filePath)] = EnvValue{envVal, false}
	return nil
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

		err := handleEnvFile(filepath.Join(dir, entryName), env)
		if err != nil {
			return nil, err
		}
	}

	return env, nil
}
