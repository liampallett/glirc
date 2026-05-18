package main

import (
	"fmt"
	"slices"
	"strings"
)

func (client *Client) handleMOTDStart(msg Message) {
	if len(msg.parameters) < 2 {
		return
	}
	client.print("%s\n", msg.parameters[1])
}

func (client *Client) handleMOTD(msg Message) {
	if len(msg.parameters) < 2 {
		return
	}
	if len(msg.parameters[1]) < 2 {
		return
	}
	client.print("%s\t\n", msg.parameters[1][2:])
}

func (client *Client) handleMOTDEnd(msg Message) {
	client.print("\n")
}

func (client *Client) handleNotice(msg Message) {
	if len(msg.parameters) < 2 {
		return
	}
	client.print("%s %s\n", msg.parameters[0], msg.parameters[1])
}

func (client *Client) handleListStart(msg Message) {
	client.print("\n")
}

func (client *Client) handleList(msg Message) {
	if len(msg.parameters) < 3 {
		return
	}
	channel := msg.parameters[1]
	userCount := msg.parameters[2]
	topic := msg.parameters[len(msg.parameters)-1]

	if topic != "" {
		client.print("- %s: %s current users - %s\n", channel, userCount, topic)
	} else {
		client.print("- %s: %s current users\n", channel, userCount)
	}
}

func (client *Client) handleListEnd(msg Message) {
	client.print("\n")
}

func (client *Client) handleNames(msg Message) {
	if len(msg.parameters) < 3 {
		return
	}
	channel := msg.parameters[2]
	members := strings.Fields(msg.parameters[len(msg.parameters)-1])
	client.channels[channel].members = nil
	for _, member := range members {
		ch, ok := client.channels[channel]
		if !ok {
			continue
		}
		ch.members = append(ch.members, member)
	}
	client.refreshNames()
}

func (client *Client) handleNamesEnd(msg Message) {
	client.print("\n")
}

func (client *Client) handlePing(msg Message) {
	pong := Message{"", "PONG", msg.parameters}
	err := client.server.send(pong)
	if err != nil {
		client.print("%s\n", err)
	}
}

func (client *Client) handlePrivmsg(msg Message) {
	if client.ignored[msg.Nick()] {
		return
	}
	if len(msg.parameters) < 2 {
		return
	}

	nick := msg.Nick()
	target := msg.parameters[0]
	text := msg.parameters[1]
	switch target[0] {
	case '#':
		if strings.HasPrefix(text, "\x01ACTION ") && strings.HasSuffix(text, "\x01") {
			client.printChannel(target, "* %s %s\n", nick, text[8:len(text)-1])
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

		if strings.HasPrefix(text, "\x01ACTION ") && strings.HasSuffix(text, "\x01") {
			client.printChannel(key, "* %s %s\n", key, text[8:len(text)-1])
		} else {
			client.printChannel(key, "<%s> %s\n", key, text)
		}
	}

}

func (client *Client) handleJoin(msg Message) {
	if client.ignored[msg.Nick()] {
		return
	}

	nick := msg.Nick()
	channel := msg.parameters[0]
	if nick == client.nick {
		client.print("you joined %s\n", channel)
		client.currentChannel = channel
		client.ui.Chat.SetTitle(channel)
		client.ui.Chat.Clear()
		if ch, ok := client.channels[channel]; ok {
			fmt.Fprintf(client.ui.Chat, strings.Join(ch.history, ""))
		}
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
		ch.members = append(ch.members, nick)
	}
	client.refreshNames()
}

func (client *Client) handlePart(msg Message) {
	if client.ignored[msg.Nick()] {
		return
	}

	nick := msg.Nick()
	channel := msg.parameters[0]
	if nick == client.nick {
		client.print("you left %s\n", channel)
		delete(client.channels, channel)
		client.currentChannel = ""
		client.ui.Chat.SetTitle(client.currentChannel)
		indices := client.ui.Channels.FindItems(channel, "", false, true)
		if len(indices) > 0 {
			client.ui.Channels.RemoveItem(indices[0])
		}
		for name := range client.channels {
			client.currentChannel = name
			break
		}
		client.ui.Chat.SetTitle(client.currentChannel)
		client.ui.Chat.Clear()
		if ch, ok := client.channels[client.currentChannel]; ok {
			fmt.Fprintf(client.ui.Chat, strings.Join(ch.history, ""))
		}
	} else {
		client.printChannel(channel, "%s left %s\n", nick, channel)
		ch, ok := client.channels[channel]
		if !ok {
			return
		}
		ch.members = slices.DeleteFunc(ch.members, func(s string) bool {
			return s == nick
		})
	}
	client.refreshNames()
}

func (client *Client) handleQuit(msg Message) {
	if client.ignored[msg.Nick()] {
		return
	}

	nick := msg.Nick()
	if len(msg.parameters) < 1 {
		if nick == client.nick {
			client.print("you quit\n")
		} else {
			client.printChannel(client.currentChannel, "%s quit\n", nick)
		}
	} else {
		text := msg.parameters[0]
		quitReason := strings.TrimPrefix(text, "Quit: ")
		if nick == client.nick {
			client.print("you quit: %s\n", quitReason)
		} else {
			for _, channel := range client.channels {
				if slices.Contains(channel.members, nick) {
					client.printChannel(channel.name, "%s quit: %s\n", nick, quitReason)
					channel.members = slices.DeleteFunc(channel.members, func(s string) bool {
						return s == nick
					})
				}
			}
		}
	}
	client.refreshNames()
}

func (client *Client) handleNick(msg Message) {
	if len(msg.parameters) < 1 {
		return
	}
	if msg.Nick() == client.nick {
		newNick := msg.parameters[0]
		client.nick = newNick
		client.print("you are now known as %s\n", newNick)
	} else {
		oldNick := msg.Nick()
		newNick := msg.parameters[0]
		for _, channel := range client.channels {
			if slices.Contains(channel.members, oldNick) {
				client.printChannel(channel.name, "%s is now known as %s\n", oldNick, newNick)
				for i := range channel.members {
					if channel.members[i] == oldNick {
						channel.members[i] = newNick
					}
				}
			}
		}
	}
	client.refreshNames()
}
