// Package otel provides OpenTelemetry-specific implementations and utilities
// for MQTT observability. This package implements the observability interfaces
// using OpenTelemetry semantic conventions.
package otelpaho

import (
	"github.com/adibaulia/paho.golang/observability"
)

// Re-export core types from observability package for convenience
type MQTTObserver = observability.MQTTObserver
type ConnectionConfig = observability.ConnectionConfig
type DisconnectReason = observability.DisconnectReason
type PublishMessage = observability.PublishMessage
type PublishResult = observability.PublishResult
type SubscriptionTopic = observability.SubscriptionTopic
type SubscribeResult = observability.SubscribeResult
type UnsubscribeResult = observability.UnsubscribeResult
type ReceivedMessage = observability.ReceivedMessage
type AuthData = observability.AuthData
type AuthResult = observability.AuthResult
type NoOpObserver = observability.NoOpObserver

// SpanAttributes defines common span attributes for MQTT operations
// following OpenTelemetry semantic conventions
type SpanAttributes struct {
	// Messaging system attributes
	MessagingSystem        string // "mqtt"
	MessagingOperation     string // "publish", "receive", "subscribe", etc.
	MessagingDestination   string // topic name
	MessagingProtocol      string // "mqtt"
	MessagingProtocolVersion string // "3.1.1", "5.0"

	// MQTT specific attributes
	MQTTClientID          string
	MQTTBrokerURL         string
	MQTTQoS               int
	MQTTRetain            bool
	MQTTDuplicate         bool
	MQTTMessageID         int
	MQTTKeepAlive         int
	MQTTCleanStart        bool

	// Network attributes
	NetworkTransport      string // "tcp", "tls"
	ServerAddress         string
	ServerPort            int

	// Error attributes
	ErrorType             string
	ErrorMessage          string
}

// ToMap converts SpanAttributes to a map of string key-value pairs
// This can be used with any observability library that accepts map[string]interface{}
func (sa SpanAttributes) ToMap() map[string]interface{} {
	attrs := make(map[string]interface{})

	if sa.MessagingSystem != "" {
		attrs["messaging.system"] = sa.MessagingSystem
	}
	if sa.MessagingOperation != "" {
		attrs["messaging.operation.name"] = sa.MessagingOperation
	}
	if sa.MessagingDestination != "" {
		attrs["messaging.destination.name"] = sa.MessagingDestination
	}
	if sa.MessagingProtocol != "" {
		attrs["messaging.protocol"] = sa.MessagingProtocol
	}
	if sa.MessagingProtocolVersion != "" {
		attrs["messaging.protocol_version"] = sa.MessagingProtocolVersion
	}

	// MQTT specific attributes
	if sa.MQTTClientID != "" {
		attrs["mqtt.client_id"] = sa.MQTTClientID
	}
	if sa.MQTTBrokerURL != "" {
		attrs["mqtt.broker.url"] = sa.MQTTBrokerURL
	}
	if sa.MQTTQoS >= 0 {
		attrs["mqtt.qos"] = sa.MQTTQoS
	}
	if sa.MQTTRetain {
		attrs["mqtt.retain"] = sa.MQTTRetain
	}
	if sa.MQTTDuplicate {
		attrs["mqtt.duplicate"] = sa.MQTTDuplicate
	}
	if sa.MQTTMessageID > 0 {
		attrs["mqtt.message_id"] = sa.MQTTMessageID
	}
	if sa.MQTTKeepAlive > 0 {
		attrs["mqtt.keep_alive"] = sa.MQTTKeepAlive
	}
	if sa.MQTTCleanStart {
		attrs["mqtt.clean_start"] = sa.MQTTCleanStart
	}

	// Network attributes
	if sa.NetworkTransport != "" {
		attrs["network.transport"] = sa.NetworkTransport
	}
	if sa.ServerAddress != "" {
		attrs["server.address"] = sa.ServerAddress
	}
	if sa.ServerPort > 0 {
		attrs["server.port"] = sa.ServerPort
	}

	// Error attributes
	if sa.ErrorType != "" {
		attrs["error.type"] = sa.ErrorType
	}
	if sa.ErrorMessage != "" {
		attrs["error.message"] = sa.ErrorMessage
	}

	return attrs
}

// MetricNames defines standard metric names for MQTT operations
type MetricNames struct {
	// Connection metrics
	ConnectionsActive    string
	ConnectionsTotal     string
	ConnectionDuration   string
	ReconnectAttempts    string

	// Message metrics
	MessagesPublished    string
	MessagesReceived     string
	MessageSize          string
	PublishLatency       string
	SubscribeLatency     string

	// Error metrics
	ErrorsTotal          string
	Timeouts             string

	// Protocol metrics
	PingLatency          string
	AuthAttempts         string
}

// DefaultMetricNames returns the default metric names following OpenTelemetry conventions
func DefaultMetricNames() MetricNames {
	return MetricNames{
		// Connection metrics
		ConnectionsActive:    "mqtt.client.connections.active",
		ConnectionsTotal:     "mqtt.client.connections.total",
		ConnectionDuration:   "mqtt.client.connection.duration",
		ReconnectAttempts:    "mqtt.client.reconnect.attempts",

		// Message metrics
		MessagesPublished:    "mqtt.client.messages.published",
		MessagesReceived:     "mqtt.client.messages.received",
		MessageSize:          "mqtt.client.message.size",
		PublishLatency:       "mqtt.client.publish.duration",
		SubscribeLatency:     "mqtt.client.subscribe.duration",

		// Error metrics
		ErrorsTotal:          "mqtt.client.errors.total",
		Timeouts:             "mqtt.client.timeouts.total",

		// Protocol metrics
		PingLatency:          "mqtt.client.ping.duration",
		AuthAttempts:         "mqtt.client.auth.attempts",
	}
}