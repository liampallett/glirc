package main

import (
	"errors"
	"fmt"
	"strings"
)

type Client struct {
	nick   string
	user   string
	server *Server

	currentChannel string
	channels       map[string]*Channel

	ignored  map[string]bool
	handlers map[string]func(Message)

	ui UI
}

func NewClient(nick, user, address string, port int, ui UI) *Client {
	server := &Server{address: address, port: port, incoming: make(chan Message)}
	client := &Client{nick: nick, user: user, server: server, ui: ui}

	go func() {
		for msg := range client.server.incoming {
			client.ui.App.QueueUpdateDraw(func() {
				msg := msg
				if handler, ok := client.handlers[msg.command]; ok {
					handler(msg)
				} else {
					client.print("%s\n", msg)
				}
			})
		}
		client.ui.App.QueueUpdateDraw(func() {
			client.print("disconnected from server")
		})
	}()

	client.channels = map[string]*Channel{}
	client.ignored = map[string]bool{}
	client.handlers = map[string]func(Message){
		"NOTICE":  client.handleNotice,
		"PING":    client.handlePing,
		"PRIVMSG": client.handlePrivmsg,
		"JOIN":    client.handleJoin,
		"PART":    client.handlePart,
		"QUIT":    client.handleQuit,
		"NICK":    client.handleNick,
		"321":     client.handleListStart,
		"322":     client.handleList,
		"323":     client.handleListEnd,
		"353":     client.handleNames,
		"366":     client.handleNamesEnd,
		"375":     client.handleMOTDStart,
		"372":     client.handleMOTD,
		"376":     client.handleMOTDEnd,
	}
	return client
}

func (client *Client) register() error {
	nick := Message{"", "NICK", []string{client.nick}}
	err := client.server.send(nick)
	if err != nil {
		return err
	}

	user := Message{"", "USER", []string{client.nick, "0", "*", client.user}}
	err = client.server.send(user)
	if err != nil {
		return err
	}

	return nil
}

func (client *Client) print(format string, args ...any) {
	_, err := fmt.Fprintf(client.ui.Chat, format, args...)
	if err != nil {
		return
	}
}

func (client *Client) printChannel(channel string, format string, args ...any) {
	text := fmt.Sprintf(format, args...)
	ch, ok := client.channels[channel]
	if !ok {
		return
	}
	ch.history = append(ch.history, text)
	if channel == client.currentChannel {
		_, err := fmt.Fprint(client.ui.Chat, text)
		if err != nil {
			return
		}
	}
}

func (client *Client) refreshNames() {
	client.ui.Members.Clear()
	ch, ok := client.channels[client.currentChannel]
	if !ok {
		return
	}
	for _, name := range ch.members {
		client.ui.Members.AddItem(name, "", 0, nil)
	}
	memberCount := len(ch.members)
	if memberCount < 1 {
		client.ui.Members.SetTitle("Members")
	} else {
		client.ui.Members.SetTitle(fmt.Sprintf("Members - %d", memberCount))
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
	case "help":
		return client.cmdHelp(args)
	case "clear":
		return client.cmdClear(args)
	case "motd":
		return client.cmdMOTD(args)
	case "list":
		return client.cmdList(args)
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
		return Message{}, errors.New("unrecognised command (see /help)")
	}
}
