package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"
	"time"
)

var (
	ErrClosedByUser   = errors.New("connection closed by user")
	ErrClosedByServer = errors.New("connection closed by server")
)

type TelnetClient interface {
	Connect() error
	io.Closer
	Send() error
	Receive() error
}

type client struct {
	address    string
	connection net.Conn
	timeout    time.Duration
	in         io.ReadCloser
	out        io.Writer
}

func (c *client) Connect() error {
	conn, err := net.DialTimeout("tcp", c.address, c.timeout)
	if err != nil {
		return err
	}

	c.connection = conn
	return nil
}

func (c *client) Send() error {
	reader := bufio.NewReader(c.in)
	for {
		text, err := reader.ReadString('\n')
		if err != nil {
			if errors.Is(err, io.EOF) {
				return ErrClosedByUser
			}

			return err
		}

		_, err = c.connection.Write([]byte(text))
		if err != nil {
			fmt.Println("Error writing:", err)
			return err
		}
	}
}

func (c *client) Receive() error {
	reader := bufio.NewReader(c.connection)
	for {
		text, err := reader.ReadString('\n')
		if err != nil {
			if errors.Is(err, io.EOF) {
				return ErrClosedByServer
			}
			return err
		}

		_, err = c.out.Write([]byte(text))
		if err != nil {
			return err
		}
	}
}

func (c *client) Close() error { return c.connection.Close() }
func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &client{
		address: address, in: in, out: out, timeout: timeout,
	}
}
