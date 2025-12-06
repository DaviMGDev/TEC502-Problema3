package utils

import (
	"cod-client/internal/api/protocol"
	"encoding/json"
	"strings"
)

func ParseCommand(input string) (string, []string) {
	if len(input) == 0 {
		return "", []string{}
	}
	if input[0] != '/' {
		return "chat", []string{input}
	}
	parts := strings.Fields(input)
	cmd := parts[0][1:]
	args := parts[1:]
	return cmd, args 
}

