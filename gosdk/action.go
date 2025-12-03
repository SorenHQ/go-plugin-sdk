package sdkv2

import (
	"context"
	"log"

	"github.com/bytedance/sonic"
	"github.com/nats-io/nats.go"
	"github.com/sorenhq/go-plugin-sdk/gosdk/models"
	"github.com/sorenhq/go-plugin-sdk/logtool"
)

type Plugin struct {
	sdk      *SorenSDK
	Intro    models.PluginIntro
	Settings *models.Settings
	Actions  []models.Action
}

func NewPlugin(sdk *SorenSDK) *Plugin {
	logtool.Init("SOREN-SDK", true)
	return &Plugin{
		sdk: sdk,
	}
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

	<-p.sdk.ctx.Done()
	log.Println("Plugin context done, exiting plugin:", p.Intro.Name)
	return nil
}

func (p *Plugin) Progress(jobId string, command models.Command, data models.JobProgress) {
	sub := p.sdk.makeJobSubject(jobId, string(command))
	dataByte, err := sonic.Marshal(data)
	if err != nil {
		log.Println("progress command ", command, " error:", err)
		return
	}
	p.sdk.conn.Publish(sub, dataByte)
}


