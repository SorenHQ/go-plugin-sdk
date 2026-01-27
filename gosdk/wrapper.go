package sdkv2

import (
	"fmt"
	"strings"

	"github.com/bytedance/sonic"
	"github.com/nats-io/nats.go"
	uuid "github.com/nu7hatch/gouuid"
	"github.com/sorenhq/go-plugin-sdk/gosdk/models"
)

// Accept Request , make a request session and return sessionId - jobId
func Accept(msg *nats.Msg) (jobId string) {
	jobBody := models.JobBodyContent{}
	uuid, err := uuid.NewV4()
	if err != nil {
		return ""
	}
	if strings.HasPrefix(GetPlugin().sdk.pluginID, "bin.*.") {
		// Extract EntityId(spaceId) part after "bin.*."
		parts := strings.Split(msg.Subject, ".")
		if len(parts) >= 3 {
			requesterSpaceId := parts[3]
			GetjobsHolder().Add(uuid.String(), requesterSpaceId)
		}
	}
	jobBody.JobId = uuid.String()
	responseByte, err := sonic.Marshal(jobBody)
	if err != nil {
		return ""
	}
	err = msg.Respond(responseByte)
	
	return uuid.String()
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
