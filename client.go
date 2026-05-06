package main

import (
	"bufio"
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"strings"
)

type Client struct {
	nick   string
	user   string
	server string
	port   int
	conn   net.Conn

	currentChannel  string
	currentChannels []string
	channelMembers  map[string][]string

	ignored  map[string]bool
	handlers map[string]func(Message)

	ui UI
}

func NewClient(nick, user, server string, port int, ui UI) *Client {
	client := &Client{nick: nick, user: user, server: server, port: port, ui: ui}
	client.ignored = map[string]bool{}
	client.channelMembers = map[string][]string{}
	client.handlers = map[string]func(Message){
		"PING":    client.handlePing,
		"PRIVMSG": client.handlePrivmsg,
		"JOIN":    client.handleJoin,
		"PART":    client.handlePart,
		"QUIT":    client.handleQuit,
		"NICK":    client.handleNick,
		"353":     client.handleNames,
		"375":     client.handleMOTDStart,
		"372":     client.handleMOTD,
		"376":     client.handleMOTDEnd,
	}
	return client
}

func (client *Client) connect() error {
	var err error
	client.conn, err = tls.Dial("tcp", fmt.Sprintf("%s:%d", client.server, client.port), nil)
	return err
}

func (client *Client) register() error {
	nick := Message{"", "NICK", []string{client.nick}}
	err := client.send(nick)
	if err != nil {
		return err
	}

	user := Message{"", "USER", []string{client.nick, "0", "*", client.user}}
	err = client.send(user)
	if err != nil {
		return err
	}

	return nil
}

func (client *Client) send(msg Message) error {
	_, err := fmt.Fprintf(client.conn, "%s", msg)
	return err
}

func (client *Client) print(format string, args ...any) {
	fmt.Fprintf(client.ui.Chat, format, args...)
}

func (client *Client) readLoop() {
	scanner := bufio.NewScanner(client.conn)

	for scanner.Scan() {
		line := scanner.Text()
		msg, err := parse(line)
		if err != nil {
			client.print("parse error: %s\n", err)
			continue
		}
		client.ui.App.QueueUpdateDraw(func() {
			if handler, ok := client.handlers[msg.command]; ok {
				handler(msg)
			} else {
				client.print("%s\n", line)
			}
		})
	}
}

func (client *Client) refreshNames() {
	client.ui.Members.Clear()
	for _, name := range client.channelMembers[client.currentChannel] {
		client.ui.Members.AddItem(name, "", 0, nil)
	}
}

func (client *Client) parseInput(line string) (Message, error) {
	if line == "" {
		return Message{}, errors.New("parsing an empty string")
	}

	if line[0] != '/' {
		if client.currentChannel == "" {
			return Message{}, errors.New("you are not in a channel")
		}
		return Message{"", "PRIVMSG", []string{client.currentChannel, line}}, nil
	}

	parts := strings.SplitN(line[1:], " ", 2)
	command := parts[0]
	args := ""
	if len(parts) > 1 {
		args = parts[1]
	}

	switch command {
	case "motd":
		return client.cmdMOTD(args)
	case "quit":
		return client.cmdQuit(args)
	case "nick":
		return client.cmdNick(args)
	case "join":
		return client.cmdJoin(args)
	case "msg":
		return client.cmdMsg(args)
	case "part":
		return client.cmdPart(args)
	case "me":
		return client.cmdMe(args)
	case "ignore":
		return client.cmdIgnore(args)
	case "unignore":
		return client.cmdUnignore(args)
	case "ignores":
		return client.cmdIgnores(args)
	default:
		return Message{}, errors.New("unrecognised command")
	}
}
