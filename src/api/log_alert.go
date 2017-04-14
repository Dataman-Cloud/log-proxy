package api

import (
	"bytes"
	"container/list"
	"errors"
	"html/template"

	"github.com/Dataman-Cloud/log-proxy/src/config"
	"github.com/Dataman-Cloud/log-proxy/src/models"
	"github.com/Dataman-Cloud/log-proxy/src/utils"

	"github.com/gin-gonic/gin"
)

// CreateAlert create keyword filter
func (s *Search) CreateLogAlertRule(ctx *gin.Context) {
	var alertRule models.LogAlertRule
	if err := ctx.BindJSON(&alertRule); err != nil {
		utils.ErrorResponse(ctx, utils.NewError(ParamError, errors.New("request body param error")))
		return
	}

	if err := verifyLogAlertRule(alertRule); err != nil {
		utils.ErrorResponse(ctx, utils.NewError(ParamError, err))
		return
	}

	ruleIndex := alertRule.App + alertRule.Source
	if s.KeywordFilter[ruleIndex] != nil {
		for e := s.KeywordFilter[ruleIndex].Front(); e != nil; e = e.Next() {
			if e.Value.(models.LogAlertRule).Keyword == alertRule.Keyword {
				utils.ErrorResponse(ctx, utils.NewError(CreateAlertError, errors.New("keyword exist")))
				return
			}
		}
	} else {
		s.KeywordFilter[ruleIndex] = list.New()
	}

	if err := s.Store.CreateLogAlertRule(&alertRule); err != nil {
		utils.ErrorResponse(ctx, utils.NewError(CreateAlertError, err))
		return
	}

	s.Kmutex.Lock()
	s.KeywordFilter[ruleIndex].PushBack(alertRule)
	s.Kmutex.Unlock()

	utils.Ok(ctx, "create success")
	return
}

// DeleteAlert delete keyword filter
func (s *Search) DeleteLogAlertRule(ctx *gin.Context) {
	ruleID := ctx.Param("id")
	alertRule, err := s.Store.GetLogAlertRule(ruleID)
	if err != nil {
		utils.ErrorResponse(ctx, utils.NewError(GetAlertError, err))
		return
	}

	if err := s.Store.DeleteLogAlertRule(ruleID); err != nil {
		utils.ErrorResponse(ctx, utils.NewError(DeleteAlertError, err))
		return
	}

	ruleIndex := alertRule.App + alertRule.Source

	s.Kmutex.Lock()

	if s.KeywordFilter[ruleIndex] == nil {
		utils.Ok(ctx, "delete success")
		return
	} else {
		for e := s.KeywordFilter[ruleIndex].Front(); e != nil; e = e.Next() {
			if e.Value.(models.LogAlertRule).Keyword == alertRule.Keyword {
				s.KeywordFilter[ruleIndex].Remove(e)
			}
		}
	}

	s.Kmutex.Unlock()

	utils.Ok(ctx, "delete success")
}

// GetAlerts get all keyword filter
func (s *Search) GetLogAlertRules(ctx *gin.Context) {
	alertRules, err := s.Store.GetLogAlertRules(ctx.MustGet("page").(models.Page))
	if err != nil {
		utils.ErrorResponse(ctx, utils.NewError(GetAlertError, err))
		return
	}

	utils.Ok(ctx, alertRules)
}

// GetAlert get keyword filter by id
func (s *Search) GetLogAlertRule(ctx *gin.Context) {
	alertRule, err := s.Store.GetLogAlertRule(ctx.Param("id"))
	if err != nil {
		utils.ErrorResponse(ctx, utils.NewError(GetAlertError, err))
		return
	}

	utils.Ok(ctx, alertRule)
}

// UpdateAlert update keyword filter
func (s *Search) UpdateLogAlertRule(ctx *gin.Context) {
	var alertRule models.LogAlertRule
	if err := ctx.BindJSON(&alertRule); err != nil {
		utils.ErrorResponse(ctx, utils.NewError(ParamError, errors.New("request body param error")))
		return
	}

	if err := verifyLogAlertRule(alertRule); err != nil {
		utils.ErrorResponse(ctx, utils.NewError(ParamError, err))
	}

	if err := s.Store.UpdateLogAlertRule(&alertRule); err != nil {
		utils.ErrorResponse(ctx, utils.NewError(UpdateAlertError, err))
		return
	}

	ruleIndex := alertRule.App + alertRule.Source

	s.Kmutex.Lock()
	if s.KeywordFilter[ruleIndex] != nil {
		for e := s.KeywordFilter[ruleIndex].Front(); e != nil; e = e.Next() {
			if e.Value.(models.LogAlertRule).Keyword == alertRule.Keyword {
				s.KeywordFilter[ruleIndex].Remove(e)
				s.KeywordFilter[ruleIndex].PushBack(alertRule)
				break
			}
		}
	} else {
		s.KeywordFilter[ruleIndex] = list.New()
		s.KeywordFilter[ruleIndex].PushBack(alertRule)
	}
	s.Kmutex.Unlock()

	utils.Ok(ctx, "update success")
}

