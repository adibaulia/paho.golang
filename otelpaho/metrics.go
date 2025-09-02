// Package otel provides metric definitions for MQTT observability
// following OpenTelemetry semantic conventions for messaging systems.
package otelpaho

import (
	"fmt"
	"strings"
	"time"

	"go.opentelemetry.io/otel/semconv/v1.37.0"
)

// Messaging Metric Names following OpenTelemetry semantic conventions
// https://opentelemetry.io/docs/specs/semconv/messaging/messaging-metrics/

const (
	// Client Operation Duration - measures the duration of messaging operations
	MessagingClientOperationDurationMetric = "messaging.client.operation.duration" // Histogram: Duration of messaging operations

	// Published Messages - measures the number of messages published
	MessagingClientPublishedMessagesMetric = "messaging.client.published.messages" // Counter: Number of messages published

	// Consumed Messages - measures the number of messages consumed
	MessagingClientConsumedMessagesMetric = "messaging.client.consumed.messages" // Counter: Number of messages consumed

	// MQTT-specific metrics (not covered by standard semantic conventions)
	MQTTConnectionsActiveMetric    = "mqtt.client.connections.active"    // Gauge: Number of active MQTT connections
	MQTTConnectionAttemptsMetric   = "mqtt.client.connection.attempts"   // Counter: Total connection attempts
	MQTTConnectionSuccessMetric    = "mqtt.client.connection.success"    // Counter: Successful connections
	MQTTConnectionFailuresMetric   = "mqtt.client.connection.failures"   // Counter: Failed connections
	MQTTReconnectionAttemptsMetric = "mqtt.client.reconnection.attempts" // Counter: Reconnection attempts

	// Message Size Metrics
	MQTTPublishSizeMetric = "mqtt.client.publish.size" // Histogram: Size of published messages
	MQTTMessageSizeMetric = "mqtt.client.message.size" // Histogram: Size of received messages

	// Queue Metrics
	MQTTPublishQueueSizeMetric = "mqtt.client.publish.queue_size" // Gauge: Number of messages in publish queue

	// Subscription Metrics
	MQTTSubscriptionsActiveMetric = "mqtt.client.subscriptions.active" // Gauge: Number of active subscriptions
	MQTTSubscribeAttemptsMetric   = "mqtt.client.subscribe.attempts"   // Counter: Subscribe attempts
	MQTTSubscribeSuccessMetric    = "mqtt.client.subscribe.success"    // Counter: Successful subscriptions
	MQTTSubscribeFailuresMetric   = "mqtt.client.subscribe.failures"   // Counter: Failed subscriptions
	MQTTUnsubscribeAttemptsMetric = "mqtt.client.unsubscribe.attempts" // Counter: Unsubscribe attempts
	MQTTUnsubscribeSuccessMetric  = "mqtt.client.unsubscribe.success"  // Counter: Successful unsubscriptions

	// Protocol Metrics
	MQTTPingRequestsMetric    = "mqtt.client.ping.requests"    // Counter: PING requests sent
	MQTTPingResponsesMetric   = "mqtt.client.ping.responses"   // Counter: PING responses received
	MQTTPingLatencyMetric     = "mqtt.client.ping.latency"     // Histogram: PING round-trip time
	MQTTPacketsSentMetric     = "mqtt.client.packets.sent"     // Counter: Total packets sent
	MQTTPacketsReceivedMetric = "mqtt.client.packets.received" // Counter: Total packets received
	MQTTPacketSizeMetric      = "mqtt.client.packet.size"      // Histogram: Packet sizes

	// Error Metrics
	MQTTErrorsMetric         = "mqtt.client.errors"          // Counter: Total errors by type
	MQTTProtocolErrorsMetric = "mqtt.client.protocol.errors" // Counter: Protocol-level errors
	MQTTNetworkErrorsMetric  = "mqtt.client.network.errors"  // Counter: Network-level errors
	MQTTTimeoutsMetric       = "mqtt.client.timeouts"        // Counter: Operation timeouts

	// Session Metrics
	MQTTSessionsActiveMetric  = "mqtt.client.sessions.active"  // Gauge: Active sessions
	MQTTSessionDurationMetric = "mqtt.client.session.duration" // Histogram: Session duration
	MQTTSessionMessagesMetric = "mqtt.client.session.messages" // Counter: Messages per session

	// Throughput Metrics
	MQTTThroughputBytesMetric    = "mqtt.client.throughput.bytes"    // Counter: Total bytes transferred
	MQTTThroughputMessagesMetric = "mqtt.client.throughput.messages" // Counter: Total messages transferred
)

