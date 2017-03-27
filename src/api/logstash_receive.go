package api

import (
	"strings"

	"github.com/Dataman-Cloud/log-proxy/src/models"
	"github.com/Dataman-Cloud/log-proxy/src/utils"

	"github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
)

const (
	LogAlertRuleEnabled  = "Enabled"
	LogAlertRuleDisabled = "Disabled"
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
	}

	utils.Ok(ctx, "ok")
	return
}
