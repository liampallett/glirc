package main

import (
	"slices"
	"strings"
)

func (client *Client) handleMOTDStart(msg Message) {
	client.print("%s\n", msg.parameters[1])
}

func (client *Client) handleMOTD(msg Message) {
	client.print("%s\n", msg.parameters[1])
}

func (client *Client) handleMOTDEnd(msg Message) {
	client.print("%s\n", msg.parameters[1])
}

func (client *Client) handleNames(msg Message) {
	channel := msg.parameters[2]
	members := strings.Fields(msg.parameters[len(msg.parameters)-1])
	for _, element := range members {
		client.channelMembers[channel] = append(client.channelMembers[channel], element)
	}
	client.refreshNames()
}

func (client *Client) handlePing(msg Message) {
	pong := Message{"", "PONG", msg.parameters}
	err := client.send(pong)
	if err != nil {
		client.print("%s\n", err)
	}
}

func (client *Client) handlePrivmsg(msg Message) {
	if client.ignored[msg.Nick()] {
		return
	}

	nick := msg.Nick()
	text := msg.parameters[1]
	if strings.HasPrefix(text, "\x01ACTION ") && strings.HasSuffix(text, "\x01") {
		client.print("* %s %s\n", nick, text[8:len(text)-1])
	} else {
		client.print("<%s> %s\n", nick, text)
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
		client.ui.Channels.AddItem(channel, "", 0, nil)
		client.channelMembers[channel] = nil
	} else {
		client.print("%s joined %s\n", nick, channel)
		client.channelMembers[channel] = append(client.channelMembers[channel], nick)
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
		indices := client.ui.Channels.FindItems(channel, "", false, true)
		if len(indices) > 0 {
			client.ui.Channels.RemoveItem(indices[0])
		}
	} else {
		client.print("%s left %s\n", nick, channel)
		client.channelMembers[channel] = slices.DeleteFunc(client.channelMembers[channel], func(s string) bool {
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
			client.print("%s quit\n", nick)
		}
	} else {
		text := msg.parameters[0]
		quitReason := strings.TrimPrefix(text, "Quit: ")
		if nick == client.nick {
			client.print("you quit: %s\n", quitReason)
		} else {
			client.print("%s quit: %s\n", nick, quitReason)
			for channel := range client.channelMembers {
				client.channelMembers[channel] = slices.DeleteFunc(client.channelMembers[channel], func(s string) bool {
					return s == nick
				})
			}
		}
	}
	client.refreshNames()
}

func (client *Client) handleNick(msg Message) {
	if msg.Nick() == client.nick {
		newNick := msg.parameters[0]
		client.nick = newNick
		client.print("you are now known as %s\n", newNick)
	} else {
		oldNick := msg.Nick()
		newNick := msg.parameters[0]
		client.print("%s is now known as %s\n", oldNick, newNick)
	}
}
