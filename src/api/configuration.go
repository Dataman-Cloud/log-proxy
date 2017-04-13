package api

import (
	"github.com/Dataman-Cloud/log-proxy/src/models"
	"github.com/Dataman-Cloud/log-proxy/src/store"
	"github.com/Dataman-Cloud/log-proxy/src/store/datastore"
	"github.com/Dataman-Cloud/log-proxy/src/utils"
	"github.com/Dataman-Cloud/log-proxy/src/utils/database"
	"github.com/gin-gonic/gin"
)

type Conf struct {
	Store store.Store
}

func NewConf() *Conf {
	return &Conf{
		Store: datastore.From(database.GetDB()),
	}
}

func (c *Conf) UpdateConf(ctx *gin.Context) {
	var (
		data models.Configuration
		err  error
	)
	conf := models.NewConfiguration()
	if err = ctx.BindJSON(&conf); err != nil {
		utils.ErrorResponse(ctx, err)
		return
	}
	err = c.Store.UpdateConf(conf)
	if err != nil {
		utils.ErrorResponse(ctx, err)
		return
	}
	data, err = c.Store.GetConf()
	if err != nil {
		utils.ErrorResponse(ctx, err)
		return
	}
	utils.Ok(ctx, data)
}

func (c *Conf) GetConf(ctx *gin.Context) {
	var (
		data models.Configuration
		err  error
	)
	data, err = c.Store.GetConf()
	if err != nil {
		utils.ErrorResponse(ctx, err)
		return
	}
	utils.Ok(ctx, data)
}
