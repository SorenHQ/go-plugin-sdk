package sdkv2

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"strings"

	nats "github.com/nats-io/nats.go"
)

// SorenSDK represents the main SDK instance for Soren v2 protocol
type SorenSDK struct {
	conn         *nats.Conn
	pluginID     string
	agentURI     string
	authKey      string
	eventChannel string
	storeChannel string
	ctx          context.Context
	cancel       context.CancelFunc
}

// Config holds the configuration for the Soren SDK
type Config struct {
	AgentURI     string
	AgentCred    string
	PluginID     string
	AuthKey      string
	EventChannel string
	StoreChannel string
}

// New creates a new Soren SDK instance
func New(config *Config) (*SorenSDK, error) {
	if config == nil {
		config = &Config{}
	}

	// Load from environment if not provided
	if config.AgentURI == "" {
		config.AgentURI = os.Getenv("AGENT_URI")
	}
	if config.PluginID == "" {
		config.PluginID = os.Getenv("PLUGIN_ID")
	}
	if config.AuthKey == "" {
		config.AuthKey = os.Getenv("SOREN_AUTH_KEY")
	}
	if config.EventChannel == "" {
		config.EventChannel = os.Getenv("SOREN_EVENT_CHANNEL")
	}
	if config.StoreChannel == "" {
		config.StoreChannel = os.Getenv("SOREN_STORE")
	}
	if config.AgentCred == "" {
		config.AgentCred = os.Getenv("AGENT_CRED")
	}
	// Validate required configuration
	if config.AgentURI == "" {
		return nil, fmt.Errorf("agent URI is required")
	}
	if config.PluginID == "" {
		return nil, fmt.Errorf("plugin ID is required")
	}
	var nc *nats.Conn
	var err error
	// Connect to NATS
	if config.AgentCred != "" {
		if strings.HasPrefix(config.AgentCred, "-----BEGIN") {
			nc, err = nats.Connect(config.AgentURI, nats.UserCredentialBytes([]byte(config.AgentCred)))
		}
		credByte, err := base64.StdEncoding.DecodeString(config.AgentCred)
		if err != nil {
			return nil, err
		}
		nc, err = nats.Connect(config.AgentURI, nats.UserCredentialBytes([]byte(credByte)))

	} else {
		nc, err = nats.Connect(config.AgentURI)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to connect to NATS: %w", err)
	}
	ctx, cancel := context.WithCancel(context.Background())

	sdk := &SorenSDK{
		conn:         nc,
		pluginID:     config.PluginID,
		agentURI:     config.AgentURI,
		authKey:      config.AuthKey,
		eventChannel: config.EventChannel,
		storeChannel: config.StoreChannel,
		ctx:          ctx,
		cancel:       cancel,
	}

	return sdk, nil
}

// NewFromEnv creates a new Soren SDK instance using environment variables
func NewFromEnv() (*SorenSDK, error) {
	return New(nil)
}

// Close closes the SDK connection and cleans up resources
func (s *SorenSDK) Close() error {
	s.cancel()
	s.conn.Close()
	return nil
}

// GetConnection returns the underlying NATS connection
func (s *SorenSDK) GetConnection() *nats.Conn {
	return s.conn
}

// GetPluginID returns the plugin ID
func (s *SorenSDK) GetPluginID() string {
	return s.pluginID
}

// GetContext returns the SDK context
func (s *SorenSDK) GetContext() context.Context {
	return s.ctx
}

// makeSubject creates a subject with the soren.v2 prefix
// For internal plugins (pluginID contains "bin.*"), uses gateway pattern: soren.v2.bin.*.{uuid}.{action}
// For other plugins, uses standard pattern: soren.v2.{pluginID}.{action}
func (s *SorenSDK) makeSubject(action string) string {
	// Check if this is an internal plugin (pluginID contains "bin.*")
	if strings.HasPrefix(s.pluginID, "bin.*.") {
		// Extract UUID part after "bin.*."
		parts := strings.Split(s.pluginID, ".")
		if len(parts) >= 3 {
			// Get the UUID part (last part after bin.*)
			uuid := parts[len(parts)-1]
			// Use gateway pattern: soren.v2.bin.*.{uuid}.{action}
			// The * wildcard will match any entityId sent by the gateway
			return fmt.Sprintf("soren.v2.bin.*.%s.%s", uuid, action)
		}
	}
	// Standard pattern for non-internal plugins
	return fmt.Sprintf("soren.v2.%s.%s", s.pluginID, action)
}

// makeSettingsSubject creates a subject with the soren.v2 prefix
func (s *SorenSDK) makeSettingsSubject() string {
	return fmt.Sprintf("soren.v2.%s.@settings", s.pluginID)
}

// makeActionsListSubject creates a subject with the soren.v2 prefix
func (s *SorenSDK) makeActionsListSubject() string {
	return fmt.Sprintf("soren.v2.%s.@actions", s.pluginID)
}

// makeIntroSubject creates a subject with the soren.v2 prefix
func (s *SorenSDK) makeIntroSubject() string {
	return fmt.Sprintf("soren.v2.%s.@intro", s.pluginID)
}

// makeActionCpu creates a subject for CPU/job processing (original purpose)
func (s *SorenSDK) makeActionCpu(action string) string {
	return fmt.Sprintf("soren.cpu.%s.%s", s.pluginID, action)
}

// makeJobSubject creates a subject for job updates (CPU pattern)
func (s *SorenSDK) makeJobSubject(jobID, jobUpdate string) string {
	return fmt.Sprintf("soren.cpu.%s.%s.%s", s.pluginID, jobID, jobUpdate)
}


// makeFormSubject creates a subject for form requests
func (s *SorenSDK) makeFormSubject(action string) string {
	return fmt.Sprintf("soren.v2.%s.%s.@form", s.pluginID, action)
}
