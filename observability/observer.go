// Package observability provides interfaces for observing MQTT client operations.
// This package defines core observability interfaces that are independent of any
// specific observability implementation (OpenTelemetry, Prometheus, etc.).
package observability

import (
	"context"
	"time"
)

// MQTTObserver defines the interface for observing MQTT client operations.
// Implementations can instrument MQTT operations with traces, metrics, and logs
// using any observability framework.
type MQTTObserver interface {
	// Connection Lifecycle Operations
	OnConnect(ctx context.Context, config ConnectionConfig) (context.Context, func(error))
	OnDisconnect(ctx context.Context, reason DisconnectReason) (context.Context, func(error))
	OnReconnect(ctx context.Context, attempt int, backoff time.Duration) (context.Context, func(error))

	// Message Operations
	OnPublish(ctx context.Context, msg PublishMessage) (context.Context, func(PublishResult))
	OnSubscribe(ctx context.Context, topics []SubscriptionTopic) (context.Context, func(SubscribeResult))
	OnUnsubscribe(ctx context.Context, topics []string) (context.Context, func(UnsubscribeResult))
	OnMessageReceived(ctx context.Context, msg ReceivedMessage) (context.Context, func(error))

	// Protocol Operations
	OnPing(ctx context.Context) (context.Context, func(error))
	OnAuth(ctx context.Context, authData AuthData) (context.Context, func(AuthResult))

	// Error and Performance Tracking
	RecordError(ctx context.Context, operation string, err error)
	RecordLatency(ctx context.Context, operation string, duration time.Duration)
	RecordMessageSize(ctx context.Context, operation string, size int64)
}

// ConnectionConfig represents connection configuration for observability
type ConnectionConfig struct {
	ClientID        string
	BrokerURL       string
	KeepAlive       time.Duration
	CleanStart      bool
	ProtocolVersion int
	TLSEnabled      bool
}

// DisconnectReason represents the reason for disconnection
type DisconnectReason struct {
	Code              uint8
	Reason            string
	InitiatedByClient bool
}

// PublishMessage represents a message being published
type PublishMessage struct {
	Topic     string
	Payload   []byte
	QoS       byte
	Retain    bool
	MessageID uint16
}

// PublishResult represents the result of a publish operation
type PublishResult struct {
	MessageID uint16
	Error     error
	Latency   time.Duration
}

// SubscriptionTopic represents a topic subscription
type SubscriptionTopic struct {
	Topic string
	QoS   byte
}

// SubscribeResult represents the result of a subscribe operation
type SubscribeResult struct {
	Topics  []SubscriptionTopic
	Results []byte // QoS granted for each topic
	Error   error
	Latency time.Duration
}

// UnsubscribeResult represents the result of an unsubscribe operation
type UnsubscribeResult struct {
	Topics  []string
	Error   error
	Latency time.Duration
}

// ReceivedMessage represents a message received from the broker
type ReceivedMessage struct {
	Topic     string
	Payload   []byte
	QoS       byte
	Retain    bool
	Duplicate bool
	MessageID uint16
}

// AuthData represents authentication data
type AuthData struct {
	Method string
	Data   []byte
}

// AuthResult represents the result of an authentication operation
type AuthResult struct {
	Success bool
	Error   error
	Latency time.Duration
}

// NoOpObserver provides a no-operation implementation of MQTTObserver
// for cases where observability is disabled
type NoOpObserver struct{}

// Ensure NoOpObserver implements MQTTObserver
var _ MQTTObserver = (*NoOpObserver)(nil)

func (n *NoOpObserver) OnConnect(ctx context.Context, config ConnectionConfig) (context.Context, func(error)) {
	return ctx, func(error) {}
}

func (n *NoOpObserver) OnDisconnect(ctx context.Context, reason DisconnectReason) (context.Context, func(error)) {
	return ctx, func(error) {}
}

func (n *NoOpObserver) OnReconnect(ctx context.Context, attempt int, backoff time.Duration) (context.Context, func(error)) {
	return ctx, func(error) {}
}

func (n *NoOpObserver) OnPublish(ctx context.Context, msg PublishMessage) (context.Context, func(PublishResult)) {
	return ctx, func(PublishResult) {}
}

func (n *NoOpObserver) OnSubscribe(ctx context.Context, topics []SubscriptionTopic) (context.Context, func(SubscribeResult)) {
	return ctx, func(SubscribeResult) {}
}

func (n *NoOpObserver) OnUnsubscribe(ctx context.Context, topics []string) (context.Context, func(UnsubscribeResult)) {
	return ctx, func(UnsubscribeResult) {}
}

func (n *NoOpObserver) OnMessageReceived(ctx context.Context, msg ReceivedMessage) (context.Context, func(error)) {
	return ctx, func(error) {}
}

func (n *NoOpObserver) OnPing(ctx context.Context) (context.Context, func(error)) {
	return ctx, func(error) {}
}

func (n *NoOpObserver) OnAuth(ctx context.Context, authData AuthData) (context.Context, func(AuthResult)) {
	return ctx, func(AuthResult) {}
}

func (n *NoOpObserver) RecordError(ctx context.Context, operation string, err error) {}

func (n *NoOpObserver) RecordLatency(ctx context.Context, operation string, duration time.Duration) {}

func (n *NoOpObserver) RecordMessageSize(ctx context.Context, operation string, size int64) {}