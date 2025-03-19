package main

import (
	"fmt"
	"os"
)

func main() {
	dir, cmd := os.Args[1], os.Args[2:]

	env, err := ReadDir(dir)
	if err != nil {
		fmt.Println(err)
		return
	}

	os.Exit(RunCmd(cmd, env))
}
