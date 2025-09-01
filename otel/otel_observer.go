// Package otel provides OpenTelemetry implementation for MQTT observability
package otel

import (
	"context"
	"fmt"
	"time"

	"github.com/eclipse/paho.golang/observability"
)

// OTelMQTTObserver implements observability.MQTTObserver using OpenTelemetry
type OTelMQTTObserver struct {
	// tracer trace.Tracer
	// meter metric.Meter
	clientID string
	metricRecorder *MetricRecorder
	spanBuilder SpanNameBuilder
}

// NewOTelMQTTObserver creates a new OpenTelemetry MQTT observer
func NewOTelMQTTObserver(clientID string) *OTelMQTTObserver {
	return &OTelMQTTObserver{
		clientID: clientID,
		metricRecorder: &MetricRecorder{},
		spanBuilder: SpanNames,
	}
}

// OnConnect is called when a connection is established
func (o *OTelMQTTObserver) OnConnect(ctx context.Context, config observability.ConnectionConfig) (context.Context, func(error)) {
	// Start span for connection
	spanName := o.spanBuilder.ForOperation(OperationConnect)
	// ctx, span := o.tracer.Start(ctx, spanName)
	// defer span.End()
	_ = spanName // Placeholder until OpenTelemetry is integrated
	
	// Build attributes using the attribute builder
	attrs := ConnectAttributes(config)
	// span.SetAttributes(convertToOTelAttributes(attrs)...)
	
	// Record connection attempt metric
	if o.metricRecorder != nil {
		o.metricRecorder.RecordConnectionAttempt(
			config.ClientID,
			config.BrokerURL,
			TransportTCP, // Default, could be determined from URL
		)
	}
	
	// Store connection start time for latency calculation
	ctx = context.WithValue(ctx, "mqtt.connect.start", time.Now())
	
	fmt.Printf("MQTT Connect: %s to %s\n", config.ClientID, config.BrokerURL)
	fmt.Printf("Attributes: %+v\n", attrs)
	
	// Return completion callback
	return ctx, func(err error) {
		if err != nil {
			o.RecordError(ctx, OperationConnect, err)
		} else {
			fmt.Printf("MQTT Connect Complete\n")
		}
	}
}

// OnDisconnect is called when a connection is terminated
func (o *OTelMQTTObserver) OnDisconnect(ctx context.Context, reason observability.DisconnectReason) (context.Context, func(error)) {
	// Start span for disconnect
	spanName := o.spanBuilder.ForOperation(OperationDisconnect)
	// ctx, span := o.tracer.Start(ctx, spanName)
	// defer span.End()
	_ = spanName // Placeholder until OpenTelemetry is integrated
	
	// Build attributes
	attrs := DisconnectAttributes(reason)
	// span.SetAttributes(convertToOTelAttributes(attrs)...)
	
	fmt.Printf("MQTT Disconnect: reason=%s, code=%d\n", reason.Reason, reason.Code)
	fmt.Printf("Attributes: %+v\n", attrs)
	
	return ctx, func(err error) {
		if err != nil {
			o.RecordError(ctx, OperationDisconnect, err)
		} else {
			fmt.Printf("MQTT Disconnect Complete\n")
		}
	}
}

