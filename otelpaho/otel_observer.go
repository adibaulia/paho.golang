// Package otel provides OpenTelemetry implementation for MQTT observability
package otelpaho

import (
	"context"
	"fmt"
	"time"

	"github.com/adibaulia/paho.golang/observability"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

// OTelMQTTObserver implements observability.MQTTObserver using OpenTelemetry
type OTelMQTTObserver struct {
	tracer         trace.Tracer
	meter          metric.Meter
	clientID       string
	metricRecorder *MetricRecorder
	spanBuilder    SpanNameBuilder
}

// NewOTelMQTTObserver creates a new OpenTelemetry MQTT observer
func NewOTelMQTTObserver(clientID string) *OTelMQTTObserver {
	return &OTelMQTTObserver{
		tracer:         otel.Tracer("mqtt-client"),
		meter:          otel.Meter("mqtt-client"),
		clientID:       clientID,
		metricRecorder: NewMetricRecorder(),
		spanBuilder:    SpanNames,
	}
}

// OnConnect is called when a connection is established
func (o *OTelMQTTObserver) OnConnect(ctx context.Context, config observability.ConnectionConfig) (context.Context, func(error)) {
	// Start span for connection using semantic convention naming
	spanName := MessagingClientOperationSpan
	ctx, span := o.tracer.Start(ctx, spanName)

	// Build attributes using semantic conventions
	attrs := ConnectAttributes(config.ClientID, config.BrokerURL, int(config.KeepAlive.Seconds()), config.CleanStart)
	span.SetAttributes(attrs...)

	// Record connection attempt metric
	if o.metricRecorder != nil {
		o.metricRecorder.RecordConnectionAttempt(
			config.ClientID,
			config.BrokerURL,
			TransportTCP, // Default, could be determined from URL
		)
	}

	// Store connection start time for latency calculation
	start := time.Now()
	ctx = context.WithValue(ctx, MQTTConnectStartKey, start)

	fmt.Printf("MQTT Connect: %s to %s\n", config.ClientID, config.BrokerURL)
	fmt.Printf("Attributes: %+v\n", attrs)

	// Return completion callback
	return ctx, func(err error) {
		defer span.End()
		if err != nil {
			o.RecordError(ctx, OperationConnect, err)
			span.RecordError(err)
		} else {
			// Record successful connection duration
			if startTime, ok := ctx.Value(MQTTConnectStartKey).(time.Time); ok {
				duration := time.Since(startTime)
				if o.metricRecorder != nil {
					o.metricRecorder.RecordConnectionSuccess(o.clientID, config.BrokerURL, TransportTCP, duration)
				}
			}
			fmt.Printf("MQTT Connect Complete\n")
		}
	}
}

// OnDisconnect is called when a connection is terminated
func (o *OTelMQTTObserver) OnDisconnect(ctx context.Context, reason observability.DisconnectReason) (context.Context, func(error)) {
	// Start span for disconnect using semantic convention naming
	spanName := MessagingClientOperationSpan
	ctx, span := o.tracer.Start(ctx, spanName)

	// Build attributes using semantic conventions
	attrs := DisconnectAttributes(reason)
	span.SetAttributes(attrs...)

	fmt.Printf("MQTT Disconnect: reason=%s, code=%d\n", reason.Reason, reason.Code)
	fmt.Printf("Attributes: %+v\n", attrs)

	return ctx, func(err error) {
		defer span.End()
		if err != nil {
			o.RecordError(ctx, OperationDisconnect, err)
			span.RecordError(err)
		} else {
			fmt.Printf("MQTT Disconnect Complete\n")
		}
	}
}

// OnReconnect is called when a reconnection occurs
func (o *OTelMQTTObserver) OnReconnect(ctx context.Context, attempt int, backoff time.Duration) (context.Context, func(error)) {
	// Start span for reconnect
	ctx, span := o.tracer.Start(ctx, MessagingClientOperationSpan)

	// Build attributes
	attrs := ReconnectAttributes(attempt, backoff)
	span.SetAttributes(attrs...)

	// Record reconnection attempt
	o.metricRecorder.RecordError(o.clientID, OperationReconnect, "reconnection_needed")

	fmt.Printf("MQTT Reconnect: attempt=%d, backoff=%v\n", attempt, backoff)
	fmt.Printf("Attributes: %+v\n", attrs)

	return ctx, func(err error) {
		defer span.End()
		if err != nil {
			o.RecordError(ctx, OperationReconnect, err)
		} else {
			fmt.Printf("MQTT Reconnect Complete\n")
		}
	}
}

// OnPublish is called when a message is published
func (o *OTelMQTTObserver) OnPublish(ctx context.Context, msg observability.PublishMessage) (context.Context, func(observability.PublishResult)) {
	// Start span for publish
	ctx, span := o.tracer.Start(ctx, MessagingClientOperationSpan)

	// Build attributes using the attribute builder
	attrs := PublishAttributes(msg)
	span.SetAttributes(attrs...)

	// Record publish metrics
	start := time.Now()
	ctx = context.WithValue(ctx, MQTTPublishStartKey, start)

	fmt.Printf("MQTT Publish: topic=%s, qos=%d, size=%d bytes\n",
		msg.Topic, msg.QoS, len(msg.Payload))
	fmt.Printf("Attributes: %+v\n", attrs)

	return ctx, func(result observability.PublishResult) {
		defer span.End()
		if result.Error != nil {
			o.RecordError(ctx, OperationPublish, result.Error)
		} else {
			fmt.Printf("MQTT Publish Complete\n")
		}
	}
}

// OnPublishComplete is called when a publish operation completes
func (o *OTelMQTTObserver) OnPublishComplete(ctx context.Context, msg observability.PublishMessage, err error) {
	// Calculate latency if start time is available
	if startTime, ok := ctx.Value(MQTTPublishStartKey).(time.Time); ok {
		latency := time.Since(startTime)

		if o.metricRecorder != nil {
			if err != nil {
				o.metricRecorder.RecordError(o.clientID, OperationPublish, fmt.Sprintf("%T", err))
			} else {
				o.metricRecorder.RecordMessagePublished(
					o.clientID,
					msg.Topic,
					int(msg.QoS),
					len(msg.Payload),
					latency,
				)
			}
		}

		if err != nil {
			fmt.Printf("MQTT Publish Failed: %v, latency=%v\n", err, latency)
		} else {
			fmt.Printf("MQTT Publish Complete: latency=%v\n", latency)
		}
	}
}

// OnSubscribe is called when subscribing to topics
func (o *OTelMQTTObserver) OnSubscribe(ctx context.Context, topics []observability.SubscriptionTopic) (context.Context, func(observability.SubscribeResult)) {
	// Start span for subscribe
	ctx, span := o.tracer.Start(ctx, MessagingClientOperationSpan)

	// Build attributes
	attrs := SubscribeAttributes(topics)
	span.SetAttributes(attrs...)

	// Record subscription attempts
	for _, topic := range topics {
		o.metricRecorder.RecordSubscription(o.clientID, topic.Topic, int(topic.QoS), true)
		fmt.Printf("MQTT Subscribe: topic=%s, qos=%d\n", topic.Topic, topic.QoS)
	}
	fmt.Printf("Attributes: %+v\n", attrs)

	return ctx, func(result observability.SubscribeResult) {
		defer span.End()
		if result.Error != nil {
			o.RecordError(ctx, OperationSubscribe, result.Error)
		} else {
			fmt.Printf("MQTT Subscribe Complete\n")
		}
	}
}

// OnUnsubscribe is called when unsubscribing from topics
func (o *OTelMQTTObserver) OnUnsubscribe(ctx context.Context, topics []string) (context.Context, func(observability.UnsubscribeResult)) {
	// Start span for unsubscribe
	ctx, span := o.tracer.Start(ctx, MessagingClientOperationSpan)

	// Build attributes
	attrs := UnsubscribeAttributes(topics)
	span.SetAttributes(attrs...)

	for _, topic := range topics {
		fmt.Printf("MQTT Unsubscribe: topic=%s\n", topic)
	}
	fmt.Printf("Attributes: %+v\n", attrs)

	return ctx, func(result observability.UnsubscribeResult) {
		defer span.End()
		if result.Error != nil {
			o.RecordError(ctx, OperationUnsubscribe, result.Error)
		} else {
			fmt.Printf("MQTT Unsubscribe Complete\n")
		}
	}
}

// OnMessageReceived is called when a message is received
func (o *OTelMQTTObserver) OnMessageReceived(ctx context.Context, msg observability.ReceivedMessage) (context.Context, func(error)) {
	// Start span for message receive
	ctx, span := o.tracer.Start(ctx, MessagingClientOperationSpan)
	defer span.End()

	// Build attributes
	attrs := ReceiveAttributes(msg)
	span.SetAttributes(attrs...)

	// Record message received metrics
	o.metricRecorder.RecordMessageReceived(
		o.clientID,
		msg.Topic,
		int(msg.QoS),
		len(msg.Payload),
	)

	fmt.Printf("MQTT Message Received: topic=%s, qos=%d, size=%d bytes\n",
		msg.Topic, msg.QoS, len(msg.Payload))
	fmt.Printf("Attributes: %+v\n", attrs)

	return ctx, func(err error) {
		if err != nil {
			o.RecordError(ctx, "message_received", err)
		} else {
			fmt.Printf("[OTEL] Message processing completed successfully\n")
		}
	}
}

// OnPing is called when a ping is sent
func (o *OTelMQTTObserver) OnPing(ctx context.Context) (context.Context, func(error)) {
	// Start span for ping
	ctx, span := o.tracer.Start(ctx, MessagingClientOperationSpan)
	defer span.End()

	// Build attributes
	attrs := PingAttributes(o.clientID)
	span.SetAttributes(attrs...)

	// Store ping start time for latency calculation
	ctx = context.WithValue(ctx, MQTTPingStartKey, time.Now())

	fmt.Printf("MQTT Ping sent\n")
	fmt.Printf("Attributes: %+v\n", attrs)

	return ctx, func(err error) {
		if err != nil {
			o.RecordError(ctx, "ping", err)
		} else {
			fmt.Printf("[OTEL] Ping completed successfully\n")
		}
	}
}

// OnPingResponse is called when a ping response is received
func (o *OTelMQTTObserver) OnPingResponse(ctx context.Context) {
	// Calculate ping latency if start time is available
	if startTime, ok := ctx.Value(MQTTPingStartKey).(time.Time); ok {
		latency := time.Since(startTime)

		// Record ping metrics
		o.metricRecorder.RecordPing(o.clientID, "", latency) // Broker URL would be stored in observer

		fmt.Printf("MQTT Ping Response: latency=%v\n", latency)
	}
}

// OnAuth is called during authentication
func (o *OTelMQTTObserver) OnAuth(ctx context.Context, authData observability.AuthData) (context.Context, func(observability.AuthResult)) {
	// Start span for auth
	ctx, span := o.tracer.Start(ctx, MessagingClientOperationSpan)
	defer span.End()

	// Build attributes
	attrs := AuthAttributes(o.clientID, authData.Method)
	span.SetAttributes(attrs...)

	fmt.Printf("MQTT Auth: method=%s\n", authData.Method)
	fmt.Printf("Attributes: %+v\n", attrs)

	// Return completion callback
	return ctx, func(result observability.AuthResult) {
		if result.Error != nil {
			o.RecordError(ctx, OperationAuth, result.Error)
		} else {
			fmt.Printf("MQTT Auth Complete: success=%v, latency=%v\n", result.Success, result.Latency)
		}
	}
}

// RecordError records an error event
func (o *OTelMQTTObserver) RecordError(ctx context.Context, operation string, err error) {
	// Record error metric and span event
	errorType := fmt.Sprintf("%T", err)
	if o.metricRecorder != nil {
		o.metricRecorder.RecordError(o.clientID, operation, errorType)
	}

	fmt.Printf("MQTT Error: operation=%s, type=%s, message=%s\n",
		operation, errorType, err.Error())
}

// RecordLatency records operation latency
func (o *OTelMQTTObserver) RecordLatency(ctx context.Context, operation string, duration time.Duration) {
	// Record latency metric
	fmt.Printf("MQTT Latency: operation=%s, duration=%v\n", operation, duration)
}

// RecordMessageSize records message size
func (o *OTelMQTTObserver) RecordMessageSize(ctx context.Context, operation string, size int64) {
	// Record message size metric
	fmt.Printf("MQTT Message Size: operation=%s, size=%d bytes\n", operation, size)
}

// Helper function to convert attribute map to OpenTelemetry attributes
func convertToOTelAttributes(attrs map[string]interface{}) []attribute.KeyValue {
	var otelAttrs []attribute.KeyValue
	for key, value := range attrs {
		switch v := value.(type) {
		case string:
			otelAttrs = append(otelAttrs, attribute.String(key, v))
		case int:
			otelAttrs = append(otelAttrs, attribute.Int(key, v))
		case bool:
			otelAttrs = append(otelAttrs, attribute.Bool(key, v))
		case float64:
			otelAttrs = append(otelAttrs, attribute.Float64(key, v))
		default:
			// Convert other types to string as fallback
			otelAttrs = append(otelAttrs, attribute.String(key, fmt.Sprintf("%v", v)))
		}
	}
	return otelAttrs
}

// Example of how to create and configure the observer with actual OpenTelemetry
// This would be in a separate file or example when dependencies are available:

/*
func NewConfiguredOTelMQTTObserver(clientID string) (*OTelMQTTObserver, error) {
	// Initialize OpenTelemetry provider
	tp, err := trace.NewTracerProvider(
		trace.WithBatcher(trace.NewBatchSpanProcessor(
			otlptracegrpc.NewClient(),
		)),
		trace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName("mqtt-client"),
			semconv.ServiceVersion("1.0.0"),
		)),
	)
	if err != nil {
		return nil, err
	}

	mp, err := metric.NewMeterProvider(
		metric.WithReader(metric.NewPeriodicReader(
			otlpmetricgrpc.NewClient(),
		)),
		metric.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName("mqtt-client"),
			semconv.ServiceVersion("1.0.0"),
		)),
	)
	if err != nil {
		return nil, err
	}

	// Set global providers
	otel.SetTracerProvider(tp)
	otel.SetMeterProvider(mp)

	// Create observer
	return &OTelMQTTObserver{
		// tracer:         tp.Tracer("mqtt-client"),
		// meter:          mp.Meter("mqtt-client"),
		clientID:       clientID,
		metricRecorder: NewMetricRecorder(),
		spanBuilder:    SpanNames,
	}, nil
}
*/
