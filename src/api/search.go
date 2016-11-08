package api

import (
	"errors"

	"github.com/Dataman-Cloud/log-proxy/src/service"
	"github.com/Dataman-Cloud/log-proxy/src/utils"

	//log "github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
)

const (
	GETAPPS_ERROR = "503-11000"
	PARAM_ERROR   = "400-11001"
	GETTASK_ERROR = "503-11002"
)

type search struct {
	Service *service.SearchService
}

func GetSearch() *search {
	return &search{
		Service: service.GetSearchService(),
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
	appName := ctx.Param("appname")
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
	appName := ctx.Param("appname")
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
