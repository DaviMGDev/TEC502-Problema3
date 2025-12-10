package utils

import (
	"strings"
)

func ParseCommand(input string) (string, []string) {
	if input == "" {
		return "", []string{}
	}
	var command string
	var parts, args []string
	if input[0] != '/' {
		return "chat", []string{input}
	}
	parts = strings.Fields(input[1:])
	command = strings.ToLower(parts[0])
	if len(parts) > 1 {
		args = parts[1:]
	} else {
		args = []string{}
	}
	return command, args
}
