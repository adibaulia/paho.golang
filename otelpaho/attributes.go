// Package otelpaho provides detailed attribute and metric definitions for MQTT observability
// following OpenTelemetry semantic conventions for messaging systems.
package otelpaho

import (
	"fmt"
	"time"
)

// MQTT Semantic Conventions
// Based on OpenTelemetry semantic conventions for messaging systems
// https://opentelemetry.io/docs/specs/semconv/messaging/

const (
	// Messaging System Attributes
	MessagingSystemKey        = "messaging.system"        // "mqtt"
	MessagingOperationNameKey = "messaging.operation.name" // "publish", "receive", "subscribe", etc.
	MessagingDestinationNameKey = "messaging.destination.name" // topic name
	MessagingProtocolKey      = "messaging.protocol"      // "mqtt"
	MessagingProtocolVersionKey = "messaging.protocol_version" // "3.1.1", "5.0"
	
	// MQTT Specific Attributes
	MQTTClientIDKey          = "mqtt.client_id"
	MQTTBrokerURLKey         = "mqtt.broker.url"
	MQTTQoSKey               = "mqtt.qos"
	MQTTRetainKey            = "mqtt.retain"
	MQTTDuplicateKey         = "mqtt.duplicate"
	MQTTMessageIDKey         = "mqtt.message_id"
	MQTTKeepAliveKey         = "mqtt.keep_alive"
	MQTTCleanStartKey        = "mqtt.clean_start"
	MQTTSessionPresentKey    = "mqtt.session_present"
	MQTTReasonCodeKey        = "mqtt.reason_code"
	MQTTReasonStringKey      = "mqtt.reason_string"
	MQTTUserPropertiesKey    = "mqtt.user_properties"
	MQTTPayloadSizeKey       = "mqtt.payload.size"
	MQTTTopicAliasKey        = "mqtt.topic_alias"
	MQTTSubscriptionIDKey    = "mqtt.subscription_id"
	
	// Network Attributes
	NetworkTransportKey      = "network.transport"  // "tcp", "tls", "ws", "wss"
	ServerAddressKey         = "server.address"
	ServerPortKey            = "server.port"
	
	// Error Attributes
	ErrorTypeKey             = "error.type"
	ErrorMessageKey          = "error.message"
	
	// Performance Attributes
	MQTTConnectionLatencyKey = "mqtt.connection.latency"
	MQTTPublishLatencyKey    = "mqtt.publish.latency"
	MQTTSubscribeLatencyKey  = "mqtt.subscribe.latency"
	MQTTMessageLatencyKey    = "mqtt.message.latency"
)

// MQTT Operation Names
const (
	OperationConnect     = "connect"
	OperationDisconnect  = "disconnect"
	OperationReconnect   = "reconnect"
	OperationPublish     = "publish"
	OperationReceive     = "receive"
	OperationSubscribe   = "subscribe"
	OperationUnsubscribe = "unsubscribe"
	OperationPing        = "ping"
	OperationAuth        = "auth"
)

// MQTT Protocol Versions
const (
	ProtocolVersion311 = "3.1.1"
	ProtocolVersion50  = "5.0"
)

// MQTT QoS Levels
const (
	QoSAtMostOnce  = 0
	QoSAtLeastOnce = 1
	QoSExactlyOnce = 2
)

// Network Transport Types
const (
	TransportTCP = "tcp"
	TransportTLS = "tls"
	TransportWS  = "ws"
	TransportWSS = "wss"
)

// AttributeBuilder helps build consistent attribute sets for MQTT operations
type AttributeBuilder struct {
	attributes map[string]interface{}
}

// NewAttributeBuilder creates a new attribute builder
func NewAttributeBuilder() *AttributeBuilder {
	return &AttributeBuilder{
		attributes: make(map[string]interface{}),
	}
}

// WithMessagingSystem sets the messaging system (always "mqtt")
func (ab *AttributeBuilder) WithMessagingSystem() *AttributeBuilder {
	ab.attributes[MessagingSystemKey] = "mqtt"
	return ab
}

