# OpenTelemetry Integration Examples

This directory contains examples demonstrating how to integrate the paho.golang MQTT clients with OpenTelemetry observability.

## Examples

### autopaho Example

Location: `autopaho/main.go`

Demonstrates how to use the `autopaho` client with OpenTelemetry observability:

- Automatic connection management with observability tracing
- Publish and subscribe operations with automatic tracing
- Connection lifecycle events tracked through observability
- Graceful shutdown with observability

### paho Example

Location: `paho/main.go`

Demonstrates how to use the `paho` client with OpenTelemetry observability:

- Manual connection management with observability tracing
- Publish and subscribe operations with tracing
- Connection lifecycle events tracked through observability
- Graceful shutdown with observability

## Prerequisites

1. **MQTT Broker**: You need a running MQTT broker. The examples default to `mqtt://localhost:1883`
2. **Go Modules**: Ensure your Go environment supports modules
3. **OpenTelemetry**: The examples use the OpenTelemetry observer from the `otel` package

## Running the Examples

### Start an MQTT Broker

You can use Mosquitto or any other MQTT broker:

```bash
# Using Docker
docker run -it -p 1883:1883 eclipse-mosquitto:2.0

# Or install locally (macOS)
brew install mosquitto
mosquitto -p 1883
```

### Run autopaho Example

```bash
cd examples/otel/autopaho
go run main.go
```

### Run paho Example

```bash
cd examples/otel/paho
go run main.go
```

## Configuration

### Environment Variables

- `MQTT_BROKER_URL`: Override the default broker URL (default: `mqtt://localhost:1883`)

### Example Usage

```bash
# Use a different broker
MQTT_BROKER_URL="mqtt://test.mosquitto.org:1883" go run main.go

# Use TLS
MQTT_BROKER_URL="mqtts://broker.example.com:8883" go run main.go
```

## What You'll See

Both examples will:

1. Connect to the MQTT broker with observability tracing
2. Subscribe to `test/topic`
3. Publish messages every 5 seconds to `test/topic`
4. Receive and display the published messages
5. Handle graceful shutdown when interrupted (Ctrl+C)

The observability integration will automatically:

- Trace connection attempts and results
- Trace publish operations with latency and success/failure
- Trace subscribe operations
- Trace message reception
- Track connection lifecycle events

## Integration Details

### autopaho Integration

The `autopaho` client integrates observability through:

- `Observer` field in `ClientConfig`
- Automatic tracing of connection management
- Built-in retry and reconnection observability

### paho Integration

The `paho` client integrates observability through:

- `Observer` field in `ClientConfig`
- Manual connection management with tracing
- Explicit observability calls for each operation

## Observability Data

The OpenTelemetry observer tracks:

- **Connection Events**: Connect, disconnect, reconnect attempts
- **Publish Operations**: Message publishing with QoS, topic, latency
- **Subscribe Operations**: Topic subscriptions with QoS
- **Message Reception**: Incoming messages with metadata
- **Error Conditions**: Connection failures, publish failures, etc.

## Extending the Examples

You can extend these examples by:

1. **Adding Custom Tracing**: Implement additional spans around your business logic
2. **Metrics Collection**: Add custom metrics using OpenTelemetry metrics
3. **Error Handling**: Enhance error handling with observability context
4. **Configuration**: Add more sophisticated configuration options
5. **Multiple Topics**: Subscribe to and publish on multiple topics

## Troubleshooting

### Connection Issues

- Ensure your MQTT broker is running and accessible
- Check firewall settings for port 1883 (or your custom port)
- Verify the broker URL format

### Import Issues

- Ensure you're running from the correct directory
- Check that all dependencies are available: `go mod tidy`

### Observability Issues

- Verify OpenTelemetry configuration
- Check that the `otel` package is properly imported
- Ensure observability backend is configured if using external systems