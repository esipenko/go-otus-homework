package main

import (
	"bufio"
	"io"
	"net"
	"time"
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
	scanner := bufio.NewScanner(c.in)
	for scanner.Scan() {
		text := scanner.Text()
		_, err := c.connection.Write([]byte(text + "\n"))
		if err != nil {
			return err
		}
	}

	return scanner.Err()
}

func (c *client) Receive() error {
	_, err := io.Copy(c.out, c.connection)
	return err
}

func (c *client) Close() error {
	if c.connection != nil {
		return c.connection.Close()
	}

	return nil
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &client{
		address: address,
		timeout: timeout,
		in:      in,
		out:     out,
	}
}
