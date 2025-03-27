package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"time"
)

var (
	ErrClosedByUser   = errors.New("connection closed by user")
	ErrClosedByServer = errors.New("connection closed by server")
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
	log.Println("Connecting to", addr)
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

		if err == nil {
			log.Println(ErrClosedByUser.Error())
			return
		}

		log.Println(err)
	}()

	go func() {
		defer cancel()
		err := client.Receive()
		if err == nil {
			log.Println(ErrClosedByServer.Error())
			return
		}
		log.Fatal(err)
	}()

	<-ctx.Done()
}
