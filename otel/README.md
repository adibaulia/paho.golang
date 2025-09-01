# MQTT Observability with OpenTelemetry

This package provides OpenTelemetry implementation for MQTT observability. It works with the core `observability` package to provide distributed tracing, metrics collection, and error tracking for MQTT applications.

> **Note**: The core MQTT observability interface is now in the `observability` package, making it independent of any specific observability implementation. This package provides the OpenTelemetry-specific implementation.

## Features

- **OpenTelemetry Implementation**: Full OpenTelemetry implementation of the core observability interface
- **Distributed Tracing**: Track MQTT operations across services with OpenTelemetry spans
- **Comprehensive Metrics**: Monitor connection health, message throughput, and error rates
- **Error Categorization**: Automatic classification of MQTT errors for better debugging
- **Performance Monitoring**: Track latency, queue sizes, and throughput
- **Seamless Integration**: Works with existing OpenTelemetry infrastructure
- **Pluggable Design**: Can be easily swapped with other observability implementations

## Quick Start

### 1. Basic Setup

```go
package main

import(
    "context"
    "log"
    "time"
    
    "github.com/eclipse/paho.golang/observability"
    "github.com/eclipse/paho.golang/otel"
)

func main() {
    // Create an OpenTelemetry observer instance
    var observer observability.MQTTObserver = otel.NewOTelMQTTObserver("my-mqtt-client")
    
    // Use with your MQTT client
    client := &SimpleMQTTClient{
        observer: observer,
    }
    
    // Connect and start observing
    config := observability.ConnectionConfig{
        ClientID:   "my-client",
        BrokerURL:  "tcp://localhost:1883",
        Transport:  "tcp",
        Username:   "user",
        KeepAlive:  60,
        CleanStart: true,
    }
    
    client.Connect(config)
}
```

### 2. Integration with Existing MQTT Client

```go
import (
    "github.com/eclipse/paho.golang/observability"
    "github.com/eclipse/paho.golang/otel"
)

type MQTTClientWithObservability struct {
    client   *paho.Client
    observer observability.MQTTObserver
}

func (c *MQTTClientWithObservability) Publish(topic string, payload []byte, qos byte) error {
    // Start observing the publish operation
    message := observability.PublishMessage{
        Topic:   topic,
        Payload: payload,
        QoS:     qos,
    }
    
    ctx, onComplete := c.observer.OnPublish(context.Background(), message)
    
    // Perform the actual publish
    err := c.client.Publish(ctx, &paho.Publish{
        Topic:   topic,
        Payload: payload,
        QoS:     qos,
    })
    
    // Complete the observation
    onComplete(observability.PublishResult{Error: err})
    
    return err
}
```

## Configuration

### OpenTelemetry Setup

```go
package main

import (
    "context"
    "log"
    
    "go.opentelemetry.io/otel"
    "go.opentelemetry.io/otel/exporters/jaeger"
    "go.opentelemetry.io/otel/exporters/prometheus"
    "go.opentelemetry.io/otel/sdk/metric"
    "go.opentelemetry.io/otel/sdk/trace"
)

func setupOpenTelemetry() {
    // Setup tracing
    jaegerExporter, err := jaeger.New(jaeger.WithCollectorEndpoint(
        jaeger.WithEndpoint("http://localhost:14268/api/traces"),
    ))
    if err != nil {
        log.Fatal(err)
    }
    
    tp := trace.NewTracerProvider(
        trace.WithBatcher(jaegerExporter),
        trace.WithResource(resource.NewWithAttributes(
            semconv.SchemaURL,
            semconv.ServiceNameKey.String("mqtt-client"),
        )),
    )
    otel.SetTracerProvider(tp)
    
    // Setup metrics
    prometheusExporter, err := prometheus.New()
    if err != nil {
        log.Fatal(err)
    }
    
    mp := metric.NewMeterProvider(
        metric.WithReader(prometheusExporter),
    )
    otel.SetMeterProvider(mp)
}
```

## Metrics Reference

### Connection Metrics

- `mqtt.connection.attempts` (Counter): Total connection attempts
- `mqtt.connection.success` (Counter): Successful connections
- `mqtt.connection.failures` (Counter): Failed connections by error category
- `mqtt.connection.latency` (Histogram): Connection establishment latency
- `mqtt.connections.active` (Gauge): Currently active connections

### Message Metrics

- `mqtt.messages.published` (Counter): Messages published by QoS level
- `mqtt.messages.publish_failed` (Counter): Failed publishes by error category
- `mqtt.messages.received` (Counter): Messages received by QoS level
- `mqtt.publish.latency` (Histogram): Publish operation latency
- `mqtt.publish.size` (Histogram): Published message size
- `mqtt.message.size` (Histogram): Received message size

### Subscription Metrics

- `mqtt.subscribe.attempts` (Counter): Subscription attempts
- `mqtt.subscribe.success` (Counter): Successful subscriptions
- `mqtt.subscribe.failures` (Counter): Failed subscriptions
- `mqtt.subscriptions.active` (Gauge): Active subscriptions
- `mqtt.unsubscribe.attempts` (Counter): Unsubscription attempts
- `mqtt.unsubscribe.success` (Counter): Successful unsubscriptions

