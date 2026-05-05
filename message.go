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
		return Message{}, errors.New("attemped to parse an empty line")
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
		index := strings.Index(line, " ")
		if index == -1 {
			parameters = append(parameters, line)
			break
		}
		parameters = append(parameters, line[0:index])
		line = strings.SplitN(line, " ", 2)[1]
	}

	return Message{prefix, command, parameters}, nil
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

		builder.WriteString(" :")
		builder.WriteString(msg.parameters[len(msg.parameters)-1])
	}

	builder.WriteString("\r\n")

	return builder.String()
}

func (msg Message) Nick() string {
	return strings.SplitN(msg.prefix, "!", 2)[0]
}
