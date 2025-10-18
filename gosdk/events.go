package sdkv2

import (
	"encoding/json"
	"fmt"
	"time"

	nats "github.com/nats-io/nats.go"
)

// EventLogger handles logging and event emission
type EventLogger struct {
	sdk *SorenSDK
}

// NewEventLogger creates a new event logger
func NewEventLogger(sdk *SorenSDK) *EventLogger {
	return &EventLogger{sdk: sdk}
}

// Log sends a log event to the Soren platform
func (e *EventLogger) Log(level, message string, details map[string]interface{}) error {
	event := PluginEvent{
		Event:     "log",
		Level:     level,
		Source:    fmt.Sprintf("remote-plugin/%s", e.sdk.pluginID),
		Message:   message,
		Timestamp: uint64(time.Now().Unix()),
		Details:   details,
	}

	return e.sendEvent(event)
}

// LogInfo sends an info level log
func (e *EventLogger) LogInfo(message string, details map[string]interface{}) error {
	return e.Log("INFO", message, details)
}

// LogDebug sends a debug level log
func (e *EventLogger) LogDebug(message string, details map[string]interface{}) error {
	return e.Log("DEBUG", message, details)
}

// LogWarning sends a warning level log
func (e *EventLogger) LogWarning(message string, details map[string]interface{}) error {
	return e.Log("WARNING", message, details)
}

// LogError sends an error level log
func (e *EventLogger) LogError(message string, details map[string]interface{}) error {
	return e.Log("ERROR", message, details)
}

// EmitEvent sends a custom event to the Soren platform
func (e *EventLogger) EmitEvent(eventType string, data map[string]interface{}) error {
	event := PluginEvent{
		Event:     eventType,
		Level:     "INFO",
		Source:    fmt.Sprintf("remote-plugin/%s", e.sdk.pluginID),
		Message:   fmt.Sprintf("Event: %s", eventType),
		Timestamp: uint64(time.Now().Unix()),
		Details:   data,
	}

	return e.sendEvent(event)
}

// sendEvent sends an event to the Soren platform
func (e *EventLogger) sendEvent(event PluginEvent) error {
	if e.sdk.eventChannel == "" {
		return fmt.Errorf("event channel not configured")
	}

	body, err := json.Marshal([]PluginEvent{event})
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	subject := fmt.Sprintf("%s.%s.log", e.sdk.eventChannel, e.sdk.pluginID)

	msg := &nats.Msg{
		Subject: subject,
		Data:    body,
		Header:  nats.Header{"Authorization": []string{e.sdk.authKey}},
	}

	resp, err := e.sdk.conn.RequestMsg(msg, 3*time.Second)
	if err != nil {
		return fmt.Errorf("failed to send event: %w", err)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(resp.Data, &response); err != nil {
		return fmt.Errorf("failed to unmarshal response: %w", err)
	}

	// Check if the response indicates success
	if result, ok := response["result"].(string); ok && result != "OK" {
		return fmt.Errorf("event sending failed: %s", result)
	}

	return nil
}

// SendMultipleEvents sends multiple events in a single request
func (e *EventLogger) SendMultipleEvents(events ...PluginEvent) error {
	if e.sdk.eventChannel == "" {
		return fmt.Errorf("event channel not configured")
	}

	body, err := json.Marshal(events)
	if err != nil {
		return fmt.Errorf("failed to marshal events: %w", err)
	}

	subject := fmt.Sprintf("%s.%s.log", e.sdk.eventChannel, e.sdk.pluginID)

	msg := &nats.Msg{
		Subject: subject,
		Data:    body,
		Header:  nats.Header{"Authorization": []string{e.sdk.authKey}},
	}

	resp, err := e.sdk.conn.RequestMsg(msg, 3*time.Second)
	if err != nil {
		return fmt.Errorf("failed to send events: %w", err)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(resp.Data, &response); err != nil {
		return fmt.Errorf("failed to unmarshal response: %w", err)
	}

	// Check if the response indicates success
	if result, ok := response["result"].(string); ok && result != "OK" {
		return fmt.Errorf("events sending failed: %s", result)
	}

	return nil
}

