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
	"github.com/adibaulia/paho.golang/otel"
	"github.com/adibaulia/paho.golang/paho"
)

// This example demonstrates how to use the autopaho client with OpenTelemetry observability.
// It shows how to:
// 1. Configure autopaho client with an observability observer
// 2. Automatically handle connections, reconnections, and message routing with tracing
// 3. Publish and subscribe to messages with automatic observability

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
	observer := otel.NewOTelMQTTObserver("autopaho-otel-example")

	// Create connection up channel
	connUpChan := make(chan struct{})

	// Configure autopaho client with observability
	cfg := autopaho.ClientConfig{
		ServerUrls:                    []*url.URL{u},
		KeepAlive:                     30,
		CleanStartOnInitialConnection: true,
		SessionExpiryInterval:         60,
		ConnectRetryDelay:             time.Second,
		ConnectTimeout:                10 * time.Second,
		OnConnectionUp: func(*autopaho.ConnectionManager, *paho.Connack) {
			select {
			case connUpChan <- struct{}{}:
			default:
			}
		},
		ClientConfig: paho.ClientConfig{
			ClientID: "autopaho-otel-example",
			Router:   paho.NewStandardRouter(),
			Observer: observer, // Add observability
		},
	}

	// Add message handler
	cfg.ClientConfig.Router.RegisterHandler("test/topic", func(p *paho.Publish) {
		fmt.Printf("Received message on %s: %s\n", p.Topic, string(p.Payload))
	})

	// Create the autopaho connection manager
	cm, err := autopaho.NewConnection(context.Background(), cfg)
	if err != nil {
		log.Fatalf("Failed to create connection manager: %v", err)
	}

	// Wait for initial connection
	select {
	case <-cm.Done():
		log.Fatalf("Connection manager terminated unexpectedly")
	case <-time.After(10 * time.Second):
		log.Fatalf("Connection timeout")
	case <-connUpChan:
		fmt.Println("Connected to MQTT broker")
	}

	// Subscribe to a topic (with observability tracing)
	subCtx, subCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer subCancel()

	suback, err := cm.Subscribe(subCtx, &paho.Subscribe{
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
				message := fmt.Sprintf("Hello from autopaho with observability! Message #%d", counter)
				
				pubCtx, pubCancel := context.WithTimeout(context.Background(), 5*time.Second)
				_, err := cm.Publish(pubCtx, &paho.Publish{
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
			case <-cm.Done():
				return
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

	_, err = cm.Unsubscribe(unsubCtx, &paho.Unsubscribe{
		Topics: []string{"test/topic"},
	})
	if err != nil {
		log.Printf("Error during unsubscribe: %v", err)
	} else {
		fmt.Println("Unsubscribed from test/topic")
	}

	// Disconnect gracefully (autopaho handles this automatically)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = cm.Disconnect(ctx)
	if err != nil {
		log.Printf("Error during disconnect: %v", err)
	} else {
		fmt.Println("Disconnected successfully")
	}
}