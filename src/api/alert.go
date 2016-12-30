package api

import (
	"errors"
	"time"

	"github.com/Dataman-Cloud/log-proxy/src/models"
	"github.com/Dataman-Cloud/log-proxy/src/utils"

	"github.com/gin-gonic/gin"
)

// CreateAlert create keyword filter
func (s *Search) CreateAlert(ctx *gin.Context) {
	alert := new(models.Alert)
	if err := ctx.BindJSON(alert); err != nil {
		utils.ErrorResponse(ctx, utils.NewError(ParamError, errors.New("request body param error")))
		return
	}

	if alert.AppID == "" {
		utils.ErrorResponse(ctx, utils.NewError(ParamError, errors.New("appid can't be empty")))
		return
	}

	if alert.Keyword == "" {
		utils.ErrorResponse(ctx, utils.NewError(ParamError, errors.New("keyword can't be empty")))
		return
	}

	for _, v := range s.KeywordFilter[alert.AppID+alert.Path] {
		if v == alert.Keyword {
			utils.ErrorResponse(ctx, utils.NewError(CreateAlertError, errors.New("keyword exist")))
			return
		}
	}

	err := s.Service.CreateAlert(alert)
	if err != nil {
		utils.ErrorResponse(ctx, utils.NewError(CreateAlertError, err))
		return
	}

	s.Kmutex.Lock()
	defer s.Kmutex.Unlock()
	s.KeywordFilter[alert.AppID+alert.Path] = append(s.KeywordFilter[alert.AppID+alert.Path], alert.Keyword)

	utils.Ok(ctx, "create success")
}

// DeleteAlert delete keyword filter
func (s *Search) DeleteAlert(ctx *gin.Context) {

	alert, err := s.Service.GetAlert(ctx.Param("id"))
	if err != nil {
		utils.ErrorResponse(ctx, utils.NewError(GetAlertError, err))
		return
	}

	err = s.Service.DeleteAlert(ctx.Param("id"))
	if err != nil {
		utils.ErrorResponse(ctx, utils.NewError(DeleteAlertError, err))
		return
	}

	s.Kmutex.Lock()
	defer s.Kmutex.Unlock()
	for i, v := range s.KeywordFilter[alert.AppID+alert.Path] {
		if v == alert.Keyword {
			s.KeywordFilter[alert.AppID+alert.Path] = append(s.KeywordFilter[alert.AppID+alert.Path][:i],
				s.KeywordFilter[alert.AppID+alert.Path][i+1:]...)
			break
		}
	}

	utils.Ok(ctx, "delete success")
}

// GetAlerts get all keyword filter
func (s *Search) GetAlerts(ctx *gin.Context) {
	results, err := s.Service.GetAlerts(ctx.MustGet("page").(models.Page))
	if err != nil {
		utils.ErrorResponse(ctx, utils.NewError(GetAlertError, err))
		return
	}

	utils.Ok(ctx, results)
}

// GetAlert get keyword filter by id
func (s *Search) GetAlert(ctx *gin.Context) {
	result, err := s.Service.GetAlert(ctx.Param("id"))
	if err != nil {
		utils.ErrorResponse(ctx, utils.NewError(GetAlertError, err))
		return
	}

	utils.Ok(ctx, result)
}

// UpdateAlert update keyword filter
func (s *Search) UpdateAlert(ctx *gin.Context) {
	alert := new(models.Alert)
	if err := ctx.BindJSON(alert); err != nil {
		utils.ErrorResponse(ctx, utils.NewError(ParamError, errors.New("request body param error")))
		return
	}

	if alert.ID == "" {
		utils.ErrorResponse(ctx, utils.NewError(ParamError, errors.New("id can't be empty")))
		return
	}

	if alert.AppID == "" {
		utils.ErrorResponse(ctx, utils.NewError(ParamError, errors.New("appid can't be empty")))
		return
	}

	if alert.Keyword == "" {
		utils.ErrorResponse(ctx, utils.NewError(ParamError, errors.New("keyword can't be empty")))
		return
	}

	result, err := s.Service.GetAlert(alert.ID)
	if err != nil {
		utils.ErrorResponse(ctx, utils.NewError(GetAlertError, err))
		return
	}

	alert.CreateTime = time.Now().Format(time.RFC3339Nano)
	err = s.Service.UpdateAlert(alert)
	if err != nil {
		utils.ErrorResponse(ctx, utils.NewError(UpdateAlertError, err))
		return
	}

	s.Kmutex.Lock()
	defer s.Kmutex.Unlock()
	for i, v := range s.KeywordFilter[result.AppID+result.Path] {
		if v == result.Keyword {
			s.KeywordFilter[result.AppID+result.Path][i] = alert.Keyword
		}
	}

	utils.Ok(ctx, "update success")
}
