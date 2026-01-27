package sdkv2

import (
	"context"
	"log"
	"strings"
	"time"

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
<<<<<<< HEAD
		sdk:         sdk,
		jobEntityId: make(map[string]string),
=======
		sdk: sdk,
>>>>>>> v0.2.3
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
	event := NewEventLogger(p.sdk)
	actionsByte, _ := sonic.Marshal(p.Actions)
	if len(actionsByte) > 0 {
		actionsList := []map[string]any{}
		if err := sonic.Unmarshal(actionsByte, &actionsList); err == nil {
			event.Log("soren-sdk-init", models.LogLevelInfo, "start plugin", map[string]any{"actions": actionsList})
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

	return p.Progress(jobId, models.ProgressCommand, models.JobProgress{Progress: 100, Details: data})

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

	sub := p.sdk.makeJobSubject(jobId, string(command))
	if entId, ok := GetjobsHolder().Get(jobId); ok {
		sub = strings.Replace(sub, "*", entId, 1)
	}
	dataByte, err := sonic.Marshal(data)
	if err != nil {
		log.Println("progress command ", command, " error:", err)
		return err
	}
	for retry := range 5 {
		msg, err := p.sdk.conn.Request(sub, dataByte, 3*time.Second)
		if err != nil {
			if err == nats.ErrNoResponders {
				if retry > 2 {
					log.Default().Printf("No responders for progress command:%s - retry :%d", command, retry)
				}
				time.Sleep(time.Duration(retry+1) * time.Second)
				continue

			}
			log.Println("progress command publish error:", err)
			log.Println("jobid : ", jobId)
			log.Println("subs : ", sub)
			log.Println("body : ", string(dataByte))

			return err
		}
		if err := p.sdk.conn.Flush(); err != nil {
			log.Println("progress command flush error:", err)
			return err
		}

		fmt.Printf("result of %s  :  %s \n", sub, string(msg.Data))
		return msg
	}
	if data.Progress == 100 {
		GetjobsHolder().Delete(jobId)
	}
	return nil
}