// Metric Units following OpenTelemetry conventions
const (
	UnitBytes         = "By"             // Bytes
	UnitMilliseconds  = "ms"             // Milliseconds
	UnitSeconds       = "s"              // Seconds
	UnitCount         = "1"              // Dimensionless count
	UnitConnections   = "{connection}"   // Connection count
	UnitMessages      = "{message}"      // Message count
	UnitPackets       = "{packet}"       // Packet count
	UnitSubscriptions = "{subscription}" // Subscription count
)

// MetricDefinition represents a metric with its metadata
type MetricDefinition struct {
	Name        string
	Description string
	Unit        string
	Type        MetricType
	Attributes  []string // Common attribute keys for this metric
}

// MetricType represents the type of metric
type MetricType string

const (
	MetricTypeCounter   MetricType = "counter"
	MetricTypeGauge     MetricType = "gauge"
	MetricTypeHistogram MetricType = "histogram"
)

// MQTT Metric Definitions
var MQTTMetrics = map[string]MetricDefinition{
	// Connection Metrics
	MQTTConnectionsActiveMetric: {
		Name:        MQTTConnectionsActiveMetric,
		Description: "Number of active MQTT connections",
		Unit:        UnitConnections,
		Type:        MetricTypeGauge,
		Attributes:  []string{MQTTClientIDKey, MQTTBrokerURLKey, string(semconv.NetworkTransportKey)},
	},
	MQTTConnectionAttemptsMetric: {
		Name:        MQTTConnectionAttemptsMetric,
		Description: "Total number of MQTT connection attempts",
		Unit:        UnitCount,
		Type:        MetricTypeCounter,
		Attributes:  []string{MQTTClientIDKey, MQTTBrokerURLKey, string(semconv.NetworkTransportKey)},
	},
	MQTTConnectionSuccessMetric: {
		Name:        MQTTConnectionSuccessMetric,
		Description: "Total number of successful MQTT connections",
		Unit:        UnitCount,
		Type:        MetricTypeCounter,
		Attributes:  []string{MQTTClientIDKey, MQTTBrokerURLKey, string(semconv.NetworkTransportKey)},
	},
	MQTTConnectionFailuresMetric: {
		Name:        MQTTConnectionFailuresMetric,
		Description: "Total number of failed MQTT connection attempts",
		Unit:        UnitCount,
		Type:        MetricTypeCounter,
		Attributes:  []string{MQTTClientIDKey, MQTTBrokerURLKey, string(semconv.NetworkTransportKey), string(semconv.ErrorTypeKey)},
	},


	// Publishing Metrics
	MessagingClientPublishedMessagesMetric: {
		Name:        MessagingClientPublishedMessagesMetric,
		Description: "Number of messages published by the client",
		Unit:        UnitMessages,
		Type:        MetricTypeCounter,
		Attributes:  []string{MQTTClientIDKey, string(semconv.MessagingDestinationNameKey), MQTTQoSKey},
	},
	MessagingClientConsumedMessagesMetric: {
		Name:        MessagingClientConsumedMessagesMetric,
		Description: "Number of messages consumed by the client",
		Unit:        UnitMessages,
		Type:        MetricTypeCounter,
		Attributes:  []string{MQTTClientIDKey, string(semconv.MessagingDestinationNameKey), MQTTQoSKey},
	},
	MessagingClientOperationDurationMetric: {
		Name:        MessagingClientOperationDurationMetric,
		Description: "Duration of messaging operations",
		Unit:        UnitMilliseconds,
		Type:        MetricTypeHistogram,
		Attributes:  []string{MQTTClientIDKey, MQTTBrokerURLKey, string(semconv.MessagingOperationNameKey)},
	},
	MQTTPublishSizeMetric: {
		Name:        MQTTPublishSizeMetric,
		Description: "Size of published MQTT messages",
		Unit:        UnitBytes,
		Type:        MetricTypeHistogram,
		Attributes:  []string{MQTTClientIDKey, string(semconv.MessagingDestinationNameKey), MQTTQoSKey},
	},

	// Message Size Metrics
	MQTTMessageSizeMetric: {
		Name:        MQTTMessageSizeMetric,
		Description: "Size of received MQTT messages",
		Unit:        UnitBytes,
		Type:        MetricTypeHistogram,
		Attributes:  []string{MQTTClientIDKey, string(semconv.MessagingDestinationNameKey), MQTTQoSKey},
	},

	// Subscription Metrics
	MQTTSubscriptionsActiveMetric: {
		Name:        MQTTSubscriptionsActiveMetric,
		Description: "Number of active MQTT subscriptions",
		Unit:        UnitSubscriptions,
		Type:        MetricTypeGauge,
		Attributes:  []string{MQTTClientIDKey, string(semconv.MessagingDestinationNameKey), MQTTQoSKey},
	},
	MQTTSubscribeAttemptsMetric: {
		Name:        MQTTSubscribeAttemptsMetric,
		Description: "Total number of MQTT subscribe attempts",
		Unit:        UnitCount,
		Type:        MetricTypeCounter,
		Attributes:  []string{MQTTClientIDKey, string(semconv.MessagingDestinationNameKey), MQTTQoSKey},
	},

	// Protocol Metrics
	MQTTPingRequestsMetric: {
		Name:        MQTTPingRequestsMetric,
		Description: "Total number of MQTT PING requests sent",
		Unit:        UnitCount,
		Type:        MetricTypeCounter,
		Attributes:  []string{MQTTClientIDKey, MQTTBrokerURLKey},
	},
	MQTTPingLatencyMetric: {
		Name:        MQTTPingLatencyMetric,
		Description: "MQTT PING round-trip time",
		Unit:        UnitMilliseconds,
		Type:        MetricTypeHistogram,
		Attributes:  []string{MQTTClientIDKey, MQTTBrokerURLKey},
	},

	// Error Metrics
	MQTTErrorsMetric: {
		Name:        MQTTErrorsMetric,
		Description: "Total number of MQTT errors by type",
		Unit:        UnitCount,
		Type:        MetricTypeCounter,
		Attributes:  []string{MQTTClientIDKey, string(semconv.ErrorTypeKey), string(semconv.MessagingOperationNameKey)},
	},
	MQTTProtocolErrorsMetric: {
		Name:        MQTTProtocolErrorsMetric,
		Description: "Total number of MQTT protocol errors",
		Unit:        UnitCount,
		Type:        MetricTypeCounter,
		Attributes:  []string{MQTTClientIDKey, MQTTReasonCodeKey, string(semconv.ErrorTypeKey)},
	},
}

// MetricRecorder provides methods to record MQTT metrics
// This implementation uses a map-based approach for demonstration
// In production, this would integrate with actual OpenTelemetry metric instruments
type MetricRecorder struct {
	counters   map[string]int64
	histograms map[string][]float64
	gauges     map[string]int64
	lastUpdate map[string]time.Time
}

// NewMetricRecorder creates a new metric recorder
func NewMetricRecorder() *MetricRecorder {
	return &MetricRecorder{
		counters:   make(map[string]int64),
		histograms: make(map[string][]float64),
		gauges:     make(map[string]int64),
		lastUpdate: make(map[string]time.Time),
	}
}

// ErrorCategory represents different types of MQTT errors
type ErrorCategory string

const (
	ErrorCategoryConnection ErrorCategory = "connection"
	ErrorCategoryProtocol   ErrorCategory = "protocol"
	ErrorCategoryNetwork    ErrorCategory = "network"
	ErrorCategoryTimeout    ErrorCategory = "timeout"
	ErrorCategoryAuth       ErrorCategory = "authentication"
	ErrorCategoryPermission ErrorCategory = "permission"
	ErrorCategoryQoS        ErrorCategory = "qos"
	ErrorCategoryPayload    ErrorCategory = "payload"
	ErrorCategoryUnknown    ErrorCategory = "unknown"
)

// CategorizeError categorizes an error based on its type and message
func (mr *MetricRecorder) CategorizeError(err error) ErrorCategory {
	if err == nil {
		return ErrorCategoryUnknown
	}

	errorMsg := err.Error()

	// Network-related errors
	if contains(errorMsg, []string{"connection refused", "network unreachable", "no route to host", "connection reset"}) {
		return ErrorCategoryNetwork
	}

	// Timeout errors
	if contains(errorMsg, []string{"timeout", "deadline exceeded", "context deadline"}) {
		return ErrorCategoryTimeout
	}

	// Authentication errors
	if contains(errorMsg, []string{"authentication", "unauthorized", "bad username", "bad password"}) {
		return ErrorCategoryAuth
	}

	// Permission errors
	if contains(errorMsg, []string{"permission", "forbidden", "not authorized", "access denied"}) {
		return ErrorCategoryPermission
	}

	// Protocol errors
	if contains(errorMsg, []string{"protocol", "malformed", "invalid packet", "reason code"}) {
		return ErrorCategoryProtocol
	}

	// QoS errors
	if contains(errorMsg, []string{"qos", "quality of service", "delivery"}) {
		return ErrorCategoryQoS
	}

	// Payload errors
	if contains(errorMsg, []string{"payload", "message size", "too large", "encoding"}) {
		return ErrorCategoryPayload
	}

	// Connection errors
	if contains(errorMsg, []string{"connection", "connect", "disconnect", "broker"}) {
		return ErrorCategoryConnection
	}

	return ErrorCategoryUnknown
}

