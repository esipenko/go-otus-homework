package main

import (
	"fmt"

	"golang.org/x/example/hello/reverse"
)

const str = "Hello, OTUS!"

func main() {
	fmt.Println(reverse.String(str))
}
