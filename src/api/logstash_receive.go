package api

import (
	"encoding/json"
	"strings"

	"github.com/Dataman-Cloud/log-proxy/src/utils"

	"github.com/gin-gonic/gin"
)

// ReceiverLog receive log data from logstash
func (s *Search) ReceiverLog(ctx *gin.Context) {
	data, err := utils.ReadRequestBody(ctx.Request)
	if err != nil {
		utils.ErrorResponse(ctx, utils.NewError(GetLogError, err))
		return
	}

	var m map[string]interface{}
	err = json.Unmarshal(data, &m)
	if err != nil {
		utils.ErrorResponse(ctx, utils.NewError(GetLogError, err))
		return
	}

	app, ok := m["app"]
	if !ok {
		utils.ErrorResponse(ctx, utils.NewError(GetLogError, "not found app"))
		return
	}

	task, ok := m["task"]
	if !ok {
		utils.ErrorResponse(ctx, utils.NewError(GetLogError, "not found task"))
		return
	}

	path, ok := m["path"]
	if !ok {
		utils.ErrorResponse(ctx, utils.NewError(GetLogError, "not found path"))
		return
	}

	user, ok := m["user"]
	if !ok {
		utils.ErrorResponse(ctx, utils.NewError(GetLogError, "not found user"))
		return
	}

	cluster, ok := m["cluster"]
	if !ok {
		utils.ErrorResponse(ctx, utils.NewError(GetLogError, "not found cluster"))
		return
	}

	message, ok := m["message"]
	if !ok {
		utils.ErrorResponse(ctx, utils.NewError(GetLogError, "not found message"))
		return
	}

	keywords, ok := s.KeywordFilter[app.(string)+path.(string)]
	if !ok {
		utils.Ok(ctx, "ok")
		return
	}
	for e := keywords.Front(); e != nil; e = e.Next() {
		if strings.Index(message.(string), e.Value.(string)) == -1 {
			continue
		}

		s.Counter.WithLabelValues(
			app.(string),
			task.(string),
			path.(string),
			e.Value.(string),
			user.(string),
			cluster.(string),
		).Inc()
	}

	utils.Ok(ctx, "ok")
	return
}