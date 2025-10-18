package sdkv2

import (
	"log"
	"strings"

	"github.com/bytedance/sonic"
	"github.com/nats-io/nats.go"
)

type Plugin struct {
	sdk *SorenSDK
	Intro PluginIntro
	Settings *Settings
	Actions []Action
}
func NewPlugin(sdk *SorenSDK) *Plugin{
	return &Plugin{
		sdk: sdk,
	}
}
func (p *Plugin)SetSettings(settings *Settings,handler func([]byte)any){
	p.Settings=settings
	if p.Settings!=nil{
		p.Settings.Handler=handler
	}
}
func (p *Plugin)SetActions(actions []Action){
	p.Actions=actions
}
func (p *Plugin)SetIntro(intro PluginIntro,handler func([]byte)any){
	p.Intro=intro
	if p.Intro.Requirements!=nil{
		p.Intro.Requirements.Handler=handler
	}
}
func (p *Plugin)AddActions(actions []Action){
	p.Actions=append(p.Actions,actions...)
}
func (p *Plugin)Start()error{
	err:=p.IntroHandler()
	if err!=nil{
		return err
	}
	err=p.SettingsHandler()
	if err!=nil{
		return err
	}
	p.ActionsHandler()

	return nil
}
func (p *Plugin)IntroHandler() error{
	p.sdk.conn.Subscribe(p.sdk.makeIntroSubject(),func(msg *nats.Msg) {
		// Handle the intro message
		introByte,err:=sonic.Marshal(p.Intro)
		if err!=nil{
			return 
		}
		msg.Respond(introByte)
	})
	if p.Intro.Requirements!=nil{
		if strings.TrimSpace(p.Intro.Requirements.ReplyTo)==""{
			log.Println("no setting service defined")
			return nil
		}
		p.sdk.conn.Subscribe(p.sdk.makeSubject(p.Intro.Requirements.ReplyTo),func(msg *nats.Msg) {
			if p.Intro.Requirements.Handler==nil{
				msg.Respond([]byte(`{"status":"not implemented"}`))
				return
			}
			result:=p.Intro.Requirements.Handler(msg.Data)
			resByte,err:=sonic.Marshal(result)
			if err!=nil{
				log.Println("submit required info on init plugin error:",err)
				return 
			}
			msg.Respond(resByte)
		})
	}
	return  nil
}

func (p *Plugin)SettingsHandler()  error{
	// show settings form handler
	p.sdk.conn.Subscribe(p.sdk.makeSettingsSubject(),func(msg *nats.Msg) {
		// Handle the settings message
		if p.Settings==nil{
			msg.Respond(nil)
			return
		}
		settingsByte,err:=sonic.Marshal(p.Settings)
		if err!=nil{
			return 
		}
		msg.Respond(settingsByte)
	})
	// settings submit handler 
	if p.Settings!=nil{
		if strings.TrimSpace(p.Settings.ReplyTo)==""{
			log.Println("no setting service defined")
			return nil
		}
		p.sdk.conn.Subscribe(p.sdk.makeSubject(p.Settings.ReplyTo),func(msg *nats.Msg) {
			if p.Settings.Handler==nil{
				msg.Respond([]byte(`{"status":"not implemented"}`))
				return
			}
			result:=p.Settings.Handler(msg.Data)
			resByte,err:=sonic.Marshal(result)
			if err!=nil{
				log.Println("settings handler response error:",err)
				return 
			}
			msg.Respond(resByte)
		})
	}
	
	return nil
}


func (p *Plugin)ActionsHandler() {
	p.sdk.conn.Subscribe(p.sdk.makeActionsListSubject(),func(msg *nats.Msg) {
		// Handle the intro message
		listBytes,err:=sonic.Marshal(p.Actions)
		if err!=nil{
			return 
		}
		msg.Respond(listBytes)
	})
	for _,action:=range p.Actions{
		p.sdk.conn.Subscribe(p.sdk.makeFormSubject(action.Method),func(msg *nats.Msg) {
			// Handle the action message
			formBody,err:=sonic.Marshal(action.Form)
			if err!=nil{
				log.Println("action form ",action.Title," error:",err)
				return 
			}
			msg.Respond(formBody)
		})
		// request handler make a jobId and respond it with the result
		p.sdk.conn.Subscribe(p.sdk.makeSubject(action.Method),func(msg *nats.Msg) {
			result:=action.RequestHandler(msg.Data)
			resByte,err:=sonic.Marshal(result)
			if err!=nil{
				log.Println("action response ",action.Title," error:",err)
				return 
			}
			msg.Respond(resByte)
		})
	}
}


func (p *Plugin)Progress(jobId string,event PluginEvent) {
	sub:=p.sdk.makeJobSubject(jobId,event.Event)
	dataByte,err:=sonic.Marshal(event)
	if err!=nil{
		log.Println("progress event ",event.Event," error:",err)
		return 
	}
	p.sdk.conn.Publish(sub,dataByte)
}