// WithOperation sets the messaging operation name
func (ab *AttributeBuilder) WithOperation(operation string) *AttributeBuilder {
	ab.attributes[MessagingOperationNameKey] = operation
	return ab
}

// WithDestination sets the messaging destination (topic)
func (ab *AttributeBuilder) WithDestination(topic string) *AttributeBuilder {
	ab.attributes[MessagingDestinationNameKey] = topic
	return ab
}

// WithProtocol sets the MQTT protocol version
func (ab *AttributeBuilder) WithProtocol(version string) *AttributeBuilder {
	ab.attributes[MessagingProtocolKey] = "mqtt"
	ab.attributes[MessagingProtocolVersionKey] = version
	return ab
}

// WithClientID sets the MQTT client ID
func (ab *AttributeBuilder) WithClientID(clientID string) *AttributeBuilder {
	ab.attributes[MQTTClientIDKey] = clientID
	return ab
}

// WithBrokerURL sets the MQTT broker URL
func (ab *AttributeBuilder) WithBrokerURL(url string) *AttributeBuilder {
	ab.attributes[MQTTBrokerURLKey] = url
	return ab
}

// WithQoS sets the MQTT QoS level
func (ab *AttributeBuilder) WithQoS(qos int) *AttributeBuilder {
	ab.attributes[MQTTQoSKey] = qos
	return ab
}

// WithRetain sets the MQTT retain flag
func (ab *AttributeBuilder) WithRetain(retain bool) *AttributeBuilder {
	if retain {
		ab.attributes[MQTTRetainKey] = retain
	}
	return ab
}

// WithDuplicate sets the MQTT duplicate flag
func (ab *AttributeBuilder) WithDuplicate(duplicate bool) *AttributeBuilder {
	if duplicate {
		ab.attributes[MQTTDuplicateKey] = duplicate
	}
	return ab
}

// WithMessageID sets the MQTT message ID
func (ab *AttributeBuilder) WithMessageID(messageID int) *AttributeBuilder {
	if messageID > 0 {
		ab.attributes[MQTTMessageIDKey] = messageID
	}
	return ab
}

// WithKeepAlive sets the MQTT keep alive interval
func (ab *AttributeBuilder) WithKeepAlive(keepAlive time.Duration) *AttributeBuilder {
	if keepAlive > 0 {
		ab.attributes[MQTTKeepAliveKey] = int(keepAlive.Seconds())
	}
	return ab
}

// WithCleanStart sets the MQTT clean start flag
func (ab *AttributeBuilder) WithCleanStart(cleanStart bool) *AttributeBuilder {
	ab.attributes[MQTTCleanStartKey] = cleanStart
	return ab
}

// WithSessionPresent sets the MQTT session present flag
func (ab *AttributeBuilder) WithSessionPresent(sessionPresent bool) *AttributeBuilder {
	ab.attributes[MQTTSessionPresentKey] = sessionPresent
	return ab
}

// WithReasonCode sets the MQTT reason code
func (ab *AttributeBuilder) WithReasonCode(code uint8) *AttributeBuilder {
	ab.attributes[MQTTReasonCodeKey] = int(code)
	return ab
}

// WithReasonString sets the MQTT reason string
func (ab *AttributeBuilder) WithReasonString(reason string) *AttributeBuilder {
	if reason != "" {
		ab.attributes[MQTTReasonStringKey] = reason
	}
	return ab
}

// WithPayloadSize sets the message payload size
func (ab *AttributeBuilder) WithPayloadSize(size int) *AttributeBuilder {
	if size > 0 {
		ab.attributes[MQTTPayloadSizeKey] = size
	}
	return ab
}

// WithNetworkTransport sets the network transport type
func (ab *AttributeBuilder) WithNetworkTransport(transport string) *AttributeBuilder {
	ab.attributes[NetworkTransportKey] = transport
	return ab
}

// WithServerAddress sets the server address
func (ab *AttributeBuilder) WithServerAddress(address string) *AttributeBuilder {
	ab.attributes[ServerAddressKey] = address
	return ab
}

// WithServerPort sets the server port
func (ab *AttributeBuilder) WithServerPort(port int) *AttributeBuilder {
	if port > 0 {
		ab.attributes[ServerPortKey] = port
	}
	return ab
}

