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
	if err := json.Unmarshal(data, &m); err != nil {
		utils.ErrorResponse(ctx, utils.NewError(GetLogError, err))
		return
	}

	app, ok := m["DM_GROUP_NAME"]
	if !ok {
		utils.ErrorResponse(ctx, utils.NewError(GetLogError, "log group not found"))
		return
	}

	path, ok := m["DM_APP_ID"]
	if !ok {
		utils.ErrorResponse(ctx, utils.NewError(GetLogError, "log appid not found"))
		return
	}

	message, ok := m["message"]
	if !ok {
		utils.ErrorResponse(ctx, utils.NewError(GetLogError, "log message not found"))
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
	}

	utils.Ok(ctx, "ok")
	return
}
