package api

import (
	"errors"
	"strconv"

	"github.com/Dataman-Cloud/log-proxy/src/service"
	"github.com/Dataman-Cloud/log-proxy/src/utils"

	"github.com/gin-gonic/gin"
)

const (
	GETAPPS_ERROR = "503-11000"
	PARAM_ERROR   = "400-11001"
	GETTASK_ERROR = "503-11002"
	INDEX_ERROR   = "503-11003"
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
	results, err := s.Service.Search(ctx.Param("appid"),
		ctx.Query("taskid"),
		ctx.Query("path"),
		ctx.Query("keyword"))
	if err != nil {
		utils.ErrorResponse(ctx, utils.NewError(INDEX_ERROR, err))
		return
	}

	utils.Ok(ctx, results)
}

func (s *search) Middleware(ctx *gin.Context) {
	if ctx.Query("from") != "" {
		s.Service.RangeFrom = ctx.Query("from")
	}

	if ctx.Query("to") != "" {
		s.Service.RangeTo = ctx.Query("to")
	}

	if size, err := strconv.Atoi(ctx.Param("size")); err == nil && size > 0 {
		s.Service.PageSize = size
	} else {
		s.Service.PageSize = 100
	}

	if page, err := strconv.Atoi(ctx.Param("page")); err == nil && page > 0 {
		s.Service.PageFrom = (page - 1) * s.Service.PageSize
	} else {
		s.Service.PageFrom = 0
	}

	ctx.Next()
}
