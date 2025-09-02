// Package otelpaho provides detailed attribute and metric definitions for MQTT observability
// following OpenTelemetry semantic conventions for messaging systems.
package otelpaho

import (
	"fmt"
	"time"

	"github.com/adibaulia/paho.golang/observability"
	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.37.0"
)

// MQTT Semantic Conventions
// Based on OpenTelemetry semantic conventions for messaging systems
// https://opentelemetry.io/docs/specs/semconv/messaging/

// MQTT Protocol Constants
const (
	// MQTT Specific Attributes
	MQTTClientIDKey       = "mqtt.client_id"
	MQTTBrokerURLKey      = "mqtt.broker.url"
	MQTTQoSKey            = "mqtt.qos"
	MQTTRetainKey         = "mqtt.retain"
	MQTTDuplicateKey      = "mqtt.duplicate"
	MQTTMessageIDKey      = "mqtt.message_id"
	MQTTKeepAliveKey      = "mqtt.keep_alive"
	MQTTCleanStartKey     = "mqtt.clean_start"
	MQTTSessionPresentKey = "mqtt.session_present"
	MQTTReasonCodeKey     = "mqtt.reason_code"
	MQTTReasonStringKey   = "mqtt.reason_string"
	MQTTUserPropertiesKey = "mqtt.user_properties"
	MQTTPayloadSizeKey    = "mqtt.payload.size"
	MQTTTopicAliasKey     = "mqtt.topic_alias"
	MQTTSubscriptionIDKey = "mqtt.subscription_id"
)

// MQTT Metadata
const (
	ProtocolTypeMQTT = "mqtt"

	// Messaging Attribute Keys
	MessagingMQTTQoSKey                = "messaging.mqtt.qos"
	MQTTMessageRetainKey               = "mqtt.message.retain"
	MessagingMessageIDKey              = "messaging.message.id"
	MessagingMessageBodySizeKey        = "messaging.message.body.size"
	MQTTMessageDuplicateKey            = "mqtt.message.duplicate"
	MQTTDisconnectInitiatedByClientKey = "mqtt.disconnect.initiated_by_client"
	MQTTReconnectAttemptKey            = "mqtt.reconnect.attempt"
	MQTTReconnectBackoffKey            = "mqtt.reconnect.backoff"
	MQTTAuthMethodKey                  = "mqtt.auth.method"

	// Span Names
	MQTTClientConnectSpan     = "MQTT Client Connect"
	MQTTClientDisconnectSpan  = "MQTT Client Disconnect"
	MQTTClientReconnectSpan   = "MQTT Client Reconnect"
	MQTTClientPublishSpan     = "MQTT Client Publish"
	MQTTClientSubscribeSpan   = "MQTT Client Subscribe"
	MQTTClientUnsubscribeSpan = "MQTT Client Unsubscribe"
	MQTTClientReceiveSpan     = "MQTT Client Receive"
	MQTTClientPingSpan        = "MQTT Client Ping"
	MQTTClientAuthSpan        = "MQTT Client Auth"

	// Context Keys
	MQTTConnectStartKey = "mqtt.connect.start"
	MQTTPublishStartKey = "mqtt.publish.start"
	MQTTPingStartKey    = "mqtt.ping.start"
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
	ab.attributes[string(semconv.MessagingSystemKey)] = ProtocolTypeMQTT
	return ab
}

// WithOperation sets the messaging operation name
func (ab *AttributeBuilder) WithOperation(operation string) *AttributeBuilder {
	ab.attributes[string(semconv.MessagingOperationNameKey)] = operation
	return ab
}

// WithDestination sets the messaging destination (topic)
func (ab *AttributeBuilder) WithDestination(topic string) *AttributeBuilder {
	ab.attributes[string(semconv.MessagingDestinationNameKey)] = ProtocolTypeMQTT
	return ab
}

// WithProtocol sets the MQTT protocol version
func (ab *AttributeBuilder) WithProtocol(version string) *AttributeBuilder {
	ab.attributes[string(semconv.NetworkProtocolNameKey)] = ProtocolTypeMQTT
	ab.attributes[string(semconv.NetworkProtocolVersionKey)] = version
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
	ab.attributes[string(semconv.NetworkTransportKey)] = transport
	return ab
}

// WithServerAddress sets the server address
func (ab *AttributeBuilder) WithServerAddress(address string) *AttributeBuilder {
	ab.attributes[string(semconv.ServerAddressKey)] = address
	return ab
}

// WithServerPort sets the server port
func (ab *AttributeBuilder) WithServerPort(port int) *AttributeBuilder {
	if port > 0 {
		ab.attributes[string(semconv.ServerPortKey)] = port
	}
	return ab
}

// WithError sets error attributes
func (ab *AttributeBuilder) WithError(err error) *AttributeBuilder {
	if err != nil {
		ab.attributes[string(semconv.ErrorTypeKey)] = fmt.Sprintf("%T", err)
		ab.attributes[string(semconv.ErrorMessageKey)] = err.Error()
	}
	return ab
}

// Build returns the final attribute map
func (ab *AttributeBuilder) Build() map[string]interface{} {
	return ab.attributes
}

// Predefined attribute builders for common operations

// ConnectAttributes creates attributes for MQTT connect operations
func ConnectAttributes(clientID, brokerURL string, keepAlive int, cleanStart bool) []attribute.KeyValue {
	return []attribute.KeyValue{
		semconv.MessagingSystemKey.String(ProtocolTypeMQTT),
		semconv.MessagingOperationNameKey.String(OperationConnect),
		attribute.String(MQTTClientIDKey, clientID),
		attribute.String(MQTTBrokerURLKey, brokerURL),
		attribute.Int(MQTTKeepAliveKey, keepAlive),
		attribute.Bool(MQTTCleanStartKey, cleanStart),
		semconv.NetworkProtocolNameKey.String(ProtocolTypeMQTT),
	}
}

