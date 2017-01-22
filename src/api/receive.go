package api

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/Dataman-Cloud/log-proxy/src/config"
	"github.com/Dataman-Cloud/log-proxy/src/models"
	"github.com/Dataman-Cloud/log-proxy/src/utils"

	"github.com/gin-gonic/gin"
)

// Receiver receive prometheus alert event
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

// ReceiverLog receive log
func (s *Search) ReceiverLog(ctx *gin.Context) {
	data, err := utils.ReadRequestBody(ctx.Request)
	if err != nil {
		utils.ErrorResponse(ctx, utils.NewError(GetLogError, err))
		return
	}

	var m map[string]interface{}
	err = json.Unmarshal(data, &m)
	if err != nil {
		utils.ErrorResponse(ctx, utils.NewError(GetLogError, err))
		return
	}

	app, ok := m["app"]
	if !ok {
		utils.ErrorResponse(ctx, utils.NewError(GetLogError, "not found app"))
		return
	}

	task, ok := m["task"]
	if !ok {
		utils.ErrorResponse(ctx, utils.NewError(GetLogError, "not found task"))
		return
	}

	path, ok := m["path"]
	if !ok {
		utils.ErrorResponse(ctx, utils.NewError(GetLogError, "not found path"))
		return
	}

	user, ok := m["user"]
	if !ok {
		utils.ErrorResponse(ctx, utils.NewError(GetLogError, "not found user"))
		return
	}

	cluster, ok := m["cluster"]
	if !ok {
		utils.ErrorResponse(ctx, utils.NewError(GetLogError, "not found cluster"))
		return
	}

	message, ok := m["message"]
	if !ok {
		utils.ErrorResponse(ctx, utils.NewError(GetLogError, "not found message"))
		return
	}

	keywords, ok := s.KeywordFilter[app.(string)+path.(string)]
	if !ok {
		utils.Ok(ctx, "ok")
		return
	}
	for e := keywords.Front(); e != nil; e = e.Next() {
		if strings.Index(message.(string), e.Value.(string)) == -1 {
			continue
		}

		s.Counter.WithLabelValues(
			app.(string),
			task.(string),
			path.(string),
			e.Value.(string),
			user.(string),
			cluster.(string),
		).Inc()
	}

	utils.Ok(ctx, "ok")
	return
}

// GetPrometheus get all prometheus
func (s *Search) GetPrometheus(ctx *gin.Context) {
	result, err := s.Service.GetPrometheus(ctx.MustGet("page").(models.Page))
	if err != nil {
		utils.ErrorResponse(ctx, utils.NewError(GetPrometheusError, err))
		return
	}

	utils.Ok(ctx, result)
}

// GetPrometheu get prometheus by id
func (s *Search) GetPrometheu(ctx *gin.Context) {
	result, err := s.Service.GetPrometheu(ctx.Param("id"))
	if err != nil {
		utils.ErrorResponse(ctx, utils.NewError(GetPrometheusError, err))
		return
	}

	utils.Ok(ctx, result)
}
