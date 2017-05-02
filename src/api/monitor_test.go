package api

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/Dataman-Cloud/log-proxy/src/utils/prometheusexpr"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func getQueryItems(ctx *gin.Context) {
	mo.GetQueryItems(ctx)
}

func getQuery(ctx *gin.Context) {
	mo.Query(ctx)
}

func getQueryApps(ctx *gin.Context) {
	mo.GetApps(ctx)
}

func getQueryAppTasks(ctx *gin.Context) {
	mo.GetAppsTasks(ctx)
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

func TestQueryMetric(t *testing.T) {
	path := "../../config/exprs/"
	prometheusexpr.Exprs(path)

	router := gin.New()
	router.GET("/api/v1/monitor/query", getQuery)
	testServer := httptest.NewServer(router)
	assert.NotNil(t, testServer)
	defer testServer.Close()

	httpClient := http.DefaultClient
	u, _ := url.Parse(testServer.URL)
	u.Path = strings.TrimRight(u.Path, "/") + "/api/v1/monitor/query"
	q := u.Query()
	q.Set("metric", "内存使用字节数")
	q.Set("app", "web-zdou-datamanmesos")
	q.Set("task", "0")
	q.Set("start", "1481853425")
	q.Set("end", "1481853425")
	q.Set("step", "30s")
	u.RawQuery = q.Encode()
	resp, _ := httpClient.Get(u.String())
	if resp.Request.FormValue("metric") != "内存使用字节数" {
		t.Errorf("Expect query param metric is 内存使用字节数, got %s", resp.Request.FormValue("metric"))
	}
}

func TestQueryExpr(t *testing.T) {
	path := "../../config/exprs/"
	prometheusexpr.Exprs(path)

	router := gin.New()
	router.GET("/api/v1/monitor/query", getQuery)
	testServer := httptest.NewServer(router)
	assert.NotNil(t, testServer)
	defer testServer.Close()

	httpClient := http.DefaultClient
	u, _ := url.Parse(testServer.URL)
	u.Path = strings.TrimRight(u.Path, "/") + "/api/v1/monitor/query"
	q := u.Query()
	q.Set("expr", "up")
	q.Set("start", "1481853425")
	q.Set("end", "1481853425")
	q.Set("step", "30s")
	u.RawQuery = q.Encode()
	resp, _ := httpClient.Get(u.String())
	if resp.Request.FormValue("expr") != "up" {
		t.Errorf("Expect query param expr is up, got %s", resp.Request.FormValue("metric"))
	}
}

func TestQueryError(t *testing.T) {
	path := "../../config/exprs/"
	prometheusexpr.Exprs(path)

	router := gin.New()
	router.GET("/api/v1/monitor/query", getQuery)
	testServer := httptest.NewServer(router)
	assert.NotNil(t, testServer)
	defer testServer.Close()

	httpClient := http.DefaultClient
	u, _ := url.Parse(testServer.URL)
	u.Path = strings.TrimRight(u.Path, "/") + "/api/v1/monitor/query"
	q := u.Query()
	u.RawQuery = q.Encode()
	resp, _ := httpClient.Get(u.String())
	if resp.StatusCode != 503 {
		t.Errorf("Expect StatusCode is 503, got %d", resp.StatusCode)
	}

	q = u.Query()
	q.Set("expr", "up")
	q.Set("metric", "内存使用字节数")
	u.RawQuery = q.Encode()
	resp, _ = httpClient.Get(u.String())
	if resp.StatusCode != 503 {
		t.Errorf("Expect StatusCode is 503, got %d", resp.StatusCode)
	}
}

func TestQueryErrorAppParam(t *testing.T) {
	path := "../../config/exprs/"
	prometheusexpr.Exprs(path)

	router := gin.New()
	router.GET("/api/v1/monitor/query", getQuery)
	testServer := httptest.NewServer(router)
	assert.NotNil(t, testServer)
	defer testServer.Close()

	httpClient := http.DefaultClient
	u, _ := url.Parse(testServer.URL)
	u.Path = strings.TrimRight(u.Path, "/") + "/api/v1/monitor/query"
	q := u.Query()
	u.RawQuery = q.Encode()
	resp, _ := httpClient.Get(u.String())
	if resp.StatusCode != 503 {
		t.Errorf("Expect StatusCode is 503, got %d", resp.StatusCode)
	}

	q = u.Query()
	q.Set("metric", "内存使用字节数")
	u.RawQuery = q.Encode()
	resp, _ = httpClient.Get(u.String())
	if resp.StatusCode != 503 {
		t.Errorf("Expect StatusCode is 503, got %d", resp.StatusCode)
	}
}

func TestQueryGetApps(t *testing.T) {
	path := "../../config/exprs/"
	prometheusexpr.Exprs(path)

	router := gin.New()
	router.GET("/api/v1/monitor/query/apps", getQueryApps)
	testServer := httptest.NewServer(router)
	assert.NotNil(t, testServer)
	defer testServer.Close()

	httpClient := http.DefaultClient
	u, _ := url.Parse(testServer.URL)
	u.Path = strings.TrimRight(u.Path, "/") + "/api/v1/monitor/query/apps"
	q := u.Query()
	q.Set("start", "1481853425")
	q.Set("end", "1481853425")
	q.Set("step", "30s")
	u.RawQuery = q.Encode()
	resp, _ := httpClient.Get(u.String())
	if resp.StatusCode != 503 {
		t.Errorf("Expect StatusCode is 503, got %d", resp.StatusCode)
	}
}

func TestQueryGetAppTasks(t *testing.T) {
	path := "../../config/exprs/"
	prometheusexpr.Exprs(path)

	router := gin.New()
	router.GET("/api/v1/monitor/query/apps/:appid/tasks", getQueryAppTasks)
	testServer := httptest.NewServer(router)
	assert.NotNil(t, testServer)
	defer testServer.Close()

	httpClient := http.DefaultClient
	u, _ := url.Parse(testServer.URL)
	u.Path = strings.TrimRight(u.Path, "/") + "/api/v1/monitor/query/apps/web-zdou-datamanmesos/tasks"
	q := u.Query()
	q.Set("start", "1481853425")
	q.Set("end", "1481853425")
	q.Set("step", "30s")
	u.RawQuery = q.Encode()
	resp, _ := httpClient.Get(u.String())
	if resp.StatusCode != 503 {
		t.Errorf("Expect StatusCode is 503, got %d", resp.StatusCode)
	}
}
