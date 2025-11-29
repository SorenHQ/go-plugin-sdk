package sdkv2

import (
	"context"
	"log"

	"github.com/bytedance/sonic"
	"github.com/sorenhq/go-plugin-sdk/gosdk/models"
)

type Plugin struct {
	sdk      *SorenSDK
	Intro    models.PluginIntro
	Settings *models.Settings
	Actions  []models.Action
}

func NewPlugin(sdk *SorenSDK) *Plugin {
	return &Plugin{
		sdk: sdk,
	}
}
func (p *Plugin) GetContext() context.Context {
	return p.sdk.ctx
}
func (p *Plugin) SetSettings(settings *models.Settings, handler func([]byte) any) {
	p.Settings = settings
	if p.Settings != nil {
		p.Settings.Handler = handler
	}
}
func (p *Plugin) SetActions(actions []models.Action) {
	p.Actions = actions
}
func (p *Plugin) SetIntro(intro models.PluginIntro, handler func([]byte) any) {
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


