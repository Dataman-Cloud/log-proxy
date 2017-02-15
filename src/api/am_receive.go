package api

import (
	"encoding/json"
	"time"

	"github.com/Dataman-Cloud/log-proxy/src/config"
	"github.com/Dataman-Cloud/log-proxy/src/utils"

	"github.com/gin-gonic/gin"
)

// Receiver receive alert event sourced from alertmanager, output to elasticsearch
func (s *Search) Receiver(ctx *gin.Context) {
	data, err := utils.ReadRequestBody(ctx.Request)
	if err != nil {
		utils.ErrorResponse(ctx, utils.NewError(GetEventsError, err))
		return
	}

	var m map[string]interface{}
	err = json.Unmarshal(data, &m)
	if err != nil {
		utils.ErrorResponse(ctx, utils.NewError(GetEventsError, err))
		return
	}

	for _, alert := range m["alerts"].([]interface{}) {
		a := alert.(map[string]interface{})
		a["alertname"] = a["labels"].(map[string]interface{})["alertname"]
		labels, err := json.Marshal(a["labels"])

		if config.GetConfig().NotificationURL != "" {
			utils.AlertNotification(config.GetConfig().NotificationURL, map[string]interface{}{
				"alarminfo": []map[string]interface{}{
					map[string]interface{}{
						"level":         a["labels"].(map[string]interface{})["severity"],
						"modelIdentify": a["alertname"],
						"content":       utils.Byte2str(labels),
						"alarmTime":     time.Now().Format(time.RFC3339Nano),
					},
				},
				"apiVersion": time.Now().Format("2006-01-02"),
			})
		}

		if err != nil {
			continue
		}
		delete(a, "labels")
		a["labels"] = utils.Byte2str(labels)
		s.Service.SavePrometheus(a)
	}

	utils.Ok(ctx, map[string]string{"status": "success"})
}