// contains checks if any of the keywords exist in the text
func contains(text string, keywords []string) bool {
	text = strings.ToLower(text)
	for _, keyword := range keywords {
		if strings.Contains(text, keyword) {
			return true
		}
	}
	return false
}

// RecordConnectionAttempt records a connection attempt
func (mr *MetricRecorder) RecordConnectionAttempt(clientID, brokerURL, transport string) {
	key := fmt.Sprintf("%s.%s.%s", MQTTConnectionAttemptsMetric, clientID, transport)
	mr.incrementCounter(key)
	fmt.Printf("[METRICS] Connection attempt: client=%s, broker=%s, transport=%s\n", clientID, brokerURL, transport)
}

// RecordConnectionSuccess records a successful connection
func (mr *MetricRecorder) RecordConnectionSuccess(clientID, brokerURL, transport string, latency time.Duration) {
	successKey := fmt.Sprintf("%s.%s.%s", MQTTConnectionSuccessMetric, clientID, transport)
	latencyKey := fmt.Sprintf("%s.%s.%s", MessagingClientOperationDurationMetric, clientID, "connect")
	activeKey := fmt.Sprintf("%s.%s", MQTTConnectionsActiveMetric, clientID)

	mr.incrementCounter(successKey)
	mr.recordHistogram(latencyKey, float64(latency.Milliseconds()))
	mr.incrementGauge(activeKey)

	fmt.Printf("[METRICS] Connection success: client=%s, broker=%s, latency=%v\n", clientID, brokerURL, latency)
}

// RecordConnectionFailure records a failed connection
func (mr *MetricRecorder) RecordConnectionFailure(clientID, brokerURL, transport string, err error) {
	errorCategory := mr.CategorizeError(err)
	key := fmt.Sprintf("%s.%s.%s.%s", MQTTConnectionFailuresMetric, clientID, transport, errorCategory)
	mr.incrementCounter(key)

	// Also record in general error metrics
	mr.RecordError(clientID, "connect", string(errorCategory))

	fmt.Printf("[METRICS] Connection failure: client=%s, broker=%s, error=%s, category=%s\n",
		clientID, brokerURL, err.Error(), errorCategory)
}

// RecordMessagePublished records a published message
func (mr *MetricRecorder) RecordMessagePublished(clientID, topic string, qos int, size int, latency time.Duration) {
	publishedKey := fmt.Sprintf("%s.%s.qos%d", MessagingClientPublishedMessagesMetric, clientID, qos)
	latencyKey := fmt.Sprintf("%s.%s.%s", MessagingClientOperationDurationMetric, clientID, "publish")
	sizeKey := fmt.Sprintf("%s.%s.qos%d", MQTTPublishSizeMetric, clientID, qos)

	mr.incrementCounter(publishedKey)
	mr.recordHistogram(latencyKey, float64(latency.Milliseconds()))
	mr.recordHistogram(sizeKey, float64(size))

	fmt.Printf("[METRICS] Message published: client=%s, topic=%s, qos=%d, size=%d, latency=%v\n",
		clientID, topic, qos, size, latency)
}

// RecordMessagePublishFailure records a failed message publish
func (mr *MetricRecorder) RecordMessagePublishFailure(clientID, topic string, qos int, err error) {
	errorCategory := mr.CategorizeError(err)
	// Record error in operation duration with error attribute
	latencyKey := fmt.Sprintf("%s.%s.%s.error", MessagingClientOperationDurationMetric, clientID, "publish")
	mr.recordHistogram(latencyKey, 0)

	// Also record in general error metrics
	mr.RecordError(clientID, "publish", string(errorCategory))

	fmt.Printf("[METRICS] Publish failure: client=%s, topic=%s, qos=%d, error=%s, category=%s\n",
		clientID, topic, qos, err.Error(), errorCategory)
}

