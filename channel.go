package main

import "slices"

type Channel struct {
	name    string
	members []string
	history []string
}

func (channel *Channel) addMember(nick string) {
	channel.members = append(channel.members, nick)
}

func (channel *Channel) removeMember(nick string) {
	channel.members = slices.DeleteFunc(channel.members, func(s string) bool {
		return s == nick
	})
}

func (channel *Channel) renameMember(oldNick, newNick string) {
	for i := range channel.members {
		if channel.members[i] == oldNick {
			channel.members[i] = newNick
		}
	}
}
