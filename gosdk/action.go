package sdkv2

import (
	"context"
	"log"
	"sync"

	"github.com/bytedance/sonic"
	"github.com/nats-io/nats.go"
	"github.com/sorenhq/go-plugin-sdk/gosdk/models"
	"github.com/sorenhq/go-plugin-sdk/logtool"
)

type Plugin struct {
	sdk         *SorenSDK
	Intro       models.PluginIntro
	Settings    *models.Settings
	Actions     []models.Action
	jobEntityId map[string]string // Maps jobId to entityId (spaceId)
	jobMutex    sync.RWMutex
}

func NewPlugin(sdk *SorenSDK) *Plugin {
	logtool.Init("SOREN-SDK", true)
	newPlugin := &Plugin{
		sdk:         sdk,
		jobEntityId: make(map[string]string),
	}
	GetPluginHolder().add(sdk.pluginID, newPlugin)
	return newPlugin
}
func (p *Plugin) GetContext() context.Context {
	return p.sdk.ctx
}
func (p *Plugin) SetSettings(settings *models.Settings, handler func(msg *nats.Msg) any) {
	p.Settings = settings
	if p.Settings != nil {
		p.Settings.Handler = handler
	}
}
func (p *Plugin) SetActions(actions []models.Action) {
	p.Actions = actions
}
func (p *Plugin) SetIntro(intro models.PluginIntro, handler func(msg *nats.Msg) any) {
	p.Intro = intro
	if p.Intro.Requirements != nil {
		p.Intro.Requirements.Handler = handler
	}
}
func (p *Plugin) AddActions(actions []models.Action) {
	p.Actions = append(p.Actions, actions...)
}
func (p *Plugin) Start() error {
	err := p.IntroHandler()
	if err != nil {
		return err
	}
	err = p.SettingsHandler()
	if err != nil {
		return err
	}
	p.ActionsHandler()
	// Log plugin startup event (only if event channel and auth key are configured)
	if p.sdk.eventChannel != "" && p.sdk.authKey != "" {
		event := NewEventLogger(p.sdk)
		actionsByte, _ := sonic.Marshal(p.Actions)
		if len(actionsByte) > 0 {
			actionsList := []map[string]any{}
			if err := sonic.Unmarshal(actionsByte, &actionsList); err == nil {
				// Log event, but don't fail if it errors (non-critical)
				if err := event.Log("soren-sdk-init", models.LogLevelInfo, "start plugin", map[string]any{"actions": actionsList}); err != nil {
					log.Printf("Failed to log startup event (non-critical): %v", err)
				}
			}
		}
	} else {
		if p.sdk.eventChannel == "" {
			log.Printf("Event channel not configured, skipping startup event log")
		}
		if p.sdk.authKey == "" {
			log.Printf("Auth key not configured, skipping startup event log")
		}
	}

	<-p.sdk.ctx.Done()
	log.Println("Plugin context done, exiting plugin:", p.Intro.Name)
	return nil
}

// StoreEntityIdForJob stores the entityId (spaceId) for a given jobId
func (p *Plugin) StoreEntityIdForJob(jobId, entityId string) {
	p.jobMutex.Lock()
	defer p.jobMutex.Unlock()
	p.jobEntityId[jobId] = entityId
}

// getEntityIdForJob retrieves the entityId for a given jobId
func (p *Plugin) getEntityIdForJob(jobId string) string {
	p.jobMutex.RLock()
	defer p.jobMutex.RUnlock()
	return p.jobEntityId[jobId]
}

func (p *Plugin) Done(jobId string, data map[string]any) any {
	// Try to get entityId for this job - if present, use gateway subject pattern
	entityId := p.getEntityIdForJob(jobId)
	var sub string
	if entityId != "" {
		// Use gateway pattern: soren.v2.bin.{entityId}.{pluginId}.{jobId}.progress
		sub = p.sdk.makeGatewayJobSubject(entityId, jobId, string(models.ProgressCommand))
	} else {
		// Fallback to CPU pattern
		sub = p.sdk.makeJobSubject(jobId, string(models.ProgressCommand))
	}

	// Gateway expects JobBodyContent structure with jobId included
	jobBody := models.JobBodyContent{
		JobId:    jobId,
		Progress: 100,
		Details:  data,
		CommitOn: "",
	}
	dataByte, err := sonic.Marshal(jobBody)
	if err != nil {
		log.Printf("Done command error marshaling: %v (subject: %s, jobId: %s)", err, sub, jobId)
		return err
	}

	log.Printf("Publishing Done result - subject: %s, jobId: %s, entityId: %s, data size: %d bytes", sub, jobId, entityId, len(dataByte))
	if len(dataByte) < 2000 {
		log.Printf("Done result data (first 2000 chars): %s", string(dataByte))
	} else {
		log.Printf("Done result data preview (first 500 chars): %s...", string(dataByte[:500]))
	}

	err = p.sdk.conn.Publish(sub, dataByte)
	if err != nil {
		log.Printf("Failed to publish done result: %v (subject: %s, jobId: %s)", err, sub, jobId)
		return err
	}

	return nil
}
func (p *Plugin) Progress(jobId string, command models.Command, data models.JobProgress) any {
	// Try to get entityId for this job - if present, use gateway subject pattern
	entityId := p.getEntityIdForJob(jobId)
	var sub string
	if entityId != "" {
		// Use gateway pattern: soren.v2.bin.{entityId}.{pluginId}.{jobId}.{command}
		sub = p.sdk.makeGatewayJobSubject(entityId, jobId, string(command))
	} else {
		// Fallback to CPU pattern
		sub = p.sdk.makeJobSubject(jobId, string(command))
	}

	dataByte, err := sonic.Marshal(data)
	if err != nil {
		log.Println("progress command ", command, " error:", err)
		return err
	}

	// Use Publish (fire-and-forget) for progress updates
	// Gateway subscribes to progress updates, so we don't need Request/Response
	err = p.sdk.conn.Publish(sub, dataByte)
	if err != nil {
		log.Printf("Failed to publish progress update: %v (subject: %s, jobId: %s)", err, sub, jobId)
		return err
	}

	return nil
}
