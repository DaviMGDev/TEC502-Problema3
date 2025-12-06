package api

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"cod-client/internal/ui"
	"cod-client/internal/utils"
	"cod-client/internal/state"
	"cod-client/internal/api/protocol"
	"fmt"
)

func HandleRegisterResponse(chat *ui.Chat, client mqtt.Client, msg mqtt.Message) {
	response, err := protocol.ParseMessage(msg.Payload())
	if err != nil {
		chat.Write("Error parsing registration response.")
		return
	}
	if response.Status == "success" {
		chat.Write(fmt.Sprintf("Registration successful."))
	} else {
		errorMsg := response.Data["error"].(string)
		chat.Write(fmt.Sprintf("Registration failed: %s", errorMsg))
	}
}

func HandleLoginResponse(chat *ui.Chat, client mqtt.Client, msg mqtt.Message) {
	response, err := protocol.ParseMessage(msg.Payload())
	if err != nil {
		chat.Write("Error parsing login response.")
		return
	}
	if response.Status == "success" {
		state.UserID = response.Data["user_id"].(string)
		username := response.Data["username"].(string)
		chat.Write(fmt.Sprintf("Login successful. Welcome, %s!", username))
	} else {
		errorMsg := response.Data["error"].(string)
		chat.Write(fmt.Sprintf("Login failed: %s", errorMsg))
	}
}

func HandleListCardsResponse(chat *ui.Chat, client mqtt.Client, msg mqtt.Message) {
	response, err := protocol.ParseMessage(msg.Payload())
	if err != nil {
		chat.Write("Error parsing list cards response.")
		return
	}
	if response.Status == "success" {
		cards := response.Data["cards"].(string)
		chat.Write("Your cards: " + cards)
	} else {
		errorMsg := response.Data["error"].(string)
		chat.Write(fmt.Sprintf("Failed to list cards: %s", errorMsg))
	}
}

func HandleBuyCardResponse(chat *ui.Chat, client mqtt.Client, msg mqtt.Message) {
	response, err := protocol.ParseMessage(msg.Payload())
	if err != nil {
		chat.Write("Error parsing buy card response.")
		return
	}
	if response.Status == "success" {
		card := response.Data["card"].(string)
		chat.Write(fmt.Sprintf("Successfully bought card: %s", card))
	} else {
		errorMsg := response.Data["error"].(string)
		chat.Write(fmt.Sprintf("Failed to buy card: %s", errorMsg))
	}
}

func HandlePlayCardResponse(chat *ui.Chat, client mqtt.Client, msg mqtt.Message) {
	response, err := protocol.ParseMessage(msg.Payload())
	if err != nil {
		chat.Write("Error parsing play card response.")
		return
	}
	if response.Status == "success" {
		card := response.Data["card"].(string)
		chat.Write(fmt.Sprintf("Successfully played card: %s", card))
	} else {
		errorMsg := response.Data["error"].(string)
		chat.Write(fmt.Sprintf("Failed to play card: %s", errorMsg))
	}
}

func HandleStartGameResponse(chat *ui.Chat, client mqtt.Client, msg mqtt.Message)
