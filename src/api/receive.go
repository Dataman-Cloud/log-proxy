package api

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Dataman-Cloud/log-proxy/src/models"
	"github.com/Dataman-Cloud/log-proxy/src/utils"

	"github.com/gin-gonic/gin"
)

func (s *search) Receiver(ctx *gin.Context) {
	data, err := utils.ReadRequestBody(ctx.Request)
	if err != nil {
		utils.ErrorResponse(ctx, utils.NewError(GET_EVENTS_ERROR, err))
		return
	}

	var m map[string]interface{}
	err = json.Unmarshal(data, &m)
	if err != nil {
		utils.ErrorResponse(ctx, utils.NewError(GET_EVENTS_ERROR, err))
		return
	}

	for _, alert := range m["alerts"].([]interface{}) {
		a := alert.(map[string]interface{})
		a["alertname"] = a["labels"].(map[string]interface{})["alertname"]
		labels, err := json.Marshal(a["labels"])
		if err != nil {
			continue
		}
		delete(a, "labels")
		a["labels"] = utils.Byte2str(labels)
		s.Service.SavePrometheus(a)
	}

	utils.Ok(ctx, map[string]string{"status": "success"})
}

func (s *search) ReceiverLog(ctx *gin.Context) {
	data, err := utils.ReadRequestBody(ctx.Request)
	if err != nil {
		utils.ErrorResponse(ctx, utils.NewError(GET_LOG_ERROR, err))
		return
	}

	var m map[string]interface{}
	err = json.Unmarshal(data, &m)
	if err != nil {
		utils.ErrorResponse(ctx, utils.NewError(GET_LOG_ERROR, err))
		return
	}

	appid, ok := m["appid"]
	if !ok {
		if err != nil {
			utils.ErrorResponse(ctx, utils.NewError(GET_LOG_ERROR, err))
			return
		}
	}

	taskid, ok := m["taskid"]
	if !ok {
		if err != nil {
			utils.ErrorResponse(ctx, utils.NewError(GET_LOG_ERROR, err))
			return
		}
	}

	path, ok := m["path"]
	if !ok {
		if err != nil {
			utils.ErrorResponse(ctx, utils.NewError(GET_LOG_ERROR, err))
			return
		}
	}

	userid, ok := m["userid"]
	if !ok {
		if err != nil {
			utils.ErrorResponse(ctx, utils.NewError(GET_LOG_ERROR, err))
			return
		}
	}

	clusterid, ok := m["clusterid"]
	if !ok {
		if err != nil {
			utils.ErrorResponse(ctx, utils.NewError(GET_LOG_ERROR, err))
			return
		}
	}

	offset, ok := m["offset"]
	if !ok {
		if err != nil {
			utils.ErrorResponse(ctx, utils.NewError(GET_LOG_ERROR, err))
			return
		}
	}

	message, ok := m["message"]
	if !ok {
		if err != nil {
			utils.ErrorResponse(ctx, utils.NewError(GET_LOG_ERROR, err))
			return
		}
	}

	keywords, ok := s.KeywordFilter[appid.(string)+path.(string)]
	fmt.Println("---------", string(data), "--------", keywords)
	if !ok {
		utils.Ok(ctx, "ok")
		return
	}

	for _, keyword := range keywords {
		if strings.Index(message.(string), keyword) == -1 {
			continue
		}

		s.Counter.WithLabelValues(
			appid.(string),
			taskid.(string),
			path.(string),
			keyword,
			userid.(string),
			clusterid.(string),
			fmt.Sprint(int64(offset.(float64)))).Inc()
	}

	utils.Ok(ctx, "ok")
	return
}

func (s *search) GetPrometheus(ctx *gin.Context) {
	result, err := s.Service.GetPrometheus(ctx.MustGet("page").(models.Page))
	if err != nil {
		utils.ErrorResponse(ctx, utils.NewError(GET_PROMETHEUS_ERROR, err))
		return
	}

	utils.Ok(ctx, result)
}

func (s *search) GetPrometheu(ctx *gin.Context) {
	result, err := s.Service.GetPrometheu(ctx.Param("id"))
	if err != nil {
		utils.ErrorResponse(ctx, utils.NewError(GET_PROMETHEUS_ERROR, err))
		return
	}

	utils.Ok(ctx, result)
}
