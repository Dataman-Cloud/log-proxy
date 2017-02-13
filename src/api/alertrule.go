package api

import (
	"errors"

	"github.com/Dataman-Cloud/log-proxy/src/models"
	"github.com/Dataman-Cloud/log-proxy/src/service"
	"github.com/Dataman-Cloud/log-proxy/src/utils"
	"github.com/Dataman-Cloud/log-proxy/src/utils/database"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type Alert struct {
	DB *gorm.DB
}

func NewAlert() *Alert {
	db := database.GetDB()
	return &Alert{DB: db}
}

// CreateAlertRule create the alert rule in Database
func (alert *Alert) CreateAlertRule(ctx *gin.Context) {

	var rule *models.Rule
	if err := ctx.BindJSON(&rule); err != nil {
		utils.ErrorResponse(ctx, err)
		return
	}

	if v := rule.Name; v == "" {
		utils.ErrorResponse(ctx, errors.New("not found Name string"))
		return
	}

	if v := rule.Alert; v == "" {
		utils.ErrorResponse(ctx, errors.New("not found Alert string"))
		return
	}

	if v := rule.ForTime; v == "" {
		utils.ErrorResponse(ctx, errors.New("not found for_time string"))
		return
	}

	alertRule := &service.AlertRule{
		DB:   alert.DB,
		Rule: rule,
	}

	err := alertRule.CreateAlertRule()
	if err != nil {
		utils.ErrorResponse(ctx, err)
		return
	}
	utils.Ok(ctx, "success")
}

func (alert *Alert) DeleteAlertRule(ctx *gin.Context) {

	alertRule := &service.AlertRule{
		DB: alert.DB,
	}

	err := alertRule.DeleteAlertRule(ctx.Param("id"))
	if err != nil {
		utils.ErrorResponse(ctx, err)
		return
	}
	utils.Ok(ctx, "success")
}

func (alert *Alert) GetAlertRule(ctx *gin.Context) {

	alertRule := &service.AlertRule{
		DB: alert.DB,
	}

	data, err := alertRule.GetAlertRule()
	if err != nil {
		utils.ErrorResponse(ctx, err)
		return
	}

	utils.Ok(ctx, data)
}

// UpdateAlertRule create the alert rule in Database
func (alert *Alert) UpdateAlertRule(ctx *gin.Context) {

	var rule *models.Rule
	if err := ctx.BindJSON(&rule); err != nil {
		utils.ErrorResponse(ctx, err)
		return
	}

	alertRule := &service.AlertRule{
		DB:   alert.DB,
		Rule: rule,
	}

	err := alertRule.UpdateAlertRule()
	if err != nil {
		utils.ErrorResponse(ctx, err)
		return
	}
	utils.Ok(ctx, "success")
}
