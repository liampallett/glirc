package main

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"net"
)

type Server struct {
	address  string
	port     int
	conn     net.Conn
	incoming chan Message
	err      error
}

func (server *Server) connect() error {
	conn, err := tls.Dial("tcp", fmt.Sprintf("%s:%d", server.address, server.port), nil)
	if err == nil {
		server.conn = conn
	}
	return err
}

func (server *Server) send(msg Message) error {
	_, err := fmt.Fprint(server.conn, msg)
	return err
}

func (server *Server) readLoop() {
	scanner := bufio.NewScanner(server.conn)

	for scanner.Scan() {
		line := scanner.Text()
		msg, err := parse(line)
		if err != nil {
			continue
		}
		server.incoming <- msg
	}
	server.err = scanner.Err()

	close(server.incoming)
}
