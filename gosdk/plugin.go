package sdkv2

import (
	"log"
	"strings"

	"github.com/bytedance/sonic"
	"github.com/nats-io/nats.go"
	"github.com/sorenhq/go-plugin-sdk/logtool"
)


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
			p.Intro.Requirements.Handler(msg)
			// result:=
			// resByte,err:=sonic.Marshal(result)
			// if err!=nil{
			// 	log.Println("submit required info on init plugin error:",err)
			// 	return 
			// }
			// msg.Respond(resByte)
		})
	}
	return  nil
}

func (p *Plugin)SettingsHandler()  error{
	// show settings form handler
	p.sdk.conn.Subscribe(p.sdk.makeSettingsSubject(),func(msg *nats.Msg) {
		logtool.GetLogger().Info("Settings Called")
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
			p.Settings.ReplyTo = "_settings.config.submit"
			// log.Println("no setting service defined")
			// return nil
		}
		p.sdk.conn.Subscribe(p.sdk.makeSubject(p.Settings.ReplyTo),func(msg *nats.Msg) {
			if p.Settings.Handler==nil{
				msg.Respond([]byte(`{"status":"not implemented"}`))
				return
			}
			p.Settings.Handler(msg)
			// resByte,err:=sonic.Marshal(result)
			// if err!=nil{
			// 	log.Println("settings handler response error:",err)
			// 	return 
			// }
			// msg.Respond(resByte)
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
		_,err:=p.sdk.conn.Subscribe(p.sdk.makeFormSubject(action.Method),func(msg *nats.Msg) {
			// Handle the action message
			formBody,err:=sonic.Marshal(action.Form)
			if err!=nil{
				log.Println("action form ",action.Title," error:",err)
				return 
			}
			msg.Respond(formBody)
		})
		if err!=nil{
			log.Printf("subscribe error: %s on %s\n",err.Error(),p.sdk.makeFormSubject(action.Method))
			return 
		}
		log.Printf("Form Builder Service : %s",p.sdk.makeFormSubject(action.Method))
		// request handler make a jobId and respond it with the result
		_,err=p.sdk.conn.Subscribe(p.sdk.makeActionCpu(action.Method),func(msg *nats.Msg) {
			action.RequestHandler(msg)
			// result:=
			// resByte,err:=sonic.Marshal(result)
			// if err!=nil{
			// 	log.Println("action response ",action.Title," error:",err)
			// 	return 
			// }
			// msg.Respond(resByte)
		})
		if err!=nil{
			log.Printf("subscribe error: %s on %s\n",err.Error(),p.sdk.makeActionCpu(action.Method))
			return 
		}
		log.Printf("Subscribed Action : %s",p.sdk.makeActionCpu(action.Method))
	}
}