// RecordMessageReceived records a received message
func (mr *MetricRecorder) RecordMessageReceived(clientID, topic string, qos int, size int) {
	receivedKey := fmt.Sprintf("%s.%s.qos%d", MessagingClientConsumedMessagesMetric, clientID, qos)
	sizeKey := fmt.Sprintf("%s.%s.qos%d", MQTTMessageSizeMetric, clientID, qos)

	mr.incrementCounter(receivedKey)
	mr.recordHistogram(sizeKey, float64(size))

	fmt.Printf("[METRICS] Message received: client=%s, topic=%s, qos=%d, size=%d\n",
		clientID, topic, qos, size)
}

// RecordSubscription records a subscription attempt and result
func (mr *MetricRecorder) RecordSubscription(clientID, topic string, qos int, success bool) {
	attemptsKey := fmt.Sprintf("%s.%s.qos%d", MQTTSubscribeAttemptsMetric, clientID, qos)
	mr.incrementCounter(attemptsKey)

	if success {
		successKey := fmt.Sprintf("%s.%s.qos%d", MQTTSubscribeSuccessMetric, clientID, qos)
		activeKey := fmt.Sprintf("%s.%s", MQTTSubscriptionsActiveMetric, clientID)
		mr.incrementCounter(successKey)
		mr.incrementGauge(activeKey)
		fmt.Printf("[METRICS] Subscription success: client=%s, topic=%s, qos=%d\n", clientID, topic, qos)
	} else {
		failureKey := fmt.Sprintf("%s.%s.qos%d", MQTTSubscribeFailuresMetric, clientID, qos)
		mr.incrementCounter(failureKey)
		fmt.Printf("[METRICS] Subscription failure: client=%s, topic=%s, qos=%d\n", clientID, topic, qos)
	}
}

// RecordUnsubscription records an unsubscription attempt
func (mr *MetricRecorder) RecordUnsubscription(clientID, topic string, success bool) {
	attemptsKey := fmt.Sprintf("%s.%s", MQTTUnsubscribeAttemptsMetric, clientID)
	mr.incrementCounter(attemptsKey)

	if success {
		successKey := fmt.Sprintf("%s.%s", MQTTUnsubscribeSuccessMetric, clientID)
		activeKey := fmt.Sprintf("%s.%s", MQTTSubscriptionsActiveMetric, clientID)
		mr.incrementCounter(successKey)
		mr.decrementGauge(activeKey)
		fmt.Printf("[METRICS] Unsubscription success: client=%s, topic=%s\n", clientID, topic)
	} else {
		fmt.Printf("[METRICS] Unsubscription failure: client=%s, topic=%s\n", clientID, topic)
	}
}

// RecordError records a general error
func (mr *MetricRecorder) RecordError(clientID, operation, errorType string) {
	generalKey := fmt.Sprintf("%s.%s.%s.%s", MQTTErrorsMetric, clientID, operation, errorType)
	mr.incrementCounter(generalKey)

	fmt.Printf("[METRICS] Error recorded: client=%s, operation=%s, type=%s\n", clientID, operation, errorType)
}

// RecordPing records a ping operation
func (mr *MetricRecorder) RecordPing(clientID, brokerURL string, latency time.Duration) {
	requestKey := fmt.Sprintf("%s.%s", MQTTPingRequestsMetric, clientID)
	latencyKey := fmt.Sprintf("%s.%s", MQTTPingLatencyMetric, clientID)

	mr.incrementCounter(requestKey)
	mr.recordHistogram(latencyKey, float64(latency.Milliseconds()))

	fmt.Printf("[METRICS] Ping recorded: client=%s, broker=%s, latency=%v\n", clientID, brokerURL, latency)
}

// RecordQueueSize records the current publish queue size
func (mr *MetricRecorder) RecordQueueSize(clientID string, size int) {
	key := fmt.Sprintf("%s.%s", MQTTPublishQueueSizeMetric, clientID)
	mr.setGauge(key, int64(size))

	fmt.Printf("[METRICS] Queue size: client=%s, size=%d\n", clientID, size)
}

// RecordThroughput records throughput metrics
func (mr *MetricRecorder) RecordThroughput(clientID string, bytes int64, messages int64) {
	bytesKey := fmt.Sprintf("%s.%s", MQTTThroughputBytesMetric, clientID)
	messagesKey := fmt.Sprintf("%s.%s", MQTTThroughputMessagesMetric, clientID)

	mr.addToCounter(bytesKey, bytes)
	mr.addToCounter(messagesKey, messages)

	fmt.Printf("[METRICS] Throughput: client=%s, bytes=%d, messages=%d\n", clientID, bytes, messages)
}

