package main

import (
	"errors"
	"strings"
)

func optionalArgs(args string) []string {
	if args != "" {
		return []string{args}
	}
	return []string{}
}

func (client *Client) cmdHelp(args string) (Message, error) {
	printCmd := func(cmd string, def commandDef) {
		if def.args != "" {
			client.printStatus("/%s - %s\n\t%s\n", cmd, def.desc, def.args)
		} else {
			client.printStatus("/%s - %s\n", cmd, def.desc)
		}
	}

	if args != "" {
		def, ok := client.commands[args]
		if !ok {
			return Message{}, errors.New("unknown command")
		}
		printCmd(args, def)
	} else {
		for cmd, def := range client.commands {
			printCmd(cmd, def)
		}
	}
	return Message{}, nil
}

func (client *Client) cmdClear(args string) (Message, error) {
	client.ui.Chat.Clear()
	if ch, ok := client.channels[client.currentChannel]; ok {
		ch.history = nil
	}
	return Message{}, nil
}

func (client *Client) cmdMOTD(args string) (Message, error) {
	return Message{"", "MOTD", optionalArgs(args)}, nil
}

func (client *Client) cmdQuit(args string) (Message, error) {
	return Message{"", "QUIT", optionalArgs(args)}, nil
}

func (client *Client) cmdList(args string) (Message, error) {
	return Message{"", "LIST", optionalArgs(args)}, nil
}

func (client *Client) cmdNick(args string) (Message, error) {
	if args != "" {
		return Message{"", "NICK", []string{args}}, nil
	}
	return Message{}, errors.New("need new nick")
}

func (client *Client) cmdJoin(args string) (Message, error) {
	if args != "" {
		return Message{"", "JOIN", []string{args}}, nil
	}
	return Message{}, errors.New("specify channel to join")
}

func (client *Client) cmdMsg(args string) (Message, error) {
	parts := strings.SplitN(args, " ", 2)
	if len(parts) > 1 {
		return Message{"", "PRIVMSG", []string{parts[0], parts[1]}}, nil
	}
	return Message{}, errors.New("specify nick and message")
}

func (client *Client) cmdPart(args string) (Message, error) {
	channel := client.currentChannel
	partMsg := ""

	if args != "" {
		if args[0] == '#' {
			parts := strings.SplitN(args, " ", 2)
			channel = parts[0]
			if len(parts) > 1 {
				partMsg = parts[1]
			}
		} else {
			partMsg = args
		}
	}

	if channel == "" {
		return Message{}, errors.New("not in a channel")
	}
	if partMsg != "" {
		return Message{"", "PART", []string{channel, partMsg}}, nil
	}
	return Message{"", "PART", []string{channel}}, nil
}

func (client *Client) cmdMe(args string) (Message, error) {
	if args != "" {
		if client.currentChannel == "" {
			return Message{}, errors.New("you are not in a channel")
		}
		return Message{"", "PRIVMSG", []string{client.currentChannel, "\x01ACTION " + args + "\x01"}}, nil
	}
	return Message{}, errors.New("specify action")
}

func (client *Client) cmdIgnore(args string) (Message, error) {
	if args != "" {
		client.ignored[args] = true
		return Message{}, nil
	}
	return Message{}, errors.New("specify nick to add to ignore list")
}

func (client *Client) cmdUnignore(args string) (Message, error) {
	if args != "" {
		delete(client.ignored, args)
		return Message{}, nil
	}
	return Message{}, errors.New("specify nick to remove from ignore list")
}

func (client *Client) cmdIgnores(args string) (Message, error) {
	if len(client.ignored) > 0 {
		for nick := range client.ignored {
			client.printStatus("%s\n", nick)
		}
		return Message{}, nil
	}

	return Message{}, errors.New("no ignored users")
}
