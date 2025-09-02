// Package otel provides an example of how to integrate MQTT observability
// This is a simplified example that demonstrates the observability interface
// without external MQTT library dependencies.
package otelpaho

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/adibaulia/paho.golang/observability"
)

// SimpleMQTTClient demonstrates how to integrate observability with MQTT operations
// This is a mock implementation for demonstration purposes
type SimpleMQTTClient struct {
	observer  observability.MQTTObserver
	clientID  string
	connected bool
}

// NewSimpleMQTTClient creates a new MQTT client with observability
func NewSimpleMQTTClient(clientID string, observer observability.MQTTObserver) *SimpleMQTTClient {
	return &SimpleMQTTClient{
		observer:  observer,
		clientID:  clientID,
		connected: false,
	}
}

// Connect establishes a connection to the MQTT broker with observability
func (c *SimpleMQTTClient) Connect(ctx context.Context, brokerURL string, keepAlive time.Duration) error {
	// Create connection configuration
	config := observability.ConnectionConfig{
		ClientID:        c.clientID,
		BrokerURL:       brokerURL,
		KeepAlive:       keepAlive,
		CleanStart:      true,
		ProtocolVersion: 4, // MQTT 3.1.1
		TLSEnabled:      false,
	}

	// Start observability for connection
	ctx, onComplete := c.observer.OnConnect(ctx, config)

	// Simulate connection attempt
	time.Sleep(50 * time.Millisecond) // Simulate network delay

	// Simulate successful connection
	c.connected = true
	onComplete(nil) // No error

	log.Printf("Connected to MQTT broker: %s", brokerURL)
	return nil
}

// Publish sends a message to the specified topic with observability
func (c *SimpleMQTTClient) Publish(ctx context.Context, topic string, payload []byte, qos byte, retain bool) error {
	if !c.connected {
		return fmt.Errorf("not connected to broker")
	}

	// Create publish message
	msg := observability.PublishMessage{
		Topic:     topic,
		Payload:   payload,
		QoS:       qos,
		Retain:    retain,
		MessageID: 1234, // Mock message ID
	}

	// Start observability for publish
	ctx, onComplete := c.observer.OnPublish(ctx, msg)

	// Simulate publish operation
	time.Sleep(10 * time.Millisecond) // Simulate network delay

	// Simulate successful publish
	result := observability.PublishResult{
		MessageID: msg.MessageID,
		Error:     nil,
		Latency:   10 * time.Millisecond,
	}
	onComplete(result)

	log.Printf("Published message to topic: %s (QoS: %d, Size: %d bytes)", topic, qos, len(payload))
	return nil
}

// Subscribe subscribes to one or more topics with observability
func (c *SimpleMQTTClient) Subscribe(ctx context.Context, subscriptions map[string]byte) error {
	if !c.connected {
		return fmt.Errorf("not connected to broker")
	}

	// Convert to subscription topics
	var topics []observability.SubscriptionTopic
	for topic, qos := range subscriptions {
		topics = append(topics, observability.SubscriptionTopic{
			Topic: topic,
			QoS:   qos,
		})
	}

	// Start observability for subscribe
	ctx, onComplete := c.observer.OnSubscribe(ctx, topics)

	// Simulate subscription operation
	time.Sleep(20 * time.Millisecond) // Simulate network delay

	// Simulate successful subscription
	results := make([]byte, len(topics))
	for i, topic := range topics {
		results[i] = topic.QoS // Grant requested QoS
	}

	result := observability.SubscribeResult{
		Topics:  topics,
		Results: results,
		Error:   nil,
		Latency: 20 * time.Millisecond,
	}
	onComplete(result)

	log.Printf("Subscribed to %d topics", len(subscriptions))
	return nil
}

