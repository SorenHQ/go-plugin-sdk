package sdkv2

// PluginIntro represents the plugin introduction response
// Subject: soren.v2.<PLUGIN_ID>.@intro
type PluginIntro struct {
	Name         string       `json:"name"`
	Author       string       `json:"author"`
	Version      string       `json:"version"`
	Requirements *Requirements `json:"requirements,omitempty"`
}

// Requirements represents the plugin requirements
type Requirements struct {
	ReplyTo    string         `json:"replyTo"`
	Jsonui     map[string]any `json:"jsonui"`
	Jsonschema map[string]any `json:"jsonschema"`
	Handler    func(data []byte) any `json:"-"`
}

// PluginAction represents a single plugin action
// Subject: soren.v2.<PLUGIN_ID>.@actions
type Action struct {
	Method      string `json:"method"`
	Description string `json:"description"`
	Title       string `json:"title"`
	Icon        Icon   `json:"icon"`
	RequestHandler func(data []byte) any `json:"-"`
	Form        ActionFormBuilder `json:"-"`
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
	Handler    func(data []byte) any `json:"-"`
}

// ActionFormBuilder represents the action form configuration
// Subject: soren.v2.<PLUGIN_ID>.<ACTION>.@form
type ActionFormBuilder struct {
	Jsonui     map[string]any `json:"jsonui"`
	Jsonschema map[string]any `json:"jsonschema"`
}

// PluginEvent represents a plugin event for logging
type PluginEvent struct {
	Event     string         `json:"event" bson:"event"`
	Level     string         `json:"level" bson:"level"`
	Source    string         `json:"source" bson:"source"`
	Message   string         `json:"message" bson:"message"`
	Timestamp uint64         `json:"timestamp" bson:"timestamp"`
	Progress int            `json:"progress" bson:"progress"`
	Details   map[string]any `json:"details" bson:"details"`
}
