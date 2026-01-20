package sdkv2

import (
	"fmt"
	"strings"

	"github.com/bytedance/sonic"
	"github.com/nats-io/nats.go"
	uuid "github.com/nu7hatch/gouuid"
	"github.com/sorenhq/go-plugin-sdk/gosdk/models"
)

// extractEntityIdFromSubject extracts entityId (spaceId) from NATS message subject
// Subject pattern: soren.v2.bin.{entityId}.{pluginId}.{action}
func extractEntityIdFromSubject(subject string) string {
	parts := strings.Split(subject, ".")
	// Look for "bin" in the subject, entityId should be right after it
	for i, part := range parts {
		if part == "bin" && i+1 < len(parts) {
			return parts[i+1]
		}
	}
	return ""
}

// Accept Request , make a request session and return sessionId - jobId
func Accept(msg *nats.Msg) (jobId string) {
	jobBody := models.JobBodyContent{}
	uuid, err := uuid.NewV4()
	if err != nil {
		return ""
	}
	jobId = uuid.String()
	jobBody.JobId = jobId
	responseByte, err := sonic.Marshal(jobBody)
	if err != nil {
		return ""
	}
	err = msg.Respond(responseByte)
	
	// Extract entityId from subject and store it for this job
	entityId := extractEntityIdFromSubject(msg.Subject)
	if entityId != "" {
		plugin := GetPlugin()
		if plugin != nil {
			plugin.StoreEntityIdForJob(jobId, entityId)
		}
	}
	
	return jobId
}
func RejectWithBody(msg *nats.Msg, body map[string]any) {
	responseBody := models.JobBodyContent{Details: map[string]any{"error": body}}
	responseByte, err := sonic.Marshal(responseBody)
	if err != nil {
		fmt.Println(err)
		return
	}
	msg.Respond(responseByte)
}
// for multi plugin handler
func GetPluginById(pluginId string) *Plugin {
	if p, ok := GetPluginHolder().get(pluginId); ok {
		return p
	}
	return nil
}
// Get First Registered Plugin
func GetPlugin() *Plugin {
	h := GetPluginHolder()
	for _, v := range h.holder {
		return v
	}
	return nil

}
