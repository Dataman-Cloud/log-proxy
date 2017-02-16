package api

import (
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/Dataman-Cloud/log-proxy/src/config"
	"github.com/Dataman-Cloud/log-proxy/src/models"
	"github.com/Dataman-Cloud/log-proxy/src/service"
	"github.com/Dataman-Cloud/log-proxy/src/utils"

	"github.com/gin-gonic/gin"
)

const (
	// ParamError param error
	ParamError = "400-11001"
	// GetAppsError get application error
	GetAppsError = "503-11000"
	// GetTaskError get task error
	GetTaskError = "503-11002"
	// IndexError search index log error
	IndexError = "503-11003"
	// GetEventsError get event history error
	GetEventsError = "503-11008"
	// GetPrometheusError get prometheus event error
	GetPrometheusError = "503-11009"
)

// Search search client struct
type Search struct {
	Service *service.SearchService
}

// NewSearch new search client
func NewSearch() *Search {
	s := &Search{
		Service: service.NewEsService(strings.Split(config.GetConfig().EsURL, ",")),
	}

	return s
}

// Ping ping
func (s *Search) Ping(ctx *gin.Context) {
	utils.Ok(ctx, "success")
}

// Applications get all applications
func (s *Search) Applications(ctx *gin.Context) {
	apps, err := s.Service.Applications(ctx.MustGet("page").(models.Page))
	if err != nil {
		utils.ErrorResponse(ctx, utils.NewError(GetAppsError, err))
		return
	}
	utils.Ok(ctx, apps)
}

// Tasks search applications tasks
func (s *Search) Tasks(ctx *gin.Context) {

	tasks, err := s.Service.Tasks(ctx.Param("app"), ctx.Query("user"), ctx.MustGet("page").(models.Page))
	if err != nil {
		utils.ErrorResponse(ctx, utils.NewError(GetTaskError, err))
		return
	}
	utils.Ok(ctx, tasks)
}

// Paths search applications paths
func (s *Search) Paths(ctx *gin.Context) {
	paths, err := s.Service.Paths(
		ctx.Query("cluster"),
		ctx.Query("user"),
		ctx.Param("app"),
		ctx.Query("task"),
		ctx.MustGet("page").(models.Page),
	)
	if err != nil {
		utils.ErrorResponse(ctx, utils.NewError(GetTaskError, err))
		return
	}
	utils.Ok(ctx, paths)
}

// Index search log by condition
func (s *Search) Index(ctx *gin.Context) {
	if ctx.Query("app") == "" {
		utils.ErrorResponse(ctx, utils.NewError(ParamError, errors.New("app can't be empty")))
		return
	}

	results, err := s.Service.Search(
		ctx.Query("cluster"),
		ctx.Query("user"),
		ctx.Query("app"),
		ctx.Query("task"),
		ctx.Query("path"),
		ctx.Query("keyword"),
		ctx.MustGet("page").(models.Page))
	if err != nil {
		utils.ErrorResponse(ctx, utils.NewError(IndexError, err))
		return
	}

	utils.Ok(ctx, results)
}

// Context search log context
func (s *Search) Context(ctx *gin.Context) {
	if ctx.Query("app") == "" {
		utils.ErrorResponse(ctx, utils.NewError(ParamError, errors.New("app can't be empty")))
		return
	}

	if ctx.Query("task") == "" {
		utils.ErrorResponse(ctx, utils.NewError(ParamError, errors.New("task can't be empty")))
		return
	}

	if ctx.Query("path") == "" {
		utils.ErrorResponse(ctx, utils.NewError(ParamError, errors.New("path can't be empty")))
		return
	}

	if ctx.Query("offset") == "" {
		utils.ErrorResponse(ctx, utils.NewError(ParamError, errors.New("offset can't be empty")))
		return
	}

	results, err := s.Service.Context(
		ctx.Query("cluster"),
		ctx.Query("user"),
		ctx.Query("app"),
		ctx.Query("task"),
		ctx.Query("path"),
		ctx.Query("offset"),
		ctx.MustGet("page").(models.Page))
	if err != nil {
		utils.ErrorResponse(ctx, utils.NewError(IndexError, err))
		return
	}

	utils.Ok(ctx, results)
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
