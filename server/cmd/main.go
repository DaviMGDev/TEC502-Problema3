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

// Define MQTT connection options
var (
	broker   = "tcp://broker.hivemq.com:1883" // HiveMQ Public Broker - Altere para o seu broker se desejar
	clientID = "your_unique_client_id"      // <<< ALtere para um ID único para seu servidor
	username = "your_mqtt_username"       // <<< ALtere para seu usuário MQTT, se o broker exigir (deixe vazio ou "")
	password = "your_mqtt_password"       // <<< ALtere para sua senha MQTT, se o broker exigir (deixe vazio ou "")
)

func main() {
	log.Println("Starting Cards of Destiny Server...")

	// 1. Initialize Repositories
	userRepo := data.NewInMemoryRepository[*domain.User]()
	matchRepo := data.NewInMemoryRepository[*domain.Match]()
	// Package repo is implicitly handled by CardService

	// 2. Create Services
	userService := services.NewUserService(userRepo)
	cardService := services.NewCardService(userRepo) // CardService needs userRepo
	gameService := services.NewGameService(matchRepo, userRepo) // GameService needs matchRepo and userRepo

	// 3. Create Handlers
	handlersImpl := handlers.NewHandlers(userService, cardService, gameService)

	// 4. Configure MQTT Client
	opts := mqtt.NewClientOptions().AddBroker(broker).SetClientID(clientID)
	opts.SetUsername(username)
	opts.SetPassword(password)
	opts.SetKeepAlive(60 * time.Second)
	opts.SetPingTimeout(1 * time.Second)
	opts.SetDefaultPublishHandler(func(client mqtt.Client, msg mqtt.Message) {
		log.Printf("Received unhandled message: TOPIC: %s, MSG: %s\n", msg.Topic(), msg.Payload())
	})

	// Add a OnConnectHandler to automatically subscribe when connected
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

	// 5. Setup signal handling for graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	// Block until a signal is received
	<-c
	log.Println("Shutting down server...")

	// Disconnect MQTT client
	if client.IsConnected() {
		client.Disconnect(250)
		log.Println("Disconnected from MQTT broker.")
	}

	log.Println("Server stopped.")
}
