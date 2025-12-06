package main

import (
	"cod-client/internal/api/protocol"
	"cod-client/internal/commands"
	"cod-client/internal/ui"
	// "fmt"
	// "os"
	// "os/exec"

	//	"slices"
	// "strings"
	// "time"

	//"cod-client/internal/state"
	"cod-client/internal/utils"
)

var (
	mux = utils.NewMux[func(chat *ui.Chat, args []string) *protocol.Message](commands.CommandUnknown)
)

func init() {
	mux.Register("start",    commands.CommandStart)
	mux.Register("surrender", commands.CommandSurrender)
	mux.Register("list",     commands.CommandListCards)
	mux.Register("play",     commands.CommandPlayCard)

	mux.Register("register", commands.CommandRegister)
	mux.Register("login",    commands.CommandLogin)
	mux.Register("logout",   commands.CommandLogout)

	mux.Register("join",     commands.CommandJoin)
	mux.Register("leave",    commands.CommandLeave)

	mux.Register("chat",     commands.CommandChat)

	mux.Register("clear",    commands.CommandClear)
	mux.Register("help",     commands.CommandHelp)
	mux.Register("exit",     commands.CommandExit)
}

func main() {
	chat := ui.NewChat()
	chat.Start(func() {
		var command string
		var args []string
		for input := range chat.Inputs {
			command, args = utils.ParseCommand(input)
			message := mux.Handle(command)(chat, args)
			
		}
	})
	select {}
}
