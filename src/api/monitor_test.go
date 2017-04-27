package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Dataman-Cloud/log-proxy/src/utils/prometheusexpr"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func getQueryItems(ctx *gin.Context) {
	mo.GetQueryItems(ctx)
}

func TestQueryItems(t *testing.T) {
	path := "../../config/exprs/"
	prometheusexpr.Exprs(path)

	router := gin.New()
	router.GET("/api/v1/monitor/query/items", getQueryItems)
	testServer := httptest.NewServer(router)
	assert.NotNil(t, testServer)
	defer testServer.Close()
	_, err := http.Get(testServer.URL + "/api/v1/monitor/query/items")
	assert.NoError(t, err)
}
