package api

import (
	"errors"
	"strconv"

	"github.com/Dataman-Cloud/log-proxy/src/models"
	"github.com/Dataman-Cloud/log-proxy/src/utils"

	"github.com/gin-gonic/gin"
)

func (s *search) CreateAlert(ctx *gin.Context) {
	alert := new(models.Alert)
	if err := ctx.BindJSON(alert); err != nil {
		utils.ErrorResponse(ctx, utils.NewError(PARAM_ERROR, errors.New("request body param error")))
		return
	}

	if alert.Period == "" {
		utils.ErrorResponse(ctx, utils.NewError(PARAM_ERROR, errors.New("period can't be empty")))
		return
	}

	if alert.AppId == "" {
		utils.ErrorResponse(ctx, utils.NewError(PARAM_ERROR, errors.New("appid can't be empty")))
		return
	}

	if alert.Keyword == "" {
		utils.ErrorResponse(ctx, utils.NewError(PARAM_ERROR, errors.New("keyword can't be empty")))
		return
	}

	err := s.Service.CreateAlert(alert)
	if err != nil {
		utils.ErrorResponse(ctx, utils.NewError(CREATE_ALERT_ERROR, err))
		return
	}

	utils.Ok(ctx, "create success")
}

func (s *search) DeleteAlert(ctx *gin.Context) {
	if ctx.Param("id") == "" {
		utils.ErrorResponse(ctx, utils.NewError(PARAM_ERROR, errors.New("alert id can't be empty")))
		return
	}

	err := s.Service.DeleteAlert(ctx.Param("id"))
	if err != nil {
		utils.ErrorResponse(ctx, utils.NewError(DELETE_ALERT_ERROR, err))
		return
	}

	utils.Ok(ctx, "delete success")
}

func (s *search) GetAlerts(ctx *gin.Context) {
	results, err := s.Service.GetAlerts()
	if err != nil {
		utils.ErrorResponse(ctx, utils.NewError(GET_ALERT_ERROR, err))
		return
	}

	utils.Ok(ctx, results)
}

func (s *search) UpdateAlert(ctx *gin.Context) {
	alert := new(models.Alert)
	if err := ctx.BindJSON(alert); err != nil {
		utils.ErrorResponse(ctx, utils.NewError(PARAM_ERROR, errors.New("request body param error")))
		return
	}

	if alert.Id == "" {
		utils.ErrorResponse(ctx, utils.NewError(PARAM_ERROR, errors.New("id can't be empty")))
		return
	}

	if alert.Period == "" {
		utils.ErrorResponse(ctx, utils.NewError(PARAM_ERROR, errors.New("period can't be empty")))
		return
	}

	if alert.AppId == "" {
		utils.ErrorResponse(ctx, utils.NewError(PARAM_ERROR, errors.New("appid can't be empty")))
		return
	}

	if alert.Keyword == "" {
		utils.ErrorResponse(ctx, utils.NewError(PARAM_ERROR, errors.New("keyword can't be empty")))
		return
	}

	err := s.Service.UpdateAlert(alert)
	if err != nil {
		utils.ErrorResponse(ctx, utils.NewError(UPDATE_ALERT_ERROR, err))
		return
	}

	utils.Ok(ctx, "update success")
}

func (s *search) Middleware(ctx *gin.Context) {
	if ctx.Query("from") != "" {
		s.Service.RangeFrom = ctx.Query("from")
	}

	if ctx.Query("to") != "" {
		s.Service.RangeTo = ctx.Query("to")
	}

	if size, err := strconv.Atoi(ctx.Query("size")); err == nil && size > 0 {
		s.Service.PageSize = size
	} else {
		s.Service.PageSize = 100
	}

	if page, err := strconv.Atoi(ctx.Query("page")); err == nil && page > 0 {
		s.Service.PageFrom = (page - 1) * s.Service.PageSize
	} else {
		s.Service.PageFrom = 0
	}

	ctx.Next()
}
