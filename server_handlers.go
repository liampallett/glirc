package main

import (
	"fmt"
	"strings"
)

func (client *Client) handlePing(msg Message) {
	pong := Message{"", "PONG", msg.parameters}
	client.send(pong)
}

func (client *Client) handlePrivmsg(msg Message) {
	nick := msg.Nick()
	text := msg.parameters[1]
	if strings.HasPrefix(text, "\x01ACTION ") && strings.HasSuffix(text, "\x01") {
		fmt.Printf("* %s %s\n", nick, text[8:len(text)-1])
	} else {
		fmt.Printf("<%s> %s\n", nick, text)
	}
}

func (client *Client) handleJoin(msg Message) {
	nick := msg.Nick()
	text := msg.parameters[0]
	if nick == client.nick {
		fmt.Printf("you joined %s\n", text)
	} else {
		fmt.Printf("%s joined %s\n", nick, text)
	}
}

func (client *Client) handlePart(msg Message) {
	nick := msg.Nick()
	text := msg.parameters[0]
	if nick == client.nick {
		fmt.Printf("you left %s\n", text)
	} else {
		fmt.Printf("%s left %s\n", nick, text)
	}
}

func (client *Client) handleQuit(msg Message) {
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
}

func (client *Client) handleNick(msg Message) {
	if msg.Nick() == client.nick {
		newNick := msg.parameters[0]
		client.nick = newNick
		fmt.Printf("you are now known as %s\n", newNick)
	} else {
		oldNick := msg.Nick()
		newNick := msg.parameters[0]
		fmt.Printf("%s is now known as %s\n", oldNick, newNick)
	}
}
