package api

import (
	"strings"

	"github.com/Dataman-Cloud/log-proxy/src/models"
	"github.com/Dataman-Cloud/log-proxy/src/service"
	"github.com/Dataman-Cloud/log-proxy/src/utils"

	"github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
)

const (
	LogAlertRuleEnabled  = "Enabled"
	LogAlertRuleDisabled = "Disabled"
	camaEventTempl       = `{{ .Cluster }}集群的应用:{{ .App }}触发日志报警, 日志信息: {{.Message}}`
)

// ReceiverLog receive log data from logstash
func (s *Search) ReceiverLog(ctx *gin.Context) {
	var event models.LogAlertEvent

	if err := ctx.BindJSON(&event); err != nil {
		logrus.Errorf("Unmarshal log alert event got error: %s", err.Error())
		utils.ErrorResponse(ctx, utils.NewError(GetLogError, err))
		return
	}

	keywords, ok := s.KeywordFilter[event.App+event.Path]
	if !ok {
		utils.Ok(ctx, "ok")
		return
	}

	for e := keywords.Front(); e != nil; e = e.Next() {
		rule := e.Value.(models.LogAlertRule)
		if rule.Status == LogAlertRuleDisabled {
			continue
		}

		keyword := rule.Keyword
		if strings.Index(event.Message, keyword) == -1 {
			continue
		}

		event.Keyword = keyword
		if err := s.Store.CreateLogAlertEvent(&event); err != nil {
			logrus.Errorf("create log alert event got error: %s", err.Error())
		}

		event.Description = rule.Description

		go s.SendLogAlertEventToCama(&event)

	}

	utils.Ok(ctx, "ok")
	return
}

func (s *Search) SendLogAlertEventToCama(event *models.LogAlertEvent) {
	camaEvent, err := s.ConvertLogAlertToCamaEvent(event, 0)
	if err != nil {
		logrus.Errorf("conver event to cama event failed. Error: %s", err.Error())
		return
	}

	service.SendCamaEvent(camaEvent)
}

func (s *Search) SendLogAlertAckEventToCama(ID string) {
	event, err := s.Store.GetLogAlertEvent(ID)
	if err != nil {
		logrus.Errorf("get event failed. Error: %s", err.Error())
		return
	}

	camaEvent, err := s.ConvertLogAlertToCamaEvent(event, 1)
	if err != nil {
		logrus.Errorf("conver event to cama event failed. Error: %s", err.Error())
		return
	}

	service.SendCamaEvent(camaEvent)
}