func verifyLogAlertRule(rule models.LogAlertRule) error {
	if rule.App == "" {
		return errors.New("appid can't be empty")
	}

	if rule.Keyword == "" {
		return errors.New("keyword can't be empty")
	}

	if rule.Source == "" {
		return errors.New("source can't be empty")
	}

	return nil
}

func (s *Search) GetLogAlertEvents(ctx *gin.Context) {
	options := make(map[string]interface{})

	if ctx.Query("cluster") != "" {
		options["cluster"] = ctx.Query("cluster")
	}

	if ctx.Query("app") != "" {
		options["app"] = ctx.Query("app")
	}

	if ctx.Query("source") != "" {
		options["path"] = ctx.Query("source")
	}

	if ctx.Query("keyword") != "" {
		options["keyword"] = ctx.Query("keyword")
	}

	if ctx.Query("start") != "" {
		options["start"] = ctx.Query("start")
	}

	if ctx.Query("end") != "" {
		options["end"] = ctx.Query("end")
	}

	events, err := s.Store.GetLogAlertEvents(options, ctx.MustGet("page").(models.Page))
	if err != nil {
		utils.ErrorResponse(ctx, utils.NewError(GetLogAlertEventsError, err))
		return
	}

	utils.Ok(ctx, events)
	return
}

func (s *Search) DeleteLogAlertEvents(ctx *gin.Context) {
	if err := s.Store.DeleteLogAlertEvents(ctx.Query("start"), ctx.Query("end")); err != nil {
		utils.ErrorResponse(ctx, utils.NewError(DeleteLogAlertEventsError, err))
		return
	}

	utils.Ok(ctx, "delete success")
	return
}

func (s *Search) GetLogAlertClusters(ctx *gin.Context) {
	clusters, err := s.Store.GetLogAlertClusters(ctx.Query("start"), ctx.Query("end"))
	if err != nil {
		utils.ErrorResponse(ctx, utils.NewError(GetLogALertClustersError, err))
		return
	}

	utils.Ok(ctx, clusters)
	return
}

func (s *Search) GetLogAlertApps(ctx *gin.Context) {
	apps, err := s.Store.GetLogAlertApps(ctx.Param("cluster"), ctx.Query("start"), ctx.Query("end"))
	if err != nil {
		utils.ErrorResponse(ctx, utils.NewError(GetLogALertAppsError, err))
		return
	}

	utils.Ok(ctx, apps)
	return
}

func (s *Search) HandleLogAlertEvent(ctx *gin.Context) {
	var handleInfo struct {
		Action string `json:action`
	}

	if err := ctx.BindJSON(&handleInfo); err != nil {
		utils.ErrorResponse(ctx, utils.NewError(ParamError, err))
		return
	}

	if handleInfo.Action == "ack" {
		if err := s.Store.AckLogAlertEvent(ctx.Param("id")); err != nil {
			utils.ErrorResponse(ctx, utils.NewError(ParamError, err))
			return
		}

		go s.SendLogAlertAckEventToCama(ctx.Param("id"))
	} else {
		utils.ErrorResponse(ctx, utils.NewError(ParamError, errors.New("invalid action"+handleInfo.Action)))
		return
	}

	utils.Ok(ctx, "success")
	return
}

func (s *Search) ConvertLogAlertToCamaEvent(alertEvent *models.LogAlertEvent, eventType int) (*models.CamaEvent, error) {
	cmdbServer, err := s.Store.GetCmdbServer(alertEvent.App)
	if err != nil {
		return nil, err
	}

	t := template.Must(template.New("camaEventTempl").Parse(camaEventTempl))
	var buf bytes.Buffer
	if err := t.Execute(&buf, alertEvent); err != nil {
		return nil, err
	}

	camaEvent := &models.CamaEvent{
		Channel:   "DOCKER",
		Recover:   eventType,
		ServerNo:  cmdbServer.CmdbAppID,
		Merger:    1,
		EventDesc: buf.String(),
		Level:     5,
		FirstTime: alertEvent.LogTime.Format(config.CamaTimeFormatString),
		LastTime:  alertEvent.LogTime.Format(config.CamaTimeFormatString),
	}

	return camaEvent, nil
}
