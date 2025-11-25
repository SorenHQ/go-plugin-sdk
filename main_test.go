package main

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/bytedance/sonic"
	"github.com/google/uuid"
	"github.com/joho/godotenv"

	sdkv2 "github.com/sorenhq/go-plugin-sdk/gosdk"
	models "github.com/sorenhq/go-plugin-sdk/gosdk/models"
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
	plugin.SetIntro(models.PluginIntro{
		Name:    "Code Analysis Plugin",
		Version: "1.1.1",
		Author:  "Soren Team",
		// Requirements: &models.Requirements{
		// 	ReplyTo:    "init.config",
		// 	Jsonui:     map[string]any{"type": "Control", "scope": "#/properties/github_access_token"},
		// 	Jsonschema: map[string]any{"properties": map[string]any{"github_access_token": map[string]any{"title":"GitHub Access Token","type": "string"}}},
		// },
	}, nil)

	plugin.SetSettings(&models.Settings{
		ReplyTo: "settings.config.submit",
		Jsonui: map[string]any{
			"type": "VerticalLayout",

			"elements": []map[string]any{
				{
					"type":  "Control",
					"scope": "#/properties/project",
				},
				{
					"type":  "Control",
					"scope": "#/properties/repository_name",
				},
				{
					"type":  "Control",
					"scope": "#/properties/access_token",
				},
			},
		},
		Jsonschema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"project": map[string]any{
					"type":        "string",
					"title":       "Your Project Name",
					"description": "Project Name",
				},
				"repository_name": map[string]any{
					"type":        "string",
					"title":       "Your Repository Name",
					"description": "Github Respository name",
				},

				"access_token": map[string]any{
					"type":        "string",
					"title":       "Fine Grained Access Token",
					"description": "Github FineGrained Access Token",
				},
			},
			"required": []string{"repository_name", "access_token","project"},
		},
	}, settingsUpdateHandler)
	plugin.AddActions([]models.Action{{
		Method: "prepare",
		Title:  "Clone/Pull Repo",
		Form: models.ActionFormBuilder{
			Jsonui:     map[string]any{"type": "Control", "scope": "#/properties/project"},
			Jsonschema: map[string]any{"properties": map[string]any{"project": map[string]any{"enum": makeEnumsProject()}}},
		},
		RequestHandler: func(data []byte) any {
			// for example in this step we register a job in local database or external system - mae a scan in Joern
			uuid,err:=uuid.NewV6()
			if err!=nil{
				return map[string]any{"jobId": uuid.String()}
			}
			return map[string]any{"details": map[string]any{"error": "service unavailable"}}
		},
	
	}, models.Action{
		Method: "scan.gen.graph",
		Title:  "Scan Code And Create Graph",
		Form: models.ActionFormBuilder{
			Jsonui:     map[string]any{"type": "Control", "scope": "#/properties/reponame"},
			Jsonschema: map[string]any{"properties": map[string]any{"reponame": map[string]any{"type": "string"}}},
		},
		RequestHandler: func(data []byte) any {
			// for example in this step we register a job in local database or external system - mae a scan in Joern
			uuid,err:=uuid.NewV6()
			if err!=nil{
				return map[string]any{"jobId": uuid.String()}
			}
			return map[string]any{"details": map[string]any{"error": "service unavailable"}}
		},
	},
	})
	event := sdkv2.NewEventLogger(sdkInstance)
	event.Log("remote-mate-pc", models.LogLevelInfo, "start plugin", nil)
	plugin.Start()
	select {}
}

func settingsUpdateHandler(data []byte) any {
	fmt.Println("New Update As Settings : ", string(data))
	settings:=map[string]any{}
	err:=sonic.Unmarshal(data,&settings)
	if err!=nil{
		fmt.Println("Error Unmarshalling Settings:",err)
		return map[string]any{"status": "error"}
	}
	os.WriteFile("my_database.json",data,0644)
	return map[string]any{"status": "accepted"}
}

func makeEnumsProject() []string {
	contentJson,err:=os.ReadFile("my_database.json")
	if err!=nil{
		return []string{}
	}
	savedSettings:=map[string]any{}
	err=sonic.Unmarshal(contentJson,&savedSettings)
	if err!=nil{
		return []string{}
	}
	if savedSettings["project"] == nil {
		return []string{}
	}
	return []string{savedSettings["project"].(string)}
}

// For Joern Wee Need Github Fine Grain Token and Repo Url
// this value can getting through settings or actions params - based on its nature of entity


