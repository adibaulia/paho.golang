package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/adibaulia/paho.golang/otel"
	"github.com/adibaulia/paho.golang/paho"
)

// This example demonstrates how to use the paho client with OpenTelemetry observability.
// It shows how to:
// 1. Configure paho client with an observability observer
// 2. Connect, publish, and subscribe to messages with automatic tracing
// 3. Handle connection lifecycle events with observability

func main() {
	// Parse broker URL
	brokerURL := "mqtt://localhost:1883"
	if envURL := os.Getenv("MQTT_BROKER_URL"); envURL != "" {
		brokerURL = envURL
	}

	u, err := url.Parse(brokerURL)
	if err != nil {
		log.Fatalf("Failed to parse broker URL: %v", err)
	}

	// Create OpenTelemetry observer
	observer := otel.NewOTelMQTTObserver("paho-otel-example")

	// Create TCP connection
	conn, err := net.Dial("tcp", u.Host)
	if err != nil {
		log.Fatalf("Failed to connect to broker: %v", err)
	}

	// Configure paho client with observability
	cfg := paho.ClientConfig{
		ClientID: "paho-otel-example",
		Conn:     conn,
		Router:   paho.NewStandardRouter(),
		Observer: observer, // Add observability
	}

	// Add message handler
	cfg.Router.RegisterHandler("test/topic", func(p *paho.Publish) {
		fmt.Printf("Received message on %s: %s\n", p.Topic, string(p.Payload))
	})

	// Create the client
	client := paho.NewClient(cfg)

	// Connect to the broker (with observability tracing)
	connectCtx, connectCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer connectCancel()

	connectPacket := &paho.Connect{
		KeepAlive:  30,
		ClientID:   "paho-otel-example",
		CleanStart: true,
	}

	connack, err := client.Connect(connectCtx, connectPacket)
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	if connack.ReasonCode != 0 {
		log.Fatalf("Connection rejected by server: %d", connack.ReasonCode)
	}
	fmt.Println("Connected to MQTT broker")

	// Subscribe to a topic (with observability tracing)
	subCtx, subCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer subCancel()

	suback, err := client.Subscribe(subCtx, &paho.Subscribe{
		Subscriptions: []paho.SubscribeOptions{
			{Topic: "test/topic", QoS: 1},
		},
	})
	if err != nil {
		log.Fatalf("Failed to subscribe: %v", err)
	}
	if len(suback.Reasons) > 0 && suback.Reasons[0] >= 0x80 {
		log.Fatalf("Subscription rejected: %d", suback.Reasons[0])
	}
	fmt.Println("Subscribed to test/topic")

	// Publish messages periodically (with observability tracing)
	go func() {
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()
		counter := 0

		for {
			select {
			case <-ticker.C:
				counter++
				message := fmt.Sprintf("Hello from paho with observability! Message #%d", counter)
				
				pubCtx, pubCancel := context.WithTimeout(context.Background(), 5*time.Second)
				_, err := client.Publish(pubCtx, &paho.Publish{
					Topic:   "test/topic",
					Payload: []byte(message),
					QoS:     1,
				})
				pubCancel()
				
				if err != nil {
					log.Printf("Failed to publish: %v", err)
				} else {
					fmt.Printf("Published: %s\n", message)
				}
			}
		}
	}()

	// Wait for interrupt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	fmt.Println("\nShutting down...")

	// Unsubscribe (with observability tracing)
	unsubCtx, unsubCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer unsubCancel()

	_, err = client.Unsubscribe(unsubCtx, &paho.Unsubscribe{
		Topics: []string{"test/topic"},
	})
	if err != nil {
		log.Printf("Error during unsubscribe: %v", err)
	} else {
		fmt.Println("Unsubscribed from test/topic")
	}

	// Disconnect gracefully (with observability tracing)
	err = client.Disconnect(&paho.Disconnect{
		ReasonCode: 0,
	})
	if err != nil {
		log.Printf("Error during disconnect: %v", err)
	} else {
		fmt.Println("Disconnected successfully")
	}

	// Close the underlying connection
	conn.Close()
}