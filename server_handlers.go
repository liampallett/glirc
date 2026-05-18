package main

import (
	"slices"
	"strings"
)

func (client *Client) handleMOTDStart(msg Message) {
	text, ok := msg.param(1)
	if !ok {
		return
	}
	client.printStatus("%s\n", text)
}

func (client *Client) handleMOTD(msg Message) {
	text, ok := msg.param(1)
	if !ok || len(text) < 2 {
		return
	}
	client.printStatus("%s\t\n", text[2:])
}

func (client *Client) handleMOTDEnd(msg Message) {
	client.printStatus("\n")
}

func (client *Client) handleNotice(msg Message) {
	text, ok := msg.param(1)
	if !ok {
		return
	}
	client.printStatus("%s %s\n", msg.parameters[0], text)
}

func (client *Client) handleListStart(msg Message) {
	client.printStatus("\n")
}

func (client *Client) handleList(msg Message) {
	if len(msg.parameters) < 3 {
		return
	}
	channel := msg.parameters[1]
	userCount := msg.parameters[2]
	topic := msg.parameters[len(msg.parameters)-1]

	if topic != "" {
		client.printStatus("- %s: %s current users - %s\n", channel, userCount, topic)
	} else {
		client.printStatus("- %s: %s current users\n", channel, userCount)
	}
}

func (client *Client) handleListEnd(msg Message) {
	client.printStatus("\n")
}

func (client *Client) handleNames(msg Message) {
	if len(msg.parameters) < 3 {
		return
	}
	channel := msg.parameters[2]
	members := strings.Fields(msg.parameters[len(msg.parameters)-1])
	ch, ok := client.channels[channel]
	if !ok {
		return
	}
	ch.members = nil
	for _, member := range members {
		ch.addMember(member)
	}
	client.refreshNames()
}

func (client *Client) handleNamesEnd(msg Message) {
	client.printStatus("\n")
}

func (client *Client) handlePing(msg Message) {
	pong := Message{"", "PONG", msg.parameters}
	err := client.server.send(pong)
	if err != nil {
		client.printStatus("%s\n", err)
	}
}

func (client *Client) handlePrivmsg(msg Message) {
	text, ok := msg.param(1)
	if !ok {
		return
	}

	nick := msg.Nick()
	target := msg.parameters[0]
	switch target[0] {
	case '#':
		if action, ok := parseAction(text); ok {
			client.printChannel(target, "* %s %s\n", nick, action)
		} else {
			client.printChannel(target, "<%s> %s\n", nick, text)
		}
	default:

		key := nick
		if nick == client.nick {
			key = target
		}
		if _, ok := client.channels[key]; !ok {
			client.channels[key] = &Channel{name: key}
			client.ui.Channels.AddItem(key, "", 0, nil)
		}

		if action, ok := parseAction(text); ok {
			client.printChannel(target, "* %s %s\n", nick, action)
		} else {
			client.printChannel(target, "<%s> %s\n", nick, text)
		}
	}

}

func (client *Client) handleJoin(msg Message) {
	nick := msg.Nick()
	channel := msg.parameters[0]
	if nick == client.nick {
		client.printStatus("you joined %s\n", channel)
		client.switchChannel(channel)
		if _, ok := client.channels[channel]; !ok {
			client.channels[channel] = &Channel{name: channel}
			client.ui.Channels.AddItem(channel, "", 0, nil)
		}
	} else {
		client.printChannel(channel, "%s joined %s\n", nick, channel)
		ch, ok := client.channels[channel]
		if !ok {
			return
		}
		ch.addMember(nick)
	}
	client.refreshNames()
}

func (client *Client) handlePart(msg Message) {
	nick := msg.Nick()
	channel := msg.parameters[0]
	if nick == client.nick {
		client.printStatus("you left %s\n", channel)
		delete(client.channels, channel)
		client.currentChannel = ""
		indices := client.ui.Channels.FindItems(channel, "", false, true)
		if len(indices) > 0 {
			client.ui.Channels.RemoveItem(indices[0])
		}
		client.currentChannel = client.anyChannel()
		client.switchChannel(client.currentChannel)
	} else {
		client.printChannel(channel, "%s left %s\n", nick, channel)
		ch, ok := client.channels[channel]
		if !ok {
			return
		}
		ch.removeMember(nick)
	}
	client.refreshNames()
}

func (client *Client) handleQuit(msg Message) {
	nick := msg.Nick()
	text, ok := msg.param(0)
	if !ok {
		if nick == client.nick {
			client.printStatus("you quit\n")
		} else {
			client.printChannel(client.currentChannel, "%s quit\n", nick)
		}
	} else {
		quitReason := strings.TrimPrefix(text, "Quit: ")
		if nick == client.nick {
			client.printStatus("you quit: %s\n", quitReason)
		} else {
			for _, channel := range client.channels {
				if slices.Contains(channel.members, nick) {
					client.printChannel(channel.name, "%s quit: %s\n", nick, quitReason)
					channel.removeMember(nick)
				}
			}
		}
	}
	client.refreshNames()
}

func (client *Client) handleNick(msg Message) {
	newNick, ok := msg.param(0)
	if !ok {
		return
	}
	if msg.Nick() == client.nick {
		client.nick = newNick
		client.printStatus("you are now known as %s\n", newNick)
	} else {
		oldNick := msg.Nick()
		for _, channel := range client.channels {
			if slices.Contains(channel.members, oldNick) {
				client.printChannel(channel.name, "%s is now known as %s\n", oldNick, newNick)
				channel.renameMember(oldNick, newNick)
			}
		}
	}
	client.refreshNames()
}
