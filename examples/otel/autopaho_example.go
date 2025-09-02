package main

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/adibaulia/paho.golang/autopaho"
	"github.com/adibaulia/paho.golang/otelpaho"
	"github.com/adibaulia/paho.golang/paho"
)

// This example demonstrates how to use the autopaho client with OpenTelemetry observability.
// It shows how to:
// 1. Set up OpenTelemetry tracing
// 2. Configure autopaho with an observability observer
// 3. Publish and subscribe to messages with automatic tracing

func exampleMain() {
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
	observer := otelpaho.NewOTelMQTTObserver("autopaho-otel-example")

	// Configure autopaho client with observability
	cfg := autopaho.ClientConfig{
		ServerUrls:                    []*url.URL{u},
		KeepAlive:                     30,
		CleanStartOnInitialConnection: true,
		SessionExpiryInterval:         60,
		ConnectRetryDelay:             time.Second,
		ConnectTimeout:                10 * time.Second,
		ClientConfig: paho.ClientConfig{
			ClientID: "autopaho-otel-example",
			Router:   paho.NewStandardRouter(),
			Observer: observer,
		},
	}

	// Add message handler
	cfg.ClientConfig.Router.RegisterHandler("test/topic", func(p *paho.Publish) {
		fmt.Printf("Received message on %s: %s\n", p.Topic, string(p.Payload))
	})

	// Create and start the client
	client, err := autopaho.NewConnection(context.Background(), cfg)
	if err != nil {
		log.Fatalf("Failed to create autopaho client: %v", err)
	}

	// Wait for connection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := client.AwaitConnection(ctx); err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	fmt.Println("Connected to MQTT broker")

	// Subscribe to a topic (with observability tracing)
	subCtx := context.Background()
	_, err = client.Subscribe(subCtx, &paho.Subscribe{
		Subscriptions: []paho.SubscribeOptions{
			{Topic: "test/topic", QoS: 1},
		},
	})
	if err != nil {
		log.Fatalf("Failed to subscribe: %v", err)
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
				message := fmt.Sprintf("Hello from autopaho with observability! Message #%d", counter)

				pubCtx := context.Background()
				_, err := client.Publish(pubCtx, &paho.Publish{
					Topic:   "test/topic",
					Payload: []byte(message),
					QoS:     1,
				})
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

	// Disconnect gracefully
	disconnectCtx, disconnectCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer disconnectCancel()

	if err := client.Disconnect(disconnectCtx); err != nil {
		log.Printf("Error during disconnect: %v", err)
	}

	fmt.Println("Disconnected successfully")
}

func main() {
	exampleMain()
}
