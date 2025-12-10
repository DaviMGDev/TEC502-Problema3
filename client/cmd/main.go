package main

import (
	"cod-client/internal/api/codmqtt"
	"cod-client/internal/api/protocol"
	"cod-client/internal/state"
	"cod-client/internal/ui"
	"cod-client/internal/utils"
	"encoding/json"
	"os"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/google/uuid"
)

var (
	clientID = uuid.New().String()
	mqttOnlineBroker = "tcp://broker.hivemq.com:1883"
	mqttClient mqtt.Client
	chatTopic = "general"
)

// func decodeEvent(command string, args []string) *protocol.Event {
// 	event := &protocol.Event{
// 		Method:		command,
// 		Timestamp:	time.Now(),
// 		Payload:	utils.Dict{},
// 	}
// 	switch command {
// 	case "chat":
// 		if len(args) >= 1 {
// 			event.Payload["message"] = args[0]
// 		}
// 	case "clear":
// 	case "help":
// 	case "exit":
// 	case "register":
// 		if len(args) >= 2 {
// 			event.Payload["username"] = args[0]
// 			event.Payload["password"] = args[1]
// 		}
// 	case "login":
// 		if len(args) >= 2 {
// 			event.Payload["username"] = args[0]
// 			event.Payload["password"] = args[1]
// 		}
// 	case "start":
// 		if len(args) >= 1 {
// 			event.Payload["user_id"] = args[0]
// 		}
// 	case "play":
// 		if len(args) >= 3 {
// 			event.Payload["match_id"] = args[0]
// 			event.Payload["user_id"] = args[1]
// 			event.Payload["card"] = args[2]
// 		}
// 	case "surrender":
// 		if len(args) >= 2 {
// 			event.Payload["match_id"] = args[0]
// 			event.Payload["user_id"] = args[1]
// 		}
// 	case "list":
// 		if len(args) >= 1 {
// 			event.Payload["user_id"] = args[0]
// 		}
// 	case "trade":
// 		if len(args) >= 3 {
// 			event.Payload["from_user_id"] = args[0]
// 			event.Payload["to_username"] = args[1]
// 			event.Payload["offered_card"] = args[2]
// 		}
// 	case "buy":
// 		if len(args) >= 1 {
// 			event.Payload["user_id"] = args[0]
// 		}
// 	default:
// 	}
// 	return event
// }

// func execEvent(command string, args []string) {
// 	if slices.Contain([]string{"exit", "help", "clear"}, strings.ToLower(command)) {
// 		switch strings.ToLower(command) {
// 		case "exit":
// 			mqttClient.Disconnect(250)
// 			ui.Chat.Write("Exiting...")
// 			time.Sleep(1 * time.Second)
// 			os.Exit(0)
// 		case "help":
// 			helpMessage := "WIP help"		
// 			ui.Chat.Write(helpMessage)
// 		case "clear":
//
// 		}
// 	}

