package main

import (
	"strings"
)

func (client *Client) handlePing(msg Message) {
	pong := Message{"", "PONG", msg.parameters}
	client.send(pong)
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
	text := msg.parameters[0]
	if nick == client.nick {
		client.print("you joined %s\n", text)
	} else {
		client.print("%s joined %s\n", nick, text)
	}
}

func (client *Client) handlePart(msg Message) {
	if client.ignored[msg.Nick()] {
		return
	}

	nick := msg.Nick()
	text := msg.parameters[0]
	if nick == client.nick {
		client.print("you left %s\n", text)
	} else {
		client.print("%s left %s\n", nick, text)
	}
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
		}
	}
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
