package mqtt

import (
	"cod-server/internal/api"
	"fmt"
	"log"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/google/uuid"
)

type MQTTAdapterInterface interface {
	Connect() error
	Publish(topic string, event api.Event) error
	Subscribe(topic string, handler mqtt.MessageHandler) error
	Disconnect()
}

type MQTTAdapter struct {
	client mqtt.Client
}

// NewMQTTAdapter cria uma nova instância do adaptador MQTT
func NewMQTTAdapter(broker, clientID string) (MQTTAdapterInterface, error) {
	if clientID == "" {
		clientID = uuid.New().String()
	}

	opts := mqtt.NewClientOptions()
	opts.AddBroker(broker)
	opts.SetClientID(clientID)
	opts.SetDefaultPublishHandler(func(client mqtt.Client, msg mqtt.Message) {
		log.Printf("[MQTT] Received message on topic %s: %s\n", msg.Topic(), msg.Payload())
	})

	opts.OnConnect = func(client mqtt.Client) {
		log.Println("[MQTT] Connected to broker")
	}
	opts.OnConnectionLost = func(client mqtt.Client, err error) {
		log.Printf("[MQTT] Connection lost: %v", err)
	}

	client := mqtt.NewClient(opts)
	return &MQTTAdapter{client: client}, nil
}

// Connect estabelece a conexão com o broker MQTT
func (a *MQTTAdapter) Connect() error {
	if token := a.client.Connect(); token.Wait() && token.Error() != nil {
		return fmt.Errorf("mqtt connection error: %w", token.Error())
	}
	return nil
}

// Publish publica um evento em um tópico MQTT
func (a *MQTTAdapter) Publish(topic string, event api.Event) error {
	payload, err := event.Json()
	if err != nil {
		return fmt.Errorf("failed to serialize event to json: %w", err)
	}

	token := a.client.Publish(topic, 1, false, payload)
	go func() {
		_ = token.Wait()
		if token.Error() != nil {
			log.Printf("[MQTT] Failed to publish to topic %s: %v", topic, token.Error())
		}
	}()

	return nil
}

// Subscribe se inscreve em um tópico MQTT
func (a *MQTTAdapter) Subscribe(topic string, handler mqtt.MessageHandler) error {
	if token := a.client.Subscribe(topic, 1, handler); token.Wait() && token.Error() != nil {
		return fmt.Errorf("failed to subscribe to topic %s: %w", topic, token.Error())
	}
	log.Printf("[MQTT] Subscribed to topic: %s", topic)
	return nil
}

// Disconnect encerra a conexão com o broker MQTT
func (a *MQTTAdapter) Disconnect() {
	log.Println("[MQTT] Disconnecting...")
	a.client.Disconnect(250)
}
