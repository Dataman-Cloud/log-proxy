package api

import (
	"encoding/json"

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

	var pro models.Prometheus
	err = json.Unmarshal(data, &pro)
	if err != nil {
		utils.ErrorResponse(ctx, utils.NewError(GET_EVENTS_ERROR, err))
		return
	}
	pro.CLs.Condition = pro.CAs.Summary
	pro.CLs.Usage = pro.CAs.Description

	s.Service.SavePrometheus(pro.CLs)
	utils.Ok(ctx, map[string]string{"status": "success"})
}

func (s *search) GetPrometheus(ctx *gin.Context) {
	result, err := s.Service.GetPrometheus(ctx.MustGet("page").(models.Page))
	if err != nil {
		utils.ErrorResponse(ctx, utils.NewError(GET_PROMETHEUS_ERROR, err))
		return
	}

	utils.Ok(ctx, result)
}
