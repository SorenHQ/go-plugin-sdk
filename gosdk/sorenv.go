package sdkv2

import (
	"context"
	"fmt"
	"os"

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

	// Validate required configuration
	if config.AgentURI == "" {
		return nil, fmt.Errorf("agent URI is required")
	}
	if config.PluginID == "" {
		return nil, fmt.Errorf("plugin ID is required")
	}

	// Connect to NATS
	nc, err := nats.Connect(config.AgentURI)
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
func (s *SorenSDK) makeSubject(action string) string {
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
// makeActionSubject creates a subject for action execution
func (s *SorenSDK) makeActionCpu(action string) string {
	return fmt.Sprintf("soren.cpu.%s.%s", s.pluginID, action)
}

// makeJobSubject creates a subject for job updates
func (s *SorenSDK) makeJobSubject(jobID, jobUpdate string) string {
	return fmt.Sprintf("soren.cpu.%s.%s.%s", s.pluginID, jobID, jobUpdate)
}

// makeFormSubject creates a subject for form requests
func (s *SorenSDK) makeFormSubject(action string) string {
	return fmt.Sprintf("soren.v2.%s.%s.@form", s.pluginID, action)
}

// makeProgressSubject creates a subject for progress updates
func (s *SorenSDK) makeProgressSubject(jobID string) string {
	return fmt.Sprintf("soren.cpu.%s.%s.*", s.pluginID, jobID)
}

// // SendProgress sends a progress update for a specific job
// func (s *SorenSDK) SendProgress(jobID string, progress int, message string, data map[string]any) error {
// 	progressMsg := ProgressMessage{
// 		JobID:     jobID,
// 		Progress:  progress,
// 		Message:   message,
// 		Data:      data,
// 		Timestamp: time.Now().Unix(),
// 	}

// 	msgData, err := sonic.Marshal(progressMsg)
// 	if err != nil {
// 		return fmt.Errorf("failed to marshal progress message: %w", err)
// 	}

// 	// Send to the progress subject
// 	subject := s.makeProgressSubject(jobID)
// 	return s.conn.Publish(subject, msgData)
// }

// // SendJobStatus sends a job status update
// func (s *SorenSDK) SendJobStatus(jobID, status string, progress int, message string, result map[string]any, errMsg string) error {
// 	jobStatus := JobStatus{
// 		JobID:     jobID,
// 		Status:    status,
// 		Progress:  progress,
// 		Message:   message,
// 		Result:    result,
// 		Error:     errMsg,
// 		Timestamp: time.Now().Unix(),
// 	}

// 	msgData, err := sonic.Marshal(jobStatus)
// 	if err != nil {
// 		return fmt.Errorf("failed to marshal job status: %w", err)
// 	}

// 	// Send to the progress subject
// 	subject := s.makeProgressSubject(jobID)
// 	return s.conn.Publish(subject, msgData)
// }

// // CompleteJob marks a job as completed with 100% progress
// func (s *SorenSDK) CompleteJob(jobID string, result map[string]any, message string) error {
// 	return s.SendJobStatus(jobID, "completed", 100, message, result, "")
// }

// // FailJob marks a job as failed
// func (s *SorenSDK) FailJob(jobID string, errMsg string) error {
// 	return s.SendJobStatus(jobID, "failed", 0, "", nil, errMsg)
// }
