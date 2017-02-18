package api

import (
	"encoding/json"
	"errors"
	"strconv"

	"github.com/Dataman-Cloud/log-proxy/src/models"
	"github.com/Dataman-Cloud/log-proxy/src/service"
	"github.com/Dataman-Cloud/log-proxy/src/service/alertevent"
	"github.com/Dataman-Cloud/log-proxy/src/utils"
	"github.com/Dataman-Cloud/log-proxy/src/utils/database"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

const (
	// Receive Alert event error
	ReceiveEventError = "503-21000"
	// Ack Alert event error
	AckEventError = "503-21001"
)

type Alert struct {
	DB           *gorm.DB
	eventManager *alertevent.EventManager
}

func NewAlert() *Alert {
	db := database.GetDB()
	return &Alert{
		DB:           db,
		eventManager: &alertevent.EventManager{DB: db},
	}
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

func (alert *Alert) ReceiveAlertEvent(ctx *gin.Context) {
	data, err := utils.ReadRequestBody(ctx.Request)
	if err != nil {
		utils.ErrorResponse(ctx, utils.NewError(ReceiveEventError, err))
		return
	}

	var m map[string]interface{}
	err = json.Unmarshal(data, &m)
	if err != nil {
		utils.ErrorResponse(ctx, utils.NewError(ReceiveEventError, err))
		return
	}

	for _, item := range m["alerts"].([]interface{}) {
		labels := item.(map[string]interface{})["labels"].(map[string]interface{})
		annotations := item.(map[string]interface{})["annotations"].(map[string]interface{})
		event := &models.AlertEvent{
			AlertName:   labels["alertname"].(string),
			Severity:    labels["severity"].(string),
			VCluster:    labels["container_label_VCLUSTER"].(string),
			App:         labels["container_label_APP"].(string),
			Slot:        labels["container_label_SLOT"].(string),
			ContainerId: labels["id"].(string),
			Description: annotations["description"].(string),
			Summary:     annotations["summary"].(string),
		}
		if err := alert.eventManager.CreateOrIncreaseEvent(event); err != nil {
			utils.ErrorResponse(ctx, utils.NewError(ReceiveEventError, err))
			return
		}
	}

	utils.Ok(ctx, map[string]string{"status": "success"})
}

func (alert *Alert) AckAlertEvent(ctx *gin.Context) {
	pk, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		utils.ErrorResponse(ctx, utils.NewError(AckEventError, err))
		return
	}
	var data map[string]interface{}
	if err := ctx.BindJSON(&data); err != nil {
		utils.ErrorResponse(ctx, utils.NewError(AckEventError, err))
		return
	}

	switch action := data["action"].(string); action {
	case "ack":
		if err = alert.eventManager.AckEvent(pk); err != nil {
			utils.ErrorResponse(ctx, utils.NewError(AckEventError, err))
			return
		}
		utils.Ok(ctx, map[string]string{"status": "success"})
	}
}

func (alert *Alert) GetAlertEvents(ctx *gin.Context) {
	switch ack := ctx.Query("ack"); ack {
	case "true":
		result := alert.eventManager.ListAckedEvent(ctx.MustGet("page").(models.Page))
		utils.Ok(ctx, result)
	case "false", "":
		result := alert.eventManager.ListUnackedEvent(ctx.MustGet("page").(models.Page))
		utils.Ok(ctx, result)
	}
}
