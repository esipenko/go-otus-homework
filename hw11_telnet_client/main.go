package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"time"
)

var timeout string

func init() {
	flag.StringVar(&timeout, "timeout", "10s", "file to read from")
}

func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) < 2 {
		return
	}

	duration, err := time.ParseDuration(timeout)
	if err != nil {
		fmt.Println("Error parsing duration:", err)
		return
	}

	addr := strings.Join(args, ":")

	client := NewTelnetClient(addr, duration, os.Stdin, os.Stdout)
	err = client.Connect()
	if err != nil {
		fmt.Println("Error connecting:", err)
		return
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)

	go func() {
		defer cancel()
		err := client.Send()
		if err != nil {
			log.Println("Error sending:", err)
			return
		}
	}()

	go func() {
		defer cancel()

		err := client.Receive()
		if err != nil {
			log.Println("Error sending:", err)
			return
		}
	}()

	<-ctx.Done()
}
