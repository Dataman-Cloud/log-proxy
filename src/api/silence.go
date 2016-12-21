package api

import (
	"net/http"

	"github.com/Dataman-Cloud/log-proxy/src/backends"
	"github.com/Dataman-Cloud/log-proxy/src/config"
	"github.com/Dataman-Cloud/log-proxy/src/utils"

	"github.com/gin-gonic/gin"
)

const (
	// SilenceParamError silence error response code
	SilenceParamError = "503-12000"
)

// GetSilences return the silences list
func (m *Monitor) GetSilences(ctx *gin.Context) {
	query := &backends.AlertManager{
		HTTPClient: http.DefaultClient,
		Server:     config.GetConfig().AlertManagerURL,
		Path:       backends.AlertsPath,
	}

	data, err := query.GetSilences()
	if err != nil {
		utils.ErrorResponse(ctx, err)
		return
	}
	utils.Ok(ctx, data)
}

// CreateSilence pass the silence varabile to the func query.CreateSilence
func (m *Monitor) CreateSilence(ctx *gin.Context) {
	query := &backends.AlertManager{
		HTTPClient: http.DefaultClient,
		Server:     config.GetConfig().AlertManagerURL,
		Path:       backends.AlertsPath,
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

// GetSilence return the silence
func (m *Monitor) GetSilence(ctx *gin.Context) {
	query := &backends.AlertManager{
		HTTPClient: http.DefaultClient,
		Server:     config.GetConfig().AlertManagerURL,
		Path:       backends.AlertsPath,
	}

	data, err := query.GetSilence(ctx.Param("id"))
	if err != nil {
		utils.ErrorResponse(ctx, err)
		return
	}
	utils.Ok(ctx, data)
}

// DeleteSilence pass the id to the func query.DeleteSilence
func (m *Monitor) DeleteSilence(ctx *gin.Context) {
	query := &backends.AlertManager{
		HTTPClient: http.DefaultClient,
		Server:     config.GetConfig().AlertManagerURL,
		Path:       backends.AlertsPath,
	}

	err := query.DeleteSilence(ctx.Param("id"))
	if err != nil {
		utils.ErrorResponse(ctx, err)
		return
	}
	utils.Ok(ctx, "success")
}

// UpdateSilence delete and create the silence by id
func (m *Monitor) UpdateSilence(ctx *gin.Context) {
	query := &backends.AlertManager{
		HTTPClient: http.DefaultClient,
		Server:     config.GetConfig().AlertManagerURL,
		Path:       backends.AlertsPath,
	}

	var silence map[string]interface{}
	if err := ctx.BindJSON(&silence); err != nil {
		utils.ErrorResponse(ctx, err)
		return
	}

	err := query.DeleteSilence(ctx.Param("id"))
	if err != nil {
		utils.ErrorResponse(ctx, err)
		return
	}

	err = query.CreateSilence(silence)
	if err != nil {
		utils.ErrorResponse(ctx, err)
		return
	}
	utils.Ok(ctx, "success")
}