// OnReconnect is called when a reconnection occurs
func (o *OTelMQTTObserver) OnReconnect(ctx context.Context, attempt int, backoff time.Duration) (context.Context, func(error)) {
	// Start span for reconnect
	spanName := o.spanBuilder.ForOperation(OperationReconnect)
	// ctx, span := o.tracer.Start(ctx, spanName)
	// defer span.End()
	_ = spanName // Placeholder until OpenTelemetry is integrated
	
	// Build attributes
	attrs := NewAttributeBuilder().
		WithMessagingSystem().
		WithOperation(OperationReconnect).
		WithClientID(o.clientID).
		Build()
	attrs["mqtt.reconnect.attempt"] = attempt
	attrs["mqtt.reconnect.backoff"] = backoff.String()
	
	// Record reconnection attempt
	o.metricRecorder.RecordError(o.clientID, OperationReconnect, "reconnection_needed")
	
	fmt.Printf("MQTT Reconnect: attempt=%d, backoff=%v\n", attempt, backoff)
	fmt.Printf("Attributes: %+v\n", attrs)
	
	return ctx, func(err error) {
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
	spanName := o.spanBuilder.ForPublish(msg.Topic)
	// ctx, span := o.tracer.Start(ctx, spanName)
	// defer span.End()
	_ = spanName // Placeholder until OpenTelemetry is integrated
	
	// Build attributes using the attribute builder
	attrs := PublishAttributes(msg)
	// span.SetAttributes(convertToOTelAttributes(attrs)...)
	
	// Record publish metrics
	start := time.Now()
	ctx = context.WithValue(ctx, "mqtt.publish.start", start)
	
	fmt.Printf("MQTT Publish: topic=%s, qos=%d, size=%d bytes\n", 
		msg.Topic, msg.QoS, len(msg.Payload))
	fmt.Printf("Attributes: %+v\n", attrs)
	
	return ctx, func(result observability.PublishResult) {
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
	if startTime, ok := ctx.Value("mqtt.publish.start").(time.Time); ok {
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
	var spanName string
	if len(topics) > 0 {
		spanName = o.spanBuilder.ForSubscribe(topics[0].Topic)
	} else {
		spanName = o.spanBuilder.ForOperation(OperationSubscribe)
	}
	// ctx, span := o.tracer.Start(ctx, spanName)
	// defer span.End()
	_ = spanName // Placeholder until OpenTelemetry is integrated
	
	// Build attributes
	attrs := SubscribeAttributes(topics)
	// span.SetAttributes(convertToOTelAttributes(attrs)...)
	
	// Record subscription attempts
	for _, topic := range topics {
		o.metricRecorder.RecordSubscription(o.clientID, topic.Topic, int(topic.QoS), true)
		fmt.Printf("MQTT Subscribe: topic=%s, qos=%d\n", topic.Topic, topic.QoS)
	}
	fmt.Printf("Attributes: %+v\n", attrs)
	
	return ctx, func(result observability.SubscribeResult) {
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
	var spanName string
	if len(topics) > 0 {
		spanName = fmt.Sprintf("mqtt.unsubscribe %s", topics[0])
	} else {
		spanName = o.spanBuilder.ForOperation(OperationUnsubscribe)
	}
	// ctx, span := o.tracer.Start(ctx, spanName)
	// defer span.End()
	_ = spanName // Placeholder until OpenTelemetry is integrated
	
	// Build attributes
	attrs := NewAttributeBuilder().
		WithMessagingSystem().
		WithOperation(OperationUnsubscribe).
		WithClientID(o.clientID).
		Build()
	
	if len(topics) > 0 {
		attrs[MessagingDestinationNameKey] = topics[0]
	}
	
	for _, topic := range topics {
		fmt.Printf("MQTT Unsubscribe: topic=%s\n", topic)
	}
	fmt.Printf("Attributes: %+v\n", attrs)
	
	return ctx, func(result observability.UnsubscribeResult) {
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
	spanName := o.spanBuilder.ForReceive(msg.Topic)
	_ = spanName // Placeholder until OpenTelemetry is integrated
	
	// Build attributes
	attrs := ReceiveAttributes(msg)
	
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
	spanName := o.spanBuilder.ForOperation(OperationPing)
	_ = spanName // Placeholder until OpenTelemetry is integrated
	
	// Build attributes
	attrs := NewAttributeBuilder().
		WithMessagingSystem().
		WithOperation(OperationPing).
		WithClientID(o.clientID).
		Build()
	
	// Store ping start time for latency calculation
	ctx = context.WithValue(ctx, "mqtt.ping.start", time.Now())
	
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
	if startTime, ok := ctx.Value("mqtt.ping.start").(time.Time); ok {
		latency := time.Since(startTime)
		
		// Record ping metrics
		o.metricRecorder.RecordPing(o.clientID, "", latency) // Broker URL would be stored in observer
		
		fmt.Printf("MQTT Ping Response: latency=%v\n", latency)
	}
}

// OnAuth is called during authentication
func (o *OTelMQTTObserver) OnAuth(ctx context.Context, authData observability.AuthData) (context.Context, func(observability.AuthResult)) {
	// Start span for auth
	spanName := o.spanBuilder.ForOperation(OperationAuth)
	// ctx, span := o.tracer.Start(ctx, spanName)
	// defer span.End()
	_ = spanName // Placeholder until OpenTelemetry is integrated
	
	// Build attributes
	attrs := NewAttributeBuilder().
		WithMessagingSystem().
		WithOperation(OperationAuth).
		WithClientID(o.clientID).
		Build()
	attrs["mqtt.auth.method"] = authData.Method
	
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
// This would be implemented when OpenTelemetry dependencies are added
// func convertToOTelAttributes(attrs map[string]interface{}) []attribute.KeyValue {
//     var otelAttrs []attribute.KeyValue
//     for key, value := range attrs {
//         switch v := value.(type) {
//         case string:
//             otelAttrs = append(otelAttrs, attribute.String(key, v))
//         case int:
//             otelAttrs = append(otelAttrs, attribute.Int(key, v))
//         case bool:
//             otelAttrs = append(otelAttrs, attribute.Bool(key, v))
//         case float64:
//             otelAttrs = append(otelAttrs, attribute.Float64(key, v))
//         }
//     }
//     return otelAttrs
// }

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