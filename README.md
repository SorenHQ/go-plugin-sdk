# Sorenv2 Protocol Plugin SDK

This repository contains the Go SDK for creating plugins that implement the Sorenv2 protocol. The SDK provides a simple and straightforward way to create plugins that can be integrated with the Soren platform.

## Quick Start

```go
import (
    sdkv2 "github.com/sorenhq/go-plugin-sdk/gosdk"
    models "github.com/sorenhq/go-plugin-sdk/gosdk/models"
)
```

## Plugin Structure

A Sorenv2 plugin consists of several key components:

### 1. Plugin Instance

Create a new plugin instance using the SDK:

```go
sdkInstance, err := sdkv2.NewFromEnv()
plugin := sdkv2.NewPlugin(sdkInstance)
```

### 2. Plugin Introduction

Set up your plugin's basic information using `SetIntro`:

```go
plugin.SetIntro(models.PluginIntro{
    Name:    "Your Plugin Name",
    Version: "1.0.0",
    Author:  "Your Name",
    Requirements: &models.Requirements{
        ReplyTo:    "init.config",
        Jsonui:     map[string]any{...},  // UI configuration
        Jsonschema: map[string]any{...},  // JSON schema for validation
    },
})
```

### 3. Plugin Settings

Configure plugin settings using `SetSettings`:

```go
plugin.SetSettings(&models.Settings{
    ReplyTo: "settings.config.submit",
    Jsonui: map[string]any{...},      // UI configuration
    Jsonschema: map[string]any{...},  // JSON schema for settings
}, settingsUpdateHandler)
```

### 4. Actions

Define plugin actions using `AddActions`:

```go
plugin.AddActions([]models.Action{
    {
        Method: "your.action.method",
        Title:  "Action Title",
        Form: models.ActionFormBuilder{
            Jsonui:     map[string]any{...},  // UI configuration
            Jsonschema: map[string]any{...},  // JSON schema for action input
        },
        RequestHandler: func(data []byte) any {
            // Handle the action here
            return map[string]any{"result": "success"}
        },
    },
})
```

### 5. Event Logging

Log events from your plugin:

```go
event := sdkv2.NewEventLogger(sdkInstance)
event.Log("source-identifier", models.LogLevelInfo, "message", nil)
```

## Components Reference

### PluginIntro

- `Name`: Name of your plugin
- `Version`: Plugin version
- `Author`: Plugin author
- `Requirements`: Initial configuration requirements
  - `ReplyTo`: Configuration endpoint
  - `Jsonui`: UI configuration for requirements
  - `Jsonschema`: JSON schema for validating requirements

### Settings

- `ReplyTo`: Settings endpoint
- `Jsonui`: UI configuration for settings
- `Jsonschema`: JSON schema for validating settings
- Settings handler function to process updates

### Action

- `Method`: Unique identifier for the action
- `Title`: Display name of the action
- `Form`: Input form configuration
  - `Jsonui`: UI configuration for the action form
  - `Jsonschema`: JSON schema for validating action input
- `RequestHandler`: Function to handle the action execution

### Event Logger

- Log levels: Info, Warning, Error
- Parameters:
  - Source identifier
  - Log level
  - Message
  - Additional data (optional)

## Starting the Plugin

To start your plugin:

```go
plugin.Start()
```

## Best Practices

1. Always implement proper error handling
2. Use environment variables for configuration
3. Implement graceful shutdown
4. Validate all inputs using JSON schema
5. Log important events and errors

## Example

See the test file for a complete working example of a plugin implementation.
