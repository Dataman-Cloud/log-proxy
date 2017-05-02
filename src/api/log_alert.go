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

func (s *Search) UpdateLogAlertRule(ctx *gin.Context) {
	var rule models.LogAlertRule
	if err := ctx.BindJSON(&rule); err != nil {
		utils.ErrorResponse(ctx, utils.NewError(UpdateLogAlertRuleError, err))
		return
	}

	if err := s.Store.UpdateLogAlertRule(&rule); err != nil {
		utils.ErrorResponse(ctx, utils.NewError(UpdateLogAlertRuleError, err))
		return
	}

	utils.Ok(ctx, rule)
	return
}

func (s *Search) GetLogAlertRule(ctx *gin.Context) {
	ruleID := ctx.Param("id")
	rule, err := s.Store.GetLogAlertRule(ruleID)
	if err != nil {
		utils.ErrorResponse(ctx, utils.NewError(GetLogAlertRuleError, err))
		return
	}

	utils.Ok(ctx, rule)
	return
}

func (s *Search) GetLogAlertRules(ctx *gin.Context) {
	options := make(map[string]interface{})
	if ctx.Query("group") != "" {
		options["group"] = ctx.Query("group")
	}

	rules, err := s.Store.GetLogAlertRules(options, ctx.MustGet("page").(models.Page))
	if err != nil {
		utils.ErrorResponse(ctx, utils.NewError(GetLogAlertRuleError, err))
		return
	}

	utils.Ok(ctx, rules)
	return
}
