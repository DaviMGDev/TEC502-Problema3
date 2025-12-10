package main

import (
	"cod-server/internal/api/codmqtt"
	"cod-server/internal/data"
	"cod-server/internal/domain"
	"cod-server/internal/handlers"
	"cod-server/internal/services"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	broker		= "tcp://broker.hivemq.com:1883"
	clientID	= "your_unique_client_id"
	username	= "your_mqtt_username"
	password	= "your_mqtt_password"
)

func main() {
	log.Println("Starting Cards of Destiny Server...")

	userRepo := data.NewInMemoryRepository[*domain.User]()
	matchRepo := data.NewInMemoryRepository[*domain.Match]()

	userService := services.NewUserService(userRepo)
	cardService := services.NewCardService(userRepo)
	gameService := services.NewGameService(matchRepo, userRepo)

	handlersImpl := handlers.NewHandlers(userService, cardService, gameService)

	opts := mqtt.NewClientOptions().AddBroker(broker).SetClientID(clientID)
	opts.SetUsername(username)
	opts.SetPassword(password)
	opts.SetKeepAlive(60 * time.Second)
	opts.SetPingTimeout(1 * time.Second)
	opts.SetDefaultPublishHandler(func(client mqtt.Client, msg mqtt.Message) {
		log.Printf("Received unhandled message: TOPIC: %s, MSG: %s\n", msg.Topic(), msg.Payload())
	})

	opts.SetOnConnectHandler(func(client mqtt.Client) {
		log.Println("MQTT Client Connected")
		mqttHandler := codmqtt.NewMQTTHandler(client, handlersImpl)
		if err := mqttHandler.SubscribeToEvents(); err != nil {
			log.Fatalf("Error subscribing to events: %v", err)
		}
		log.Println("Subscribed to all event topics.")
	})

	opts.SetConnectionLostHandler(func(client mqtt.Client, err error) {
		log.Printf("MQTT Connection Lost: %v", err)
	})

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatalf("Failed to connect to MQTT broker: %v", token.Error())
	}
	log.Println("Connected to MQTT broker.")

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	<-c
	log.Println("Shutting down server...")

	if client.IsConnected() {
		client.Disconnect(250)
		log.Println("Disconnected from MQTT broker.")
	}

	log.Println("Server stopped.")
}