// WithError sets error attributes
func (ab *AttributeBuilder) WithError(err error) *AttributeBuilder {
	if err != nil {
		ab.attributes[ErrorTypeKey] = fmt.Sprintf("%T", err)
		ab.attributes[ErrorMessageKey] = err.Error()
	}
	return ab
}

// Build returns the final attribute map
func (ab *AttributeBuilder) Build() map[string]interface{} {
	return ab.attributes
}

// Predefined attribute builders for common operations

// ConnectAttributes creates attributes for MQTT connect operations
func ConnectAttributes(config ConnectionConfig) map[string]interface{} {
	return NewAttributeBuilder().
		WithMessagingSystem().
		WithOperation(OperationConnect).
		WithProtocol(ProtocolVersion50). // Default to 5.0, can be overridden
		WithClientID(config.ClientID).
		WithBrokerURL(config.BrokerURL).
		WithKeepAlive(config.KeepAlive).
		WithCleanStart(config.CleanStart).
		Build()
}

// PublishAttributes creates attributes for MQTT publish operations
func PublishAttributes(msg PublishMessage) map[string]interface{} {
	return NewAttributeBuilder().
		WithMessagingSystem().
		WithOperation(OperationPublish).
		WithDestination(msg.Topic).
		WithProtocol(ProtocolVersion50).
		WithQoS(int(msg.QoS)).
		WithRetain(msg.Retain).
		WithMessageID(int(msg.MessageID)).
		WithPayloadSize(len(msg.Payload)).
		Build()
}

// SubscribeAttributes creates attributes for MQTT subscribe operations
func SubscribeAttributes(topics []SubscriptionTopic) map[string]interface{} {
	builder := NewAttributeBuilder().
		WithMessagingSystem().
		WithOperation(OperationSubscribe).
		WithProtocol(ProtocolVersion50)
	
	// For multiple topics, use the first one as the primary destination
	if len(topics) > 0 {
		builder.WithDestination(topics[0].Topic).
			WithQoS(int(topics[0].QoS))
	}
	
	return builder.Build()
}

// ReceiveAttributes creates attributes for MQTT message receive operations
func ReceiveAttributes(msg ReceivedMessage) map[string]interface{} {
	return NewAttributeBuilder().
		WithMessagingSystem().
		WithOperation(OperationReceive).
		WithDestination(msg.Topic).
		WithProtocol(ProtocolVersion50).
		WithQoS(int(msg.QoS)).
		WithRetain(msg.Retain).
		WithDuplicate(msg.Duplicate).
		WithMessageID(int(msg.MessageID)).
		WithPayloadSize(len(msg.Payload)).
		Build()
}

// DisconnectAttributes creates attributes for MQTT disconnect operations
func DisconnectAttributes(reason DisconnectReason) map[string]interface{} {
	builder := NewAttributeBuilder().
		WithMessagingSystem().
		WithOperation(OperationDisconnect).
		WithProtocol(ProtocolVersion50).
		WithReasonCode(reason.Code).
		WithReasonString(reason.Reason)
	
	builder.attributes["mqtt.disconnect.initiated_by_client"] = reason.InitiatedByClient
	
	return builder.Build()
}

// SpanNameBuilder creates consistent span names for MQTT operations
type SpanNameBuilder struct{}

// ForOperation creates a span name for the given operation
func (SpanNameBuilder) ForOperation(operation string) string {
	return fmt.Sprintf("mqtt.%s", operation)
}

// ForPublish creates a span name for publish operations with topic
func (SpanNameBuilder) ForPublish(topic string) string {
	return fmt.Sprintf("mqtt.publish %s", topic)
}

// ForSubscribe creates a span name for subscribe operations with topic
func (SpanNameBuilder) ForSubscribe(topic string) string {
	return fmt.Sprintf("mqtt.subscribe %s", topic)
}

// ForReceive creates a span name for receive operations with topic
func (SpanNameBuilder) ForReceive(topic string) string {
	return fmt.Sprintf("mqtt.receive %s", topic)
}

// Global span name builder instance
var SpanNames = SpanNameBuilder{}