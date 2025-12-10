package codmqtt

import (
	"encoding/json"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"cod-client/internal/api/protocol"
)

// TimeoutError represents an error that occurs when an operation times out
type TimeoutError string

func (e TimeoutError) Error() string {
	return string(e)
}

// NewTimeoutError creates a new timeout error
func NewTimeoutError(s string) TimeoutError {
	return TimeoutError(s)
}

// PublishEvent serializes an event to JSON and publishes it to the specified MQTT topic
func PublishEvent(client mqtt.Client, topic string, event *protocol.Event) error {
	jsonData, err := json.Marshal(event)
	if err != nil {
		return err
	}
	
	token := client.Publish(topic, 1, false, jsonData)
	token.WaitTimeout(5 * time.Second) // Wait up to 5 seconds for the publish to complete
	
	return token.Error()
}

// SubscribeEvent subscribes to an MQTT topic with a callback handler
func SubscribeEvent(client mqtt.Client, topic string, callback func(mqtt.Client, mqtt.Message)) error {
	token := client.Subscribe(topic, 1, callback)
	token.WaitTimeout(5 * time.Second)
	
	return token.Error()
}

// WaitForResponse waits for a response on a specific topic with timeout
func WaitForResponse(client mqtt.Client, topic string, timeout time.Duration) (*protocol.Event, error) {
	responseChan := make(chan *protocol.Event, 1)
	errorChan := make(chan error, 1)
	
	// Subscribe to the topic
	unsubscribeFunc, err := subscribeOnce(client, topic, responseChan, errorChan)
	if err != nil {
		return nil, err
	}
	
	// Wait for response or timeout
	select {
	case event := <-responseChan:
		unsubscribeFunc()
		return event, nil
	case err := <-errorChan:
		unsubscribeFunc()
		return nil, err
	case <-time.After(timeout):
		unsubscribeFunc()
		return nil, NewTimeoutError("Response timeout")
	}
}

// Helper function for single-use subscription
func subscribeOnce(client mqtt.Client, topic string, responseChan chan<- *protocol.Event, errorChan chan<- error) (func(), error) {
	// Create a unique topic pattern to avoid conflicts with other subscriptions
	callback := func(c mqtt.Client, msg mqtt.Message) {
		var event protocol.Event
		if err := json.Unmarshal(msg.Payload(), &event); err != nil {
			errorChan <- err
			return
		}
		
		responseChan <- &event
	}
	
	token := client.Subscribe(topic, 1, callback)
	token.WaitTimeout(5 * time.Second)
	
	if token.Error() != nil {
		return nil, token.Error()
	}
	
	// Return unsubscribe function
	return func() {
		client.Unsubscribe(topic)
	}, nil
}