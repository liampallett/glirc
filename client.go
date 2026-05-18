package main

import (
	"errors"
	"fmt"
	"strings"
)

type commandDef struct {
	desc string
	args string
	fn   func(string) (Message, error)
}

type Client struct {
	nick   string
	user   string
	server *Server

	currentChannel string
	channels       map[string]*Channel

	ignored  map[string]bool
	handlers map[string]func(Message)
	commands map[string]commandDef

	ui UI
}

func NewClient(nick, user, address string, port int, ui UI) *Client {
	server := &Server{address: address, port: port, incoming: make(chan Message)}
	client := &Client{nick: nick, user: user, server: server, ui: ui}

	go func() {
		for msg := range client.server.incoming {
			client.ui.App.QueueUpdateDraw(func() {
				if client.ignored[msg.Nick()] {
					return
				}
				if handler, ok := client.handlers[msg.command]; ok {
					handler(msg)
				} else {
					client.printStatus("%s\n", msg)
				}
			})
		}
		client.ui.App.QueueUpdateDraw(func() {
			client.printStatus("disconnected from server")
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
	client.commands = map[string]commandDef{
		"help":     {"display available commands", "optional command name", client.cmdHelp},
		"clear":    {"clear the chat window", "", client.cmdClear},
		"motd":     {"display current server message of the day", "", client.cmdMOTD},
		"list":     {"lists all channels and their topics", "filter with >, <", client.cmdList},
		"quit":     {"quit the application", "optional quit message", client.cmdQuit},
		"nick":     {"change your nickname displayed on the server", "new nickname", client.cmdNick},
		"join":     {"join the specified channel", "#channel name", client.cmdJoin},
		"msg":      {"privately message a user on the server", "username, message", client.cmdMsg},
		"part":     {"leave a channel", "channel to leave (default: current), optional parting message", client.cmdPart},
		"me":       {"send a message from yourself", "message", client.cmdMe},
		"ignore":   {"add a user to your ignore list (will not see messages, join, part, quit, etc.)", "user nick", client.cmdIgnore},
		"unignore": {"remove a user from your ignore list", "user nick", client.cmdUnignore},
		"ignores":  {"display ignore list", "", client.cmdIgnores},
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
	return client.server.send(user)
}

func (client *Client) printStatus(format string, args ...any) {
	_, err := fmt.Fprintf(client.ui.Chat, format, args...)
	if err != nil {
		return
	}
}

func (client *Client) anyChannel() string {
	for name := range client.channels {
		return name
	}
	return ""
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
	client.ui.Members.SetTitle(memberCountTitle(memberCount))
}

func (client *Client) switchChannel(name string) {
	client.currentChannel = name
	client.ui.Chat.SetTitle(name)
	client.ui.Chat.Clear()
	client.refreshNames()
	if ch, ok := client.channels[name]; ok {
		fmt.Fprintf(client.ui.Chat, strings.Join(ch.history, ""))
	}
}

func memberCountTitle(n int) string {
	if n < 1 {
		return "Members"
	}
	return fmt.Sprintf("Members - %d", n)
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

	if def, ok := client.commands[command]; ok {
		return def.fn(args)
	}
	return Message{}, errors.New("unrecognised command (see /help)")
}