// Helper methods for metric recording
func (mr *MetricRecorder) incrementCounter(key string) {
	if mr.counters == nil {
		mr.counters = make(map[string]int64)
	}
	if mr.lastUpdate == nil {
		mr.lastUpdate = make(map[string]time.Time)
	}
	mr.counters[key]++
	mr.lastUpdate[key] = time.Now()
}

func (mr *MetricRecorder) addToCounter(key string, value int64) {
	if mr.counters == nil {
		mr.counters = make(map[string]int64)
	}
	if mr.lastUpdate == nil {
		mr.lastUpdate = make(map[string]time.Time)
	}
	mr.counters[key] += value
	mr.lastUpdate[key] = time.Now()
}

func (mr *MetricRecorder) recordHistogram(key string, value float64) {
	if mr.histograms == nil {
		mr.histograms = make(map[string][]float64)
	}
	if mr.lastUpdate == nil {
		mr.lastUpdate = make(map[string]time.Time)
	}
	if mr.histograms[key] == nil {
		mr.histograms[key] = make([]float64, 0)
	}
	mr.histograms[key] = append(mr.histograms[key], value)
	mr.lastUpdate[key] = time.Now()
}

func (mr *MetricRecorder) incrementGauge(key string) {
	if mr.gauges == nil {
		mr.gauges = make(map[string]int64)
	}
	if mr.lastUpdate == nil {
		mr.lastUpdate = make(map[string]time.Time)
	}
	mr.gauges[key]++
	mr.lastUpdate[key] = time.Now()
}

func (mr *MetricRecorder) decrementGauge(key string) {
	if mr.gauges == nil {
		mr.gauges = make(map[string]int64)
	}
	if mr.lastUpdate == nil {
		mr.lastUpdate = make(map[string]time.Time)
	}
	if mr.gauges[key] > 0 {
		mr.gauges[key]--
	}
	mr.lastUpdate[key] = time.Now()
}

func (mr *MetricRecorder) setGauge(key string, value int64) {
	if mr.gauges == nil {
		mr.gauges = make(map[string]int64)
	}
	if mr.lastUpdate == nil {
		mr.lastUpdate = make(map[string]time.Time)
	}
	mr.gauges[key] = value
	mr.lastUpdate[key] = time.Now()
}

// GetMetricsSummary returns a summary of all recorded metrics
func (mr *MetricRecorder) GetMetricsSummary() map[string]interface{} {
	summary := make(map[string]interface{})

	// Safely handle nil maps
	if mr.counters != nil {
		summary["counters"] = mr.counters
	} else {
		summary["counters"] = make(map[string]int64)
	}

	if mr.gauges != nil {
		summary["gauges"] = mr.gauges
	} else {
		summary["gauges"] = make(map[string]int64)
	}

	// Calculate histogram statistics
	histogramStats := make(map[string]map[string]float64)
	if mr.histograms != nil {
		for key, values := range mr.histograms {
			if len(values) > 0 {
				stats := calculateHistogramStats(values)
				histogramStats[key] = stats
			}
		}
	}
	summary["histograms"] = histogramStats

	return summary
}

// calculateHistogramStats calculates basic statistics for histogram values
func calculateHistogramStats(values []float64) map[string]float64 {
	if len(values) == 0 {
		return map[string]float64{}
	}

	stats := make(map[string]float64)

	// Calculate sum and count
	var sum float64
	for _, v := range values {
		sum += v
	}
	stats["count"] = float64(len(values))
	stats["sum"] = sum
	stats["avg"] = sum / float64(len(values))

	// Find min and max
	min, max := values[0], values[0]
	for _, v := range values {
		if v < min {
			min = v
		}
		if v > max {
			max = v
		}
	}
	stats["min"] = min
	stats["max"] = max

	return stats
}

// HistogramBuckets defines standard histogram buckets for different metric types
var HistogramBuckets = struct {
	Latency  []float64 // Latency buckets in milliseconds
	Size     []float64 // Size buckets in bytes
	Duration []float64 // Duration buckets in seconds
}{
	Latency:  []float64{0.1, 0.5, 1, 2, 5, 10, 25, 50, 100, 250, 500, 1000, 2500, 5000, 10000},
	Size:     []float64{64, 256, 1024, 4096, 16384, 65536, 262144, 1048576, 4194304, 16777216},
	Duration: []float64{0.1, 0.5, 1, 5, 10, 30, 60, 300, 600, 1800, 3600},
}
