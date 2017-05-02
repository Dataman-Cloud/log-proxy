package api

import (
	"container/list"
	"errors"

	"github.com/Dataman-Cloud/log-proxy/src/models"
	"github.com/Dataman-Cloud/log-proxy/src/utils"

	"github.com/gin-gonic/gin"
)

func (s *Search) CreateLogAlertRule(ctx *gin.Context) {
	var rule models.LogAlertRule
	// NOTE: BindJSON required:binding can make sure required is not empty
	if err := ctx.BindJSON(&rule); err != nil {
		utils.ErrorResponse(ctx, utils.NewError(CreateLogAlertRuleError, err))
		return
	}

	ruleIndex := getLogAlertRuleIndex(rule)
	if s.KeywordFilter[ruleIndex] != nil {
		for e := s.KeywordFilter[ruleIndex].Front(); e != nil; e = e.Next() {
			if e.Value.(models.LogAlertRule).Keyword == rule.Keyword {
				utils.ErrorResponse(ctx, utils.NewError(CreateLogAlertRuleError, errors.New("duplicate keyword")))
				return
			}
		}
	} else {
		s.KeywordFilter[ruleIndex] = list.New()
	}

	if err := s.Store.CreateLogAlertRule(&rule); err != nil {
		utils.ErrorResponse(ctx, utils.NewError(CreateLogAlertRuleError, err))
		return
	}

	s.Kmutex.Lock()
	s.KeywordFilter[ruleIndex].PushBack(rule)
	s.Kmutex.Unlock()

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

	ruleIndex := getLogAlertRuleIndex(rule)

	s.Kmutex.Lock()
	defer s.Kmutex.Unlock()
	if s.KeywordFilter[ruleIndex] != nil {
		for e := s.KeywordFilter[ruleIndex].Front(); e != nil; e = e.Next() {
			if e.Value.(models.LogAlertRule).Keyword == rule.Keyword {
				s.KeywordFilter[ruleIndex].Remove(e)
				s.KeywordFilter[ruleIndex].PushBack(rule)
				break
			}
		}
	} else {
		s.KeywordFilter[ruleIndex] = list.New()
		s.KeywordFilter[ruleIndex].PushBack(rule)
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

func (s *Search) DeleteLogAlertRule(ctx *gin.Context) {
	ruleID := ctx.Param("id")
	rule, err := s.Store.GetLogAlertRule(ruleID)
	if err != nil {
		utils.ErrorResponse(ctx, utils.NewError(GetLogAlertRuleError, err))
		return
	}

	if err := s.Store.DeleteLogAlertRule(ruleID); err != nil {
		utils.ErrorResponse(ctx, utils.NewError(DeleteLogAlertRuleError, err))
		return
	}

	s.Kmutex.Lock()
	defer s.Kmutex.Unlock()
	ruleIndex := getLogAlertRuleIndex(rule)
	if s.KeywordFilter[ruleIndex] == nil {
		utils.Ok(ctx, "delete success")
		return
	} else {
		for e := s.KeywordFilter[ruleIndex].Front(); e != nil; e = e.Next() {
			if e.Value.(models.LogAlertRule).Keyword == rule.Keyword {
				s.KeywordFilter[ruleIndex].Remove(e)
			}
		}
	}

	utils.Ok(ctx, "success")
	return
}
