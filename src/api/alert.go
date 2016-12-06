package api

import (
	"errors"
	"time"

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

	s.Kmutex.Lock()
	defer s.Kmutex.Unlock()
	s.KeywordFilter[alert.AppId+alert.Path] = append(s.KeywordFilter[alert.AppId+alert.Path], alert.Keyword)

	utils.Ok(ctx, "create success")
}

func (s *search) DeleteAlert(ctx *gin.Context) {
	if ctx.Param("id") == "" {
		utils.ErrorResponse(ctx, utils.NewError(PARAM_ERROR, errors.New("alert id can't be empty")))
		return
	}

	alert, err := s.Service.GetAlert(ctx.Param("id"))
	if err != nil {
		utils.ErrorResponse(ctx, utils.NewError(GET_ALERT_ERROR, err))
		return
	}

	err = s.Service.DeleteAlert(ctx.Param("id"))
	if err != nil {
		utils.ErrorResponse(ctx, utils.NewError(DELETE_ALERT_ERROR, err))
		return
	}

	s.Kmutex.Lock()
	defer s.Kmutex.Unlock()
	for i, v := range s.KeywordFilter[alert.AppId+alert.Path] {
		if v == alert.Keyword {
			s.KeywordFilter[alert.AppId+alert.Path] = append(s.KeywordFilter[alert.AppId+alert.Path][:i],
				s.KeywordFilter[alert.AppId+alert.Path][i+1:]...)
			break
		}
	}

	utils.Ok(ctx, "delete success")
}

func (s *search) GetAlerts(ctx *gin.Context) {
	results, err := s.Service.GetAlerts(ctx.MustGet("page").(models.Page))
	if err != nil {
		utils.ErrorResponse(ctx, utils.NewError(GET_ALERT_ERROR, err))
		return
	}

	utils.Ok(ctx, results)
}

func (s *search) GetAlert(ctx *gin.Context) {
	result, err := s.Service.GetAlert(ctx.Param("id"))
	if err != nil {
		utils.ErrorResponse(ctx, utils.NewError(GET_ALERT_ERROR, err))
		return
	}

	utils.Ok(ctx, result)
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

	if alert.AppId == "" {
		utils.ErrorResponse(ctx, utils.NewError(PARAM_ERROR, errors.New("appid can't be empty")))
		return
	}

	if alert.Keyword == "" {
		utils.ErrorResponse(ctx, utils.NewError(PARAM_ERROR, errors.New("keyword can't be empty")))
		return
	}

	result, err := s.Service.GetAlert(ctx.Param("id"))
	if err != nil {
		utils.ErrorResponse(ctx, utils.NewError(GET_ALERT_ERROR, err))
		return
	}

	alert.CreateTime = time.Now().Format(time.RFC3339Nano)
	err = s.Service.UpdateAlert(alert)
	if err != nil {
		utils.ErrorResponse(ctx, utils.NewError(UPDATE_ALERT_ERROR, err))
		return
	}

	s.Kmutex.Lock()
	defer s.Kmutex.Unlock()
	for i, v := range s.KeywordFilter[result.AppId+result.Path] {
		if v == alert.Keyword {
			s.KeywordFilter[result.AppId+result.Path] = append(s.KeywordFilter[result.AppId+result.Path][:i],
				s.KeywordFilter[result.AppId+result.Path][i+1:]...)
			break
		}
	}
	s.KeywordFilter[result.AppId+result.Path] = append(s.KeywordFilter[result.AppId+result.Path], alert.Keyword)

	utils.Ok(ctx, "update success")
}
