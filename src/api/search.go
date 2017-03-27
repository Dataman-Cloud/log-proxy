package api

import (
	"container/list"
	"errors"
	"strings"
	"sync"

	"github.com/Dataman-Cloud/log-proxy/src/config"
	"github.com/Dataman-Cloud/log-proxy/src/models"
	"github.com/Dataman-Cloud/log-proxy/src/service"
	"github.com/Dataman-Cloud/log-proxy/src/store"
	"github.com/Dataman-Cloud/log-proxy/src/store/datastore"
	"github.com/Dataman-Cloud/log-proxy/src/utils"
	"github.com/Dataman-Cloud/log-proxy/src/utils/database"
	"github.com/Sirupsen/logrus"

	"github.com/gin-gonic/gin"
)

const (
	// GetAppsError get application error
	GetAppsError = "503-11000"
	// ParamError param error
	ParamError = "400-11001"
	// GetTaskError get task error
	GetTaskError = "503-11002"
	// IndexError search index log error
	IndexError = "503-11003"
	// CreateAlertError create keyword error
	CreateAlertError = "503-11004"
	// DeleteAlertError delete keyword error
	DeleteAlertError = "503-11005"
	// GetAlertError get keyword error
	GetAlertError = "503-11006"
	// UpdateAlertError update keyword error
	UpdateAlertError = "503-11007"
	// GetEventsError get event history error
	GetEventsError = "503-11008"
	// GetPrometheusError get prometheus event error
	GetPrometheusError = "503-11009"
	// GetLogError get log error
	GetLogError = "503-11010"

	GetClustersError = "503-11011"

	GetLogAlertEventsError = "503-11012"

	DeleteLogAlertEventsError = "503-11013"
)

// Search search client struct
type Search struct {
	Store         store.Store
	Service       *service.SearchService
	KeywordFilter map[string]*list.List
	Kmutex        *sync.Mutex
}

// GetSearch new search client
func GetSearch() *Search {
	s := &Search{
		Store:         datastore.From(database.GetDB()),
		Service:       service.NewEsService(strings.Split(config.GetConfig().EsURL, ",")),
		KeywordFilter: make(map[string]*list.List),
		Kmutex:        new(sync.Mutex),
	}

	//TODO: rules count maybe greate than 10000
	alerts, err := s.Store.GetLogAlertRules(models.Page{
		PageFrom: 0,
		PageSize: 10000,
	})

	if err != nil {
		logrus.Errorf("init rules of log alert failed, error: %s", err.Error())
		return s
	}

	s.Kmutex.Lock()
	defer s.Kmutex.Unlock()
	if alerts == nil {
		return s
	}

	for _, alertRule := range alerts["rules"].([]*models.LogAlertRule) {
		ruleIndex := alertRule.App + alertRule.Source
		if s.KeywordFilter[ruleIndex] == nil {
			s.KeywordFilter[ruleIndex] = list.New()
		}
		s.KeywordFilter[ruleIndex].PushBack(alertRule)
	}

	return s
}

// Ping ping
func (s *Search) Ping(ctx *gin.Context) {
	utils.Ok(ctx, "success")
}

func (s *Search) Clusters(ctx *gin.Context) {
	clusters, err := s.Service.Clusters(ctx.MustGet("page").(models.Page))
	if err != nil {
		utils.ErrorResponse(ctx, utils.NewError(GetClustersError, err))
		return
	}

	utils.Ok(ctx, clusters)
	return
}

// Applications get all applications
func (s *Search) Applications(ctx *gin.Context) {
	apps, err := s.Service.Applications(ctx.Param("cluster"), ctx.MustGet("page").(models.Page))
	if err != nil {
		utils.ErrorResponse(ctx, utils.NewError(GetAppsError, err))
		return
	}
	utils.Ok(ctx, apps)
}

// Tasks search applications tasks
func (s *Search) Tasks(ctx *gin.Context) {
	appName := ctx.Param("app")
	tasks, err := s.Service.Tasks(ctx.Param("cluster"), appName, ctx.MustGet("page").(models.Page))
	if err != nil {
		utils.ErrorResponse(ctx, utils.NewError(GetTaskError, err))
		return
	}

	borgToken, err := utils.LoginBorg(config.GetConfig().BorgUser, config.GetConfig().BorgPassword, config.BorgLoginURL())
	if err != nil {
		utils.ErrorResponse(ctx, utils.NewError(GetTaskError, err))
		return
	}

	runningTasks, err := utils.ListAppTaskFromBorg(borgToken, config.BorgAppTasksURL(appName))
	if err != nil {
		utils.ErrorResponse(ctx, utils.NewError(GetTaskError, err))
		return
	}

	var taskInfoList []models.TaskInfo
	for _, runningTask := range runningTasks {
		if count, ok := tasks[runningTask.ID]; ok {
			taskInfo := models.TaskInfo{runningTask.ID, models.TaskRunning, count}
			taskInfoList = append(taskInfoList, taskInfo)
			delete(tasks, runningTask.ID)
		} else {
			continue
		}
	}

	for taskID, count := range tasks {
		taskInfo := models.TaskInfo{taskID, models.TaskDied, count}
		taskInfoList = append(taskInfoList, taskInfo)
	}

	utils.Ok(ctx, taskInfoList)
}

// Paths search applications paths
func (s *Search) Source(ctx *gin.Context) {
	paths, err := s.Service.Source(
		ctx.Param("cluster"),
		ctx.Param("app"),
		ctx.Query("task"),
		ctx.MustGet("page").(models.Page),
	)
	if err != nil {
		utils.ErrorResponse(ctx, utils.NewError(GetTaskError, err))
		return
	}
	utils.Ok(ctx, paths)
}

// Index search log by condition
func (s *Search) Index(ctx *gin.Context) {
	if ctx.Param("app") == "" {
		utils.ErrorResponse(ctx, utils.NewError(ParamError, errors.New("app can't be empty")))
		return
	}

	results, err := s.Service.Search(
		ctx.Param("cluster"),
		ctx.Param("app"),
		ctx.Query("task"),
		ctx.Query("path"),
		ctx.Query("keyword"),
		ctx.MustGet("page").(models.Page))
	if err != nil {
		utils.ErrorResponse(ctx, utils.NewError(IndexError, err))
		return
	}

	utils.Ok(ctx, results)
}

// Context search log context
func (s *Search) Context(ctx *gin.Context) {
	if ctx.Param("app") == "" {
		utils.ErrorResponse(ctx, utils.NewError(ParamError, errors.New("app can't be empty")))
		return
	}

	if ctx.Query("task") == "" {
		utils.ErrorResponse(ctx, utils.NewError(ParamError, errors.New("task can't be empty")))
		return
	}

	if ctx.Query("source") == "" {
		utils.ErrorResponse(ctx, utils.NewError(ParamError, errors.New("source can't be empty")))
		return
	}

	if ctx.Query("offset") == "" {
		utils.ErrorResponse(ctx, utils.NewError(ParamError, errors.New("offset can't be empty")))
		return
	}

	results, err := s.Service.Context(
		ctx.Param("cluster"),
		ctx.Param("app"),
		ctx.Query("task"),
		ctx.Query("source"),
		ctx.Query("offset"),
		ctx.MustGet("page").(models.Page))
	if err != nil {
		utils.ErrorResponse(ctx, utils.NewError(IndexError, err))
		return
	}

	utils.Ok(ctx, results)
}
