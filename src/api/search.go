package api

import (
	"errors"
	"strconv"

	"github.com/Dataman-Cloud/log-proxy/src/service"
	"github.com/Dataman-Cloud/log-proxy/src/utils"

	"github.com/gin-gonic/gin"
)

const (
	GETAPPS_ERROR      = "503-11000"
	PARAM_ERROR        = "400-11001"
	GETTASK_ERROR      = "503-11002"
	INDEX_ERROR        = "503-11003"
	CREATE_ALERT_ERROR = "503-11004"
	DELETE_ALERT_ERROR = "503-11005"
	GET_ALERT_ERROR    = "503-11006"
	UPDATE_ALERT_ERROR = "503-11007"
)

type search struct {
	Service *service.SearchService
}

func GetSearch() *search {
	return &search{
		Service: service.NewSearchService(),
	}
}

func (s *search) Ping(ctx *gin.Context) {
	utils.Ok(ctx, "success")
}

func (s *search) Applications(ctx *gin.Context) {
	apps, err := s.Service.Applications()
	if err != nil {
		utils.ErrorResponse(ctx, utils.NewError(GETAPPS_ERROR, err))
		return
	}
	utils.Ok(ctx, apps)
}

func (s *search) Tasks(ctx *gin.Context) {
	appName := ctx.Param("appid")
	if appName == "" {
		utils.ErrorResponse(ctx, utils.NewError(PARAM_ERROR, errors.New("param error")))
		return
	}
	tasks, err := s.Service.Tasks(appName)
	if err != nil {
		utils.ErrorResponse(ctx, utils.NewError(GETTASK_ERROR, err))
		return
	}
	utils.Ok(ctx, tasks)
}

func (s *search) Paths(ctx *gin.Context) {
	appName := ctx.Param("appid")
	if appName == "" {
		utils.ErrorResponse(ctx, utils.NewError(PARAM_ERROR, errors.New("param error")))
		return
	}

	taskId := ctx.Param("taskid")
	if taskId == "" {
		utils.ErrorResponse(ctx, utils.NewError(PARAM_ERROR, errors.New("param error")))
		return
	}

	paths, err := s.Service.Paths(appName, taskId)
	if err != nil {
		utils.ErrorResponse(ctx, utils.NewError(GETTASK_ERROR, err))
		return
	}
	utils.Ok(ctx, paths)
}

func (s *search) Index(ctx *gin.Context) {
	if ctx.Query("appid") == "" {
		utils.ErrorResponse(ctx, utils.NewError(PARAM_ERROR, errors.New("appid can't be empty")))
		return
	}
	results, err := s.Service.Search(ctx.Query("appid"),
		ctx.Query("taskid"),
		ctx.Query("path"),
		ctx.Query("keyword"))
	if err != nil {
		utils.ErrorResponse(ctx, utils.NewError(INDEX_ERROR, err))
		return
	}

	utils.Ok(ctx, results)
}

func (s *search) Context(ctx *gin.Context) {
	if ctx.Query("appid") == "" {
		utils.ErrorResponse(ctx, utils.NewError(PARAM_ERROR, errors.New("appid can't be empty")))
		return
	}

	if ctx.Query("taskid") == "" {
		utils.ErrorResponse(ctx, utils.NewError(PARAM_ERROR, errors.New("taskid can't be empty")))
		return
	}

	if ctx.Query("path") == "" {
		utils.ErrorResponse(ctx, utils.NewError(PARAM_ERROR, errors.New("path can't be empty")))
		return
	}

	if ctx.Query("offset") == "" {
		utils.ErrorResponse(ctx, utils.NewError(PARAM_ERROR, errors.New("offset can't be empty")))
		return
	}

	results, err := s.Service.Context(ctx.Query("appid"),
		ctx.Query("taskid"),
		ctx.Query("path"),
		ctx.Query("offset"))
	if err != nil {
		utils.ErrorResponse(ctx, utils.NewError(INDEX_ERROR, err))
		return
	}

	utils.Ok(ctx, results)
}

func (s *search) Middleware(ctx *gin.Context) {

	if ctx.Query("from") != "" {
		if from, err := strconv.ParseInt(ctx.Query("from"), 10, 64); err == nil {
			s.Service.RangeFrom = from
		} else {
			s.Service.RangeFrom = ctx.Query("from")
		}
	} else {
		s.Service.RangeFrom = nil
	}

	if ctx.Query("to") != "" {
		if to, err := strconv.ParseInt(ctx.Query("to"), 10, 64); err == nil {
			s.Service.RangeTo = to
		} else {
			s.Service.RangeTo = ctx.Query("to")
		}
	} else {
		s.Service.RangeTo = nil
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