// PublishAttributes creates attributes for MQTT publish operations
func PublishAttributes(msg observability.PublishMessage) []attribute.KeyValue {
	return []attribute.KeyValue{
		semconv.MessagingSystemKey.String(ProtocolTypeMQTT),
		semconv.MessagingOperationNameKey.String(OperationPublish),
		semconv.MessagingDestinationNameKey.String(ProtocolTypeMQTT),
		semconv.NetworkProtocolNameKey.String(ProtocolTypeMQTT),
		attribute.Int(MessagingMQTTQoSKey, int(msg.QoS)),
		attribute.Bool(MQTTMessageRetainKey, msg.Retain),
		attribute.Int(MessagingMessageIDKey, int(msg.MessageID)),
		attribute.Int(MessagingMessageBodySizeKey, len(msg.Payload)),
	}
}

// SubscribeAttributes creates attributes for MQTT subscribe operations
func SubscribeAttributes(topics []observability.SubscriptionTopic) []attribute.KeyValue {
	attrs := []attribute.KeyValue{
		semconv.MessagingSystemKey.String(ProtocolTypeMQTT),
		semconv.MessagingOperationNameKey.String(OperationSubscribe),
		semconv.NetworkProtocolNameKey.String(ProtocolTypeMQTT),
	}

	// For multiple topics, use the first one as the primary destination
	if len(topics) > 0 {
		attrs = append(attrs,
			semconv.MessagingDestinationNameKey.String(ProtocolTypeMQTT),
			attribute.Int(MessagingMQTTQoSKey, int(topics[0].QoS)),
		)
	}

	return attrs
}

// UnsubscribeAttributes creates attributes for MQTT unsubscribe operations
func UnsubscribeAttributes(topics []string) []attribute.KeyValue {
	attrs := []attribute.KeyValue{
		semconv.MessagingSystemKey.String(ProtocolTypeMQTT),
		semconv.MessagingOperationNameKey.String(OperationUnsubscribe),
		semconv.NetworkProtocolNameKey.String(ProtocolTypeMQTT),
	}

	// For multiple topics, use the first one as the primary destination
	if len(topics) > 0 {
		attrs = append(attrs,
			semconv.MessagingDestinationNameKey.String(ProtocolTypeMQTT),
		)
	}

	return attrs
}

// ReceiveAttributes creates attributes for MQTT message receive operations
func ReceiveAttributes(msg observability.ReceivedMessage) []attribute.KeyValue {
	attrs := []attribute.KeyValue{
		semconv.MessagingSystemKey.String(ProtocolTypeMQTT),
		semconv.MessagingOperationNameKey.String(OperationReceive),
		semconv.MessagingDestinationNameKey.String(ProtocolTypeMQTT),
		semconv.NetworkProtocolNameKey.String(ProtocolTypeMQTT),
		attribute.Int(MessagingMQTTQoSKey, int(msg.QoS)),
		attribute.Bool(MQTTMessageRetainKey, msg.Retain),
		attribute.Bool(MQTTMessageDuplicateKey, msg.Duplicate),
		attribute.Int(MessagingMessageIDKey, int(msg.MessageID)),
		attribute.Int(MessagingMessageBodySizeKey, len(msg.Payload)),
	}

	return attrs
}

// DisconnectAttributes creates attributes for MQTT disconnect operations
func DisconnectAttributes(reason observability.DisconnectReason) []attribute.KeyValue {
	return []attribute.KeyValue{
		semconv.MessagingSystemKey.String(ProtocolTypeMQTT),
		semconv.MessagingOperationNameKey.String(OperationDisconnect),
		semconv.NetworkProtocolNameKey.String(ProtocolTypeMQTT),
		attribute.Int(MQTTReasonCodeKey, int(reason.Code)),
		attribute.String(MQTTReasonStringKey, reason.Reason),
		attribute.Bool(MQTTDisconnectInitiatedByClientKey, reason.InitiatedByClient),
	}
}

// ReconnectAttributes creates attributes for MQTT reconnect operations
func ReconnectAttributes(attempt int, backoff time.Duration) []attribute.KeyValue {
	return []attribute.KeyValue{
		semconv.MessagingSystemKey.String(ProtocolTypeMQTT),
		semconv.MessagingOperationNameKey.String(OperationReconnect),
		semconv.NetworkProtocolNameKey.String(ProtocolTypeMQTT),
		attribute.Int(MQTTReconnectAttemptKey, attempt),
		attribute.String(MQTTReconnectBackoffKey, backoff.String()),
	}
}

// PingAttributes creates attributes for MQTT ping operations
func PingAttributes(clientID string) []attribute.KeyValue {
	return []attribute.KeyValue{
		semconv.MessagingSystemKey.String(ProtocolTypeMQTT),
		semconv.MessagingOperationNameKey.String(OperationPing),
		semconv.NetworkProtocolNameKey.String(ProtocolTypeMQTT),
		attribute.String(MQTTClientIDKey, clientID),
	}
}

// AuthAttributes creates attributes for MQTT authentication operations
func AuthAttributes(clientID, method string) []attribute.KeyValue {
	return []attribute.KeyValue{
		semconv.MessagingSystemKey.String(ProtocolTypeMQTT),
		semconv.MessagingOperationNameKey.String(OperationAuth),
		semconv.NetworkProtocolNameKey.String(ProtocolTypeMQTT),
		attribute.String(MQTTClientIDKey, clientID),
		attribute.String(MQTTAuthMethodKey, method),
	}
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
