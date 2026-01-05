package sdkv2

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	nats "github.com/nats-io/nats.go"
	"github.com/sorenhq/go-plugin-sdk/gosdk/models"
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
func (e *EventLogger) Log(source string, level models.LogLevel, message string, details map[string]any) error {
	event := models.PluginEvent{
		Event:     models.EventTypeLog,
		Level:     level,
		Source:    fmt.Sprintf("%s - %s", e.sdk.pluginID, source),
		Message:   message,
		Timestamp: uint64(time.Now().Unix()),
		Details:   details,
	}

	return e.sendEvent(event)
}

// EmitEvent sends a custom event to the Soren platform
func (e *EventLogger) EmitEvent(eventType models.EventType, data map[string]any) error {
	event := models.PluginEvent{
		Event:     eventType,
		Level:     models.LogLevelInfo,
		Source:    e.sdk.pluginID,
		Message:   fmt.Sprintf("Event: %s", eventType),
		Timestamp: uint64(time.Now().Unix()),
		Details:   data,
	}

	return e.sendEvent(event)
}

// sendEvent sends an event to the Soren platform
func (e *EventLogger) sendEvent(event models.PluginEvent) error {
	if e.sdk.eventChannel == "" {
		return fmt.Errorf("event channel not configured")
	}

	body, err := json.Marshal([]models.PluginEvent{event})
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}
	subject := fmt.Sprintf("%s.%s.log", e.sdk.eventChannel, e.sdk.pluginID)
	if len(strings.Split(e.sdk.pluginID, ".")) > 0 {
		subject = fmt.Sprintf("%s.log", e.sdk.eventChannel)

	}

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
func (e *EventLogger) SendMultipleEvents(events ...models.PluginEvent) error {
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
