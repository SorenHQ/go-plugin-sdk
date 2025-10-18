package main

import (
	"fmt"
	"log"
	"testing"

	"github.com/joho/godotenv"

	sdkv2 "github.com/sorenhq/go-plugin-sdk/gosdk"
)

func TestMain(t *testing.T) {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println(err)
	}
	sdkInstance, err := sdkv2.NewFromEnv()
	if err != nil {
		log.Fatalf("Failed to create SDK: %v", err)
	}
	defer sdkInstance.Close()
	plugin := sdkv2.NewPlugin(sdkInstance)
	plugin.SetIntro(sdkv2.PluginIntro{
		Name:    "Code Analysis Plugin",
		Version: "1.1.0",
		Author:  "Mate-Soren",
		Requirements: &sdkv2.Requirements{
			ReplyTo:    fmt.Sprintf("soren.v2.%s.intro.requirements", sdkInstance.GetPluginID()),
			Jsonui:     map[string]any{"type": "Control", "scope": "#/properties/apiKey"},
			Jsonschema: map[string]any{"properties": map[string]any{"apiKey": map[string]any{"type": "string"}}},
		},
	},nil)

	plugin.SetSettings(&sdkv2.Settings{
		ReplyTo: "settings.config.submit",
		Jsonui: map[string]any{
			"type":  "Control",
			"scope": "#/properties/start",
		},
		Jsonschema: map[string]any{
			"properties": map[string]any{
				"start": map[string]any{
					"type":        "string",
					"title":       "Start Path",
					"description": "The path to start the analysis from",
				},
			},
			"required": []string{"start"},
		},
	}, settingsUpdateHandler)
	plugin.AddActions([]sdkv2.Action{{
		Method: "analyse.code",
		Title:  "Code Analyser",
		Form: sdkv2.ActionFormBuilder{
			Jsonui:     map[string]any{"type": "Control", "scope": "#/properties/reponame"},
			Jsonschema: map[string]any{"properties": map[string]any{"reponame": map[string]any{"type": "string"}}},
		},
		RequestHandler: func(data []byte) any {
			// for example in this step we register a job in local database or external system - mae a scan in Joern
			return map[string]any{"jobId": "AAAAA-2222"}
		},
	}, sdkv2.Action{
		Method: "scan.code",
		Title:  "Code Scanner",
		Form: sdkv2.ActionFormBuilder{
			Jsonui:     map[string]any{"type": "Control", "scope": "#/properties/reponame"},
			Jsonschema: map[string]any{"properties": map[string]any{"reponame": map[string]any{"type": "string"}}},
		},
		RequestHandler: func(data []byte) any {
			// for example in this step we register a job in local database or external system - mae a scan in Joern
			return map[string]any{"jobId": "B-AAAAA-33321022"}
		},
	}})
	plugin.Start()
	select {}
}


func settingsUpdateHandler(data []byte) any {
	fmt.Println("New Update As Settings : ",string(data))
	return map[string]any{"status": "accepted"}
}