### Error Metrics

- `mqtt.errors` (Counter): General errors by operation and category
- `mqtt.ping.requests` (Counter): Ping requests sent
- `mqtt.ping.latency` (Histogram): Ping response latency

### Queue Metrics

- `mqtt.publish.queue.size` (Gauge): Current publish queue size
- `mqtt.throughput.bytes` (Counter): Total bytes transferred
- `mqtt.throughput.messages` (Counter): Total messages transferred

## Error Categories

The observability system automatically categorizes errors into the following types:

- **connection**: Connection-related errors
- **protocol**: MQTT protocol violations
- **network**: Network connectivity issues
- **timeout**: Operation timeouts
- **authentication**: Authentication failures
- **permission**: Authorization/permission errors
- **qos**: Quality of Service related errors
- **payload**: Message payload issues
- **unknown**: Uncategorized errors

## Advanced Usage

### Custom Metric Recording

```go
// Get the metric recorder for custom metrics
recorder := observer.GetMetricRecorder()

// Record custom queue size
recorder.RecordQueueSize("my-client", 42)

// Record throughput
recorder.RecordThroughput("my-client", 1024, 10)

// Get metrics summary
summary := recorder.GetMetricsSummary()
fmt.Printf("Metrics: %+v\n", summary)
```

### Error Handling with Categorization

```go
// The observer automatically categorizes errors
err := errors.New("connection refused")
category := recorder.CategorizeError(err)
// category will be "network"

// Record the error with automatic categorization
recorder.RecordError("my-client", "connect", string(category))
```

### Integration with autopaho

```go
package main

import (
    "github.com/eclipse/paho.golang/autopaho"
    "github.com/eclipse/paho.golang/observability"
    "github.com/eclipse/paho.golang/otel"
)

func main() {
    var observer observability.MQTTObserver = otel.NewOTelMQTTObserver("autopaho-client")
    
    config := autopaho.ClientConfig{
        BrokerUrls: []*url.URL{brokerURL},
        KeepAlive:  30,
        OnConnectionUp: func(cm *autopaho.ConnectionManager, connAck *paho.Connack) {
            ctx, onComplete := observer.OnConnect(context.Background(), observability.ConnectionConfig{
                ClientID:  "autopaho-client",
                BrokerURL: brokerURL.String(),
                Transport: "tcp",
            })
            onComplete(observability.ConnectionResult{Connected: true})
        },
        OnConnectError: func(err error) {
            ctx, onComplete := observer.OnDisconnect(context.Background(), "autopaho-client")
            onComplete(observability.DisconnectResult{Error: err})
        },
    }
    
    cm, err := autopaho.NewConnection(context.Background(), config)
    if err != nil {
        log.Fatal(err)
    }
    
    // Wrap publish operations
    originalPublish := cm.Publish
    cm.Publish = func(ctx context.Context, pub *paho.Publish) (*paho.PublishResponse, error) {
        message := observability.PublishMessage{
            Topic:   pub.Topic,
            Payload: pub.Payload,
            QoS:     pub.QoS,
        }
        
        ctx, onComplete := observer.OnPublish(ctx, message)
        
        resp, err := originalPublish(ctx, pub)
        
        onComplete(observability.PublishResult{Error: err})
        
        return resp, err
    }
}
```

## Monitoring and Alerting

### Prometheus Queries

```promql
# Connection success rate
rate(mqtt_connection_success_total[5m]) / rate(mqtt_connection_attempts_total[5m])

# Average publish latency
rate(mqtt_publish_latency_sum[5m]) / rate(mqtt_publish_latency_count[5m])

# Error rate by category
sum(rate(mqtt_errors_total[5m])) by (error_type)

# Active connections
mqtt_connections_active

# Queue size trending
mqtt_publish_queue_size
```

### Grafana Dashboard

Create dashboards with the following panels:

1. **Connection Health**: Success rate, active connections, connection latency
2. **Message Flow**: Publish/receive rates, message sizes, queue sizes
3. **Error Analysis**: Error rates by category, failed operations
4. **Performance**: Latency percentiles, throughput metrics

## Best Practices

1. **Resource Management**: Always properly close observers and connections
2. **Sampling**: Use appropriate sampling rates for high-volume applications
3. **Error Handling**: Don't let observability failures affect your main application logic
4. **Metric Cardinality**: Be mindful of high-cardinality labels (like topic names)
5. **Performance**: Observability should have minimal impact on application performance

## Troubleshooting

### Common Issues

1. **High Memory Usage**: Check for metric cardinality explosion
2. **Missing Metrics**: Verify OpenTelemetry setup and exporters
3. **Performance Impact**: Review sampling configuration
4. **Error Categorization**: Check error message patterns for proper categorization

### Debug Mode

```go
// Enable debug logging
observer := otel.NewOTelMQTTObserver("debug-client")
observer.SetDebugMode(true)
```

## Contributing

Contributions are welcome! Please ensure:

1. All new metrics follow OpenTelemetry semantic conventions
2. Error categories are well-defined and documented
3. Examples are provided for new features
4. Tests cover new functionality

## License

This project is licensed under the same license as the Eclipse Paho MQTT Go client.