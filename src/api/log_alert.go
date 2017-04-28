package api

import (
	"github.com/Dataman-Cloud/log-proxy/src/models"
	"github.com/Dataman-Cloud/log-proxy/src/utils"

	"github.com/gin-gonic/gin"
)

func (s *Search) CreateLogAlertRule(ctx *gin.Context) {
	var rule models.LogAlertRule
	if err := ctx.BindJSON(&rule); err != nil {
		utils.ErrorResponse(ctx, utils.NewError(CreateLogAlertRuleError, err))
		return
	}

	if err := s.Store.CreateLogAlertRule(&rule); err != nil {
		utils.ErrorResponse(ctx, utils.NewError(CreateLogAlertRuleError, err))
		return
	}

	utils.Ok(ctx, rule)
	return
}
