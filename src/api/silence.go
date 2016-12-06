package api

import (
	"net/http"

	"github.com/Dataman-Cloud/log-proxy/src/backends"
	"github.com/Dataman-Cloud/log-proxy/src/config"
	"github.com/Dataman-Cloud/log-proxy/src/utils"

	"github.com/gin-gonic/gin"
)

const (
	SILENCE_PARAM_ERROR = "503-12000"
)

func (m *monitor) GetSilences(ctx *gin.Context) {
	query := &backends.AlertManager{
		HttpClient: http.DefaultClient,
		Server:     config.GetConfig().ALERTMANAGER_URL,
		Path:       ALERTSPATH,
	}

	data, err := query.GetSilences()
	if err != nil {
		utils.ErrorResponse(ctx, err)
		return
	}
	utils.Ok(ctx, data)
}

func (m *monitor) CreateSilence(ctx *gin.Context) {
	query := &backends.AlertManager{
		HttpClient: http.DefaultClient,
		Server:     config.GetConfig().ALERTMANAGER_URL,
		Path:       ALERTSPATH,
	}

	var silence map[string]interface{}
	if err := ctx.BindJSON(&silence); err != nil {
		utils.ErrorResponse(ctx, err)
		return
	}

	err := query.CreateSilence(silence)
	if err != nil {
		utils.ErrorResponse(ctx, err)
		return
	}
	utils.Ok(ctx, "success")
}

func (m *monitor) GetSilence(ctx *gin.Context) {
	query := &backends.AlertManager{
		HttpClient: http.DefaultClient,
		Server:     config.GetConfig().ALERTMANAGER_URL,
		Path:       ALERTSPATH,
	}

	data, err := query.GetSilence(ctx.Param("id"))
	if err != nil {
		utils.ErrorResponse(ctx, err)
		return
	}
	utils.Ok(ctx, data)
}

func (m *monitor) DeleteSilence(ctx *gin.Context) {
	query := &backends.AlertManager{
		HttpClient: http.DefaultClient,
		Server:     config.GetConfig().ALERTMANAGER_URL,
		Path:       ALERTSPATH,
	}

	err := query.DeleteSilence(ctx.Param("id"))
	if err != nil {
		utils.ErrorResponse(ctx, err)
		return
	}
	utils.Ok(ctx, "success")
}
