package sdkv2

import (
	"log"

	"github.com/bytedance/sonic"
)

type Plugin struct {
	sdk      *SorenSDK
	Intro    PluginIntro
	Settings *Settings
	Actions  []Action
}

func NewPlugin(sdk *SorenSDK) *Plugin {
	return &Plugin{
		sdk: sdk,
	}
}
func (p *Plugin) SetSettings(settings *Settings, handler func([]byte) any) {
	p.Settings = settings
	if p.Settings != nil {
		p.Settings.Handler = handler
	}
}
func (p *Plugin) SetActions(actions []Action) {
	p.Actions = actions
}
func (p *Plugin) SetIntro(intro PluginIntro, handler func([]byte) any) {
	p.Intro = intro
	if p.Intro.Requirements != nil {
		p.Intro.Requirements.Handler = handler
	}
}
func (p *Plugin) AddActions(actions []Action) {
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

func (p *Plugin) Progress(jobId string, event PluginEvent) {
	sub := p.sdk.makeJobSubject(jobId, event.Event)
	dataByte, err := sonic.Marshal(event)
	if err != nil {
		log.Println("progress event ", event.Event, " error:", err)
		return
	}
	p.sdk.conn.Publish(sub, dataByte)
}
