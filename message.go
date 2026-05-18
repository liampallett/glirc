package main

import (
	"errors"
	"strings"
)

type Message struct {
	prefix     string
	command    string
	parameters []string
}

func parse(line string) (Message, error) {
	if line == "" {
		return Message{}, errors.New("attempted to parse an empty line")
	}

	var prefix string
	var command string
	var parameters []string

	if line[0] == ':' {
		spaceIndex := strings.Index(line, " ")
		if spaceIndex == -1 {
			return Message{}, errors.New("prefix with no command")
		}
		prefix = line[1:spaceIndex]
		line = line[spaceIndex+1:]
	}

	spaceIndex := strings.Index(line, " ")
	if spaceIndex == -1 {
		command = line
		return Message{prefix, command, parameters}, nil
	}
	command = line[0:spaceIndex]
	line = line[spaceIndex+1:]

	for line != "" {
		if line[0] == ':' {
			parameters = append(parameters, line[1:])
			break
		}
		spaceIndex := strings.Index(line, " ")
		if spaceIndex == -1 {
			parameters = append(parameters, line)
			break
		}
		parameters = append(parameters, line[0:spaceIndex])
		line = line[spaceIndex+1:]
	}

	return Message{prefix, command, parameters}, nil
}

func parseAction(text string) (string, bool) {
	if strings.HasPrefix(text, "\x01ACTION ") && strings.HasSuffix(text, "\x01") {
		return text[8 : len(text)-1], true
	}
	return "", false
}

func (msg Message) String() string {
	var builder strings.Builder

	if msg.prefix != "" {
		builder.WriteByte(':')
		builder.WriteString(msg.prefix)
		builder.WriteByte(' ')
	}

	builder.WriteString(msg.command)

	if len(msg.parameters) > 0 {
		for _, element := range msg.parameters[:len(msg.parameters)-1] {
			builder.WriteByte(' ')
			builder.WriteString(element)
		}
		last := msg.parameters[len(msg.parameters)-1]
		if strings.ContainsAny(last, " ") || strings.HasPrefix(last, ":") {
			builder.WriteString(" :")
		} else {
			builder.WriteString(" ")
		}
		builder.WriteString(last)
	}

	builder.WriteString("\r\n")

	return builder.String()
}

func (msg Message) Nick() string {
	if !strings.Contains(msg.prefix, "!") {
		return ""
	}
	return strings.SplitN(msg.prefix, "!", 2)[0]
}

func (msg Message) param(n int) (string, bool) {
	if n < len(msg.parameters) {
		return msg.parameters[n], true
	}
	return "", false
}