// SimulateMessageReceived simulates receiving a message with observability
func (c *SimpleMQTTClient) SimulateMessageReceived(ctx context.Context, topic string, payload []byte, qos byte) {
	if !c.connected {
		return
	}

	// Create received message
	msg := observability.ReceivedMessage{
		Topic:     topic,
		Payload:   payload,
		QoS:       qos,
		Retain:    false,
		Duplicate: false,
		MessageID: 5678, // Mock message ID
	}

	// Record message received
	ctx, onComplete := c.observer.OnMessageReceived(ctx, msg)
	onComplete(nil) // No error processing message

	log.Printf("Received message on topic %s: %s", topic, string(payload))
}

// Ping sends a ping request with observability
func (c *SimpleMQTTClient) Ping(ctx context.Context) error {
	if !c.connected {
		return fmt.Errorf("not connected to broker")
	}

	// Start observability for ping
	ctx, onComplete := c.observer.OnPing(ctx)

	// Simulate ping operation
	time.Sleep(5 * time.Millisecond) // Simulate network delay

	// Simulate successful ping
	onComplete(nil) // No error

	log.Println("Ping successful")
	return nil
}

// Disconnect closes the connection with observability
func (c *SimpleMQTTClient) Disconnect(ctx context.Context) error {
	if !c.connected {
		return fmt.Errorf("not connected to broker")
	}

	// Create disconnect reason
	reason := observability.DisconnectReason{
		Code:              0,
		Reason:            "client_disconnect",
		InitiatedByClient: true,
	}

	// Record disconnect
	ctx, onComplete := c.observer.OnDisconnect(ctx, reason)

	// Simulate disconnect
	c.connected = false
	onComplete(nil) // No error

	log.Println("Disconnected from MQTT broker")
	return nil
}

// ExampleUsage demonstrates how to use the MQTT observability interface
func ExampleUsage() {
	log.Println("=== MQTT Observability Example ===")

	// Create observer
	observer := NewOTelMQTTObserver("example-client-001")

	// Create MQTT client with observability
	client := NewSimpleMQTTClient("example-client-001", observer)

	ctx := context.Background()

	// Connect to broker
	log.Println("\n1. Connecting to broker...")
	if err := client.Connect(ctx, "tcp://localhost:1883", 30*time.Second); err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}

	// Subscribe to topics
	log.Println("\n2. Subscribing to topics...")
	subscriptions := map[string]byte{
		"sensors/temperature": 1,
		"sensors/humidity":    1,
		"alerts/#":            2,
	}
	if err := client.Subscribe(ctx, subscriptions); err != nil {
		log.Fatalf("Failed to subscribe: %v", err)
	}

	// Publish some messages
	log.Println("\n3. Publishing messages...")
	messages := []struct {
		topic   string
		payload string
		qos     byte
	}{
		{"sensors/temperature", `{"value": 23.5, "unit": "C"}`, 1},
		{"sensors/humidity", `{"value": 65.2, "unit": "%"}`, 1},
		{"alerts/high_temp", `{"message": "Temperature above threshold"}`, 2},
	}

	for _, msg := range messages {
		if err := client.Publish(ctx, msg.topic, []byte(msg.payload), msg.qos, false); err != nil {
			log.Printf("Failed to publish to %s: %v", msg.topic, err)
		}
		time.Sleep(100 * time.Millisecond) // Small delay between messages
	}

	// Simulate receiving messages
	log.Println("\n4. Simulating received messages...")
	client.SimulateMessageReceived(ctx, "sensors/temperature", []byte(`{"value": 24.1, "unit": "C"}`), 1)
	client.SimulateMessageReceived(ctx, "sensors/humidity", []byte(`{"value": 67.8, "unit": "%"}`), 1)

	// Send a ping
	log.Println("\n5. Sending ping...")
	if err := client.Ping(ctx); err != nil {
		log.Printf("Ping failed: %v", err)
	}

	// Wait a bit to see all the observability output
	time.Sleep(500 * time.Millisecond)

	// Disconnect
	log.Println("\n6. Disconnecting...")
	if err := client.Disconnect(ctx); err != nil {
		log.Printf("Failed to disconnect: %v", err)
	}

	log.Println("\n=== Example completed successfully ===")
}

// RunExample is a convenience function to run the example
func RunExample() {
	ExampleUsage()
}
