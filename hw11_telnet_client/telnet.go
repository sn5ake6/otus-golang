package main

import (
	"bufio"
	"io"
	"net"
	"time"
)

type TelnetClienter interface {
	Connect() error
	Close() error
	Send() error
	Receive() error
}

type TelnetClient struct {
	address string
	timeout time.Duration
	in      io.Reader
	out     io.Writer
	conn    net.Conn
}

func (cl *TelnetClient) Connect() error {
	conn, err := net.DialTimeout("tcp", cl.address, cl.timeout)
	if err != nil {
		return err
	}

	cl.conn = conn

	return nil
}

func (cl *TelnetClient) Receive() error {
	return cl.transmit(cl.conn, cl.out)
}

func (cl *TelnetClient) Send() error {
	return cl.transmit(cl.in, cl.conn)
}

func (cl *TelnetClient) transmit(from io.Reader, to io.Writer) error {
	scanner := bufio.NewScanner(from)
	for scanner.Scan() {
		_, err := to.Write(append(scanner.Bytes(), '\n'))
		if err != nil {
			return err
		}
	}

	return scanner.Err()
}

func (cl *TelnetClient) Close() error {
	if cl.conn != nil {
		return cl.conn.Close()
	}

	return nil
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClienter {
	return &TelnetClient{
		address: address,
		timeout: timeout,
		in:      in,
		out:     out,
	}
}
