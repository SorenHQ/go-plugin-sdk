package models

import "github.com/nats-io/nats.go"

// PluginIntro represents the plugin introduction response
// Subject: soren.v2.<PLUGIN_ID>.@intro
type PluginIntro struct {
	Name         string        `json:"name"`
	Author       string        `json:"author"`
	Version      string        `json:"version"`
	Requirements *Requirements `json:"requirements,omitempty"`
}

// Requirements represents the plugin requirements
type Requirements struct {
	ReplyTo    string                  `json:"replyTo"`
	Jsonui     map[string]any          `json:"jsonui"`
	Jsonschema map[string]any          `json:"jsonschema"`
	Handler    func(msg *nats.Msg) any `json:"-"`
}

// PluginAction represents a single plugin action
// Subject: soren.v2.<PLUGIN_ID>.@actions
type Action struct {
	Method         string                  `json:"method"`
	Description    string                  `json:"description"`
	Title          string                  `json:"title"`
	Icon           Icon                    `json:"icon"`
	RequestHandler func(msg *nats.Msg)  `json:"-"`
	Form           ActionFormBuilder       `json:"-"`
}

// Icon represents an icon for an action
type Icon struct {
	Ref  string `json:"ref"`
	Icon string `json:"icon"`
}

// Settings represents the settings form configuration
// Subject: soren.v2.<PLUGIN_ID>.@settings
type Settings struct {
	ReplyTo    string         `json:"replyTo"`
	Jsonui     map[string]any `json:"jsonui"`
	Jsonschema map[string]any `json:"jsonschema"`
	Data       map[string]any `json:"data"` // Current settings data
	Handler func(msg *nats.Msg) any `json:"-"`
}

// ActionFormBuilder represents the action form configuration
// Subject: soren.v2.<PLUGIN_ID>.<ACTION>.@form
type ActionFormBuilder struct {
	Jsonui     map[string]any `json:"jsonui"`
	Jsonschema map[string]any `json:"jsonschema"`
}

// PluginEvent represents a plugin event for logging
type PluginEvent struct {
	Event     EventType      `json:"event" bson:"event"`
	Level     LogLevel       `json:"level" bson:"level"`
	Source    string         `json:"source" bson:"source"`
	Message   string         `json:"message" bson:"message"`
	Timestamp uint64         `json:"timestamp" bson:"timestamp"`
	Details   map[string]any `json:"details" bson:"details"`
}

type JobProgress struct {
	Progress int            `json:"progress" bson:"progress"`
	Frame    Frame          `json:"frame" bson:"frame"`
	Details  map[string]any `json:"details"`
}

type Frame struct {
	Title   string `json:"title" bson:"title"`
	Content string `json:"content" bson:"content"`
}

type JobBodyContent struct {
	JobId    string         `json:"jobId"`
	Progress int            `json:"progress"`
	Details  map[string]any `json:"details"`
	CommitOn string         `json:"commit_on"`
}

type ActionRequestContent struct {
	Registry map[string]any `json:"_registry"`
	Body     map[string]any `json:"body"`
}