func main() {
	chat := ui.NewChat()
	appState := state.NewAppState()

	mqttClient = mqtt.NewClient(mqtt.NewClientOptions().AddBroker(mqttOnlineBroker).SetClientID(clientID))
	if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	// Subscribe to chat topic
	mqttClient.Subscribe("chat/"+chatTopic, 0, func(client mqtt.Client, msg mqtt.Message) {
		var event protocol.Event
		if err := json.Unmarshal(msg.Payload(), &event); err != nil {
			chat.Write("Error parsing chat message: " + err.Error())
			return
		}
		if message, ok := event.Payload["message"].(string); ok {
			chat.Write(message)
		} else {
			chat.Write("Received invalid chat message format")
		}
	})

	// Subscribe to authentication response topic
	mqttClient.Subscribe("auth/responses/"+clientID, 0, func(client mqtt.Client, msg mqtt.Message) {
		var event protocol.Event
		if err := json.Unmarshal(msg.Payload(), &event); err != nil {
			chat.Write("Error parsing auth response: " + err.Error())
			return
		}

		if event.Method == "login" || event.Method == "register" {
			if userID, ok := event.Payload["user_id"].(string); ok {
				appState.SetUserID(userID)
			}
			if username, ok := event.Payload["username"].(string); ok {
				appState.SetUsername(username)
			}
			if success, ok := event.Payload["success"].(bool); ok && success {
				chat.Write("Successfully authenticated as: " + appState.GetUsername())
			} else {
				if message, ok := event.Payload["message"].(string); ok {
					chat.Write("Authentication failed: " + message)
				} else {
					chat.Write("Authentication failed")
				}
			}
		}
	})

	chat.Start(func() {
		for input := range chat.Inputs {
			command, args := utils.ParseCommand(input)
			func() {
				switch command {
				case "clear":
					chat.Clear()
				case "exit":
					mqttClient.Disconnect(250)
					chat.Write("Exiting...")
					time.Sleep(1 * time.Second)
					os.Exit(0)
				case "help":
					helpMessage := "Available commands: /chat, /register, /login, /start, /play, /surrender, /list, /trade, /buy"
					chat.Write(helpMessage)
				case "chat":
					if len(args) >= 1 {
						event := &protocol.Event{
							Method:    "chat",
							Timestamp: time.Now(),
							Payload:   utils.Dict{
								"message":  args[0],
								"username": appState.GetUsername(),
								"client_id": clientID,
							},
						}
						if err := codmqtt.PublishEvent(mqttClient, "chat/"+chatTopic, event); err != nil {
							chat.Write("Error sending message: " + err.Error())
						}
					}
				case "register", "login":
					if len(args) >= 2 {
						event := &protocol.Event{
							Method:    command,
							Timestamp: time.Now(),
							Payload:   utils.Dict{
								"username":  args[0],
								"password":  args[1],
								"client_id": clientID,
							},
						}
						if err := codmqtt.PublishEvent(mqttClient, "auth/requests/"+command, event); err != nil {
							chat.Write("Error sending " + command + " request: " + err.Error())
						}
					} else {
						chat.Write("Usage: /" + command + " <username> <password>")
					}
				case "start":
					if !appState.IsAuthenticated() {
						chat.Write("You must be logged in to start a game")
						return
					}
					event := &protocol.Event{
						Method:    command,
						Timestamp: time.Now(),
						Payload:   utils.Dict{
							"user_id": appState.GetUserID(),
						},
					}
					if err := codmqtt.PublishEvent(mqttClient, "game/requests/start", event); err != nil {
						chat.Write("Error sending start request: " + err.Error())
					}
				case "play":
					if !appState.IsAuthenticated() {
						chat.Write("You must be logged in to play")
						return
					}
					if len(args) >= 1 {
						event := &protocol.Event{
							Method:    command,
							Timestamp: time.Now(),
							Payload:   utils.Dict{
								"match_id": appState.GetMatchID(),
								"user_id":  appState.GetUserID(),
								"card":     args[0],
							},
						}
						if err := codmqtt.PublishEvent(mqttClient, "game/requests/play", event); err != nil {
							chat.Write("Error sending play request: " + err.Error())
						}
					} else {
						chat.Write("Usage: /play <card>")
					}
				case "surrender":
					if !appState.IsAuthenticated() {
						chat.Write("You must be logged in to surrender")
						return
					}
					event := &protocol.Event{
						Method:    command,
						Timestamp: time.Now(),
						Payload:   utils.Dict{
							"match_id": appState.GetMatchID(),
							"user_id":  appState.GetUserID(),
						},
					}
					if err := codmqtt.PublishEvent(mqttClient, "game/requests/surrender", event); err != nil {
						chat.Write("Error sending surrender request: " + err.Error())
					}
				case "list":
					if !appState.IsAuthenticated() {
						chat.Write("You must be logged in to list cards")
						return
					}
					event := &protocol.Event{
						Method:    command,
						Timestamp: time.Now(),
						Payload:   utils.Dict{
							"user_id": appState.GetUserID(),
						},
					}
					if err := codmqtt.PublishEvent(mqttClient, "trade/requests/list", event); err != nil {
						chat.Write("Error sending list request: " + err.Error())
					}
				case "trade":
					if !appState.IsAuthenticated() {
						chat.Write("You must be logged in to trade")
						return
					}
					if len(args) >= 2 {
						event := &protocol.Event{
							Method:    command,
							Timestamp: time.Now(),
							Payload:   utils.Dict{
								"from_user_id": appState.GetUserID(),
								"to_username":  args[0],
								"offered_card": args[1],
							},
						}
						if err := codmqtt.PublishEvent(mqttClient, "trade/requests/trade", event); err != nil {
							chat.Write("Error sending trade request: " + err.Error())
						}
					} else {
						chat.Write("Usage: /trade <to_username> <offered_card>")
					}
				case "buy":
					if !appState.IsAuthenticated() {
						chat.Write("You must be logged in to buy cards")
						return
					}
					event := &protocol.Event{
						Method:    command,
						Timestamp: time.Now(),
						Payload:   utils.Dict{
							"user_id": appState.GetUserID(),
						},
					}
					if err := codmqtt.PublishEvent(mqttClient, "store/requests/buy", event); err != nil {
						chat.Write("Error sending buy request: " + err.Error())
					}
				default:
					chat.Write("Unknown command. Type /help for a list of commands.")
				}
			}()
		}
	})
	select {}
}
