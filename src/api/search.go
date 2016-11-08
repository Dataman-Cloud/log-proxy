package api

import (
	"github.com/Dataman-Cloud/log-proxy/src/service"
	"github.com/Dataman-Cloud/log-proxy/src/utils"

	"github.com/gin-gonic/gin"
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
	s.Service.Applications()
	utils.Ok(ctx, "applications")
}
