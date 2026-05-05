package main

import (
	"bufio"
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"os"
	"strings"
)

type Client struct {
	nick           string
	user           string
	server         string
	port           int
	conn           net.Conn
	currentChannel string
}

func (client *Client) connect() error {
	var err error
	client.conn, err = tls.Dial("tcp", fmt.Sprintf("%s:%d", client.server, client.port), nil)
	return err
}

func (client *Client) register() {
	nick := Message{"", "NICK", []string{client.nick}}
	fmt.Fprintf(client.conn, "%s", nick)

	user := Message{"", "USER", []string{client.nick, "0", "*", client.user}}
	fmt.Fprintf(client.conn, "%s", user)
}

func (client *Client) run() {
	buffServer := make(chan string)
	go func() {
		scanner := bufio.NewScanner(client.conn)
		for scanner.Scan() {
			msg := scanner.Text()
			buffServer <- msg
		}
		close(buffServer)
	}()

	buffClient := make(chan string)
	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			msg := scanner.Text()
			buffClient <- msg
		}
	}()

	for {
		select {
		case line, ok := <-buffServer:
			if !ok {
				return
			}
			msg := parse(line)
			switch msg.command {
			case "PING":
				pong := Message{"", "PONG", msg.parameters}
				fmt.Fprintf(client.conn, "%s", pong)
			case "PRIVMSG":
				nick := msg.Nick()
				text := msg.parameters[1]
				fmt.Printf("<%s> %s\n", nick, text)
			case "JOIN":
				nick := msg.Nick()
				text := msg.parameters[0]
				if nick == client.nick {
					fmt.Printf("you joined %s\n", text)
				} else {
					fmt.Printf("%s joined %s\n", nick, text)
				}
			case "PART":
				nick := msg.Nick()
				text := msg.parameters[0]
				if nick == client.nick {
					fmt.Printf("you left %s\n", text)
				} else {
					fmt.Printf("%s left %s\n", nick, text)
				}
			case "QUIT":
				nick := msg.Nick()
				if len(msg.parameters) < 1 {
					if nick == client.nick {
						fmt.Printf("you quit\n")
					} else {
						fmt.Printf("%s quit\n", nick)
					}
				} else {
					text := msg.parameters[0]
					quitReason := strings.TrimPrefix(text, "Quit: ")
					if nick == client.nick {
						fmt.Printf("you quit: %s\n", quitReason)
					} else {
						fmt.Printf("%s quit: %s\n", nick, quitReason)
					}
				}
			default:
				fmt.Println(line)
			}
		case line := <-buffClient:
			msg, err := client.parseInput(line)
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Fprintf(client.conn, "%s", msg)
				if msg.command == "PRIVMSG" {
					fmt.Printf("<%s> %s\n", client.nick, msg.parameters[1])
				}
			}
		}
	}
}

func (client *Client) parseInput(line string) (Message, error) {
	var msg Message

	if line[0] != '/' {
		if client.currentChannel == "" {
			return Message{}, errors.New("you are not in a channel")
		}
		msg = Message{"", "PRIVMSG", []string{client.currentChannel, line}}
	} else {
		lineParts := strings.SplitN(line[1:], " ", 2)
		rawCommand := lineParts[0]
		switch rawCommand {
		case "join":
			if len(lineParts) < 2 {
				return Message{}, errors.New("specify channel to join")
			}
			client.currentChannel = lineParts[1]
			msg = Message{"", "JOIN", []string{client.currentChannel}}
		case "part":
			if client.currentChannel == "" {
				return Message{}, errors.New("trying to part from no channel")
			}
			msg = Message{"", "PART", []string{client.currentChannel}}
			client.currentChannel = ""
		case "quit":
			if len(lineParts) > 1 {
				msg = Message{"", "QUIT", []string{lineParts[1]}}
			} else {
				msg = Message{"", "QUIT", []string{}}
			}
		case "msg":
			if len(lineParts) < 2 {
				return Message{}, errors.New("need target nick and message")
			}
			messageParts := strings.SplitN(lineParts[1], " ", 2)
			if len(messageParts) < 2 {
				return Message{}, errors.New("need target nick and message")
			}
			msg = Message{"", "PRIVMSG", []string{messageParts[0], messageParts[1]}}
		default:
			return Message{}, errors.New("unrecognised command")
		}
	}

	return msg, nil
}
