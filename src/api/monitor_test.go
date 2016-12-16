package api

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func getQuery(ctx *gin.Context) {
	mo.Query(ctx)
}

func getQueryInfo(ctx *gin.Context) {
	mo.QueryInfo(ctx)
}

func getQueryNodes(ctx *gin.Context) {
	mo.QueryNodes(ctx)
}

func getMonitorAlerts(ctx *gin.Context) {
	mo.GetAlerts(ctx)
}

func getMonitorAlertsGroups(ctx *gin.Context) {
	mo.GetAlertsGroups(ctx)
}

func getMonitorAlertsStatus(ctx *gin.Context) {
	mo.GetAlertsStatus(ctx)
}

func TestQueryMetric(t *testing.T) {
	httpClient := http.DefaultClient
	u, _ := url.Parse(apiUrl)
	u.Path = strings.TrimRight(u.Path, "/") + "/api/v1/monitor/query"
	q := u.Query()
	q.Set("metric", "memory")
	q.Set("appid", "work-web")
	q.Set("taskid", "4f84")
	q.Set("start", "1481853425")
	q.Set("end", "1481853425")
	q.Set("step", "30s")
	u.RawQuery = q.Encode()
	resp, _ := httpClient.Get(u.String())
	if resp.Request.FormValue("metric") != "memory" {
		t.Errorf("Expect query param metric is memory, got %s", resp.Request.FormValue("metric"))
	}
}

func TestQueryExpr(t *testing.T) {
	expr := "sum(container_memory_usage_bytes{id='/'}) by (instance)"
	httpClient := http.DefaultClient
	u, _ := url.Parse(apiUrl)
	u.Path = strings.TrimRight(u.Path, "/") + "/api/v1/monitor/query"
	q := u.Query()
	q.Set("expr", expr)
	q.Set("start", "1481853425")
	q.Set("end", "1481853425")
	q.Set("step", "30s")
	u.RawQuery = q.Encode()
	resp, _ := httpClient.Get(u.String())
	if resp.Request.FormValue("expr") != expr {
		t.Errorf("Expect query param metric is %s, got %s", expr, resp.Request.FormValue("metric"))
	}
}

func TestQueryParamConflict(t *testing.T) {
	metric := "memory"
	expr := "sum(container_memory_usage_bytes{id='/'}) by (instance)"
	httpClient := http.DefaultClient
	u, _ := url.Parse(apiUrl)
	u.Path = strings.TrimRight(u.Path, "/") + "/api/v1/monitor/query"
	q := u.Query()
	q.Set("metric", metric)
	q.Set("expr", expr)
	q.Set("start", "1481853425")
	q.Set("end", "1481853425")
	q.Set("step", "30s")
	u.RawQuery = q.Encode()
	resp, _ := httpClient.Get(u.String())
	if resp.Request.FormValue("expr") != expr {
		t.Errorf("Expect query param metric is %s, got %s", expr, resp.Request.FormValue("metric"))
	}
	if resp.Request.FormValue("metric") != metric {
		t.Errorf("Expect query param metric is %s, got %s", metric, resp.Request.FormValue("memory"))
	}
}

func TestQueryParamMissing(t *testing.T) {
	expectResult := 503
	httpClient := http.DefaultClient
	u, _ := url.Parse(apiUrl)
	u.Path = strings.TrimRight(u.Path, "/") + "/api/v1/monitor/query"
	q := u.Query()
	q.Set("start", "1481853425")
	q.Set("end", "1481853425")
	q.Set("step", "30s")
	u.RawQuery = q.Encode()
	resp, _ := httpClient.Get(u.String())
	if resp.StatusCode != expectResult {
		t.Errorf("Expect query param metric is %v, got %v", expectResult, resp.StatusCode)
	}
}
func TestQueryInfo(t *testing.T) {
	httpClient := http.DefaultClient
	u, _ := url.Parse(apiUrl)
	u.Path = strings.TrimRight(u.Path, "/") + "/api/v1/monitor/info"
	q := u.Query()
	q.Set("clusterid", "work")
	u.RawQuery = q.Encode()
	resp, _ := httpClient.Get(u.String())
	if resp.Request.FormValue("clusterid") != "work" {
		t.Errorf("Expect query param metric is work, got %s", resp.Request.FormValue("clusterid"))
	}
}

func TestQueryInfoConflict(t *testing.T) {
	expectResult := 503
	httpClient := http.DefaultClient
	u, _ := url.Parse(apiUrl)
	u.Path = strings.TrimRight(u.Path, "/") + "/api/v1/monitor/info"
	q := u.Query()
	q.Set("clusterid", "work")
	q.Set("appid", "work-web")
	u.RawQuery = q.Encode()
	resp, _ := httpClient.Get(u.String())
	if resp.StatusCode != expectResult {
		t.Errorf("Expect query param metric is %v, got %v", expectResult, resp.StatusCode)
	}
}

func TestQueryNodes(t *testing.T) {
	httpClient := http.DefaultClient
	u, _ := url.Parse(apiUrl)
	u.Path = strings.TrimRight(u.Path, "/") + "/api/v1/monitor/nodes"
	q := u.Query()
	q.Set("clusterid", "work")
	u.RawQuery = q.Encode()
	resp, _ := httpClient.Get(u.String())
	if resp.Request.FormValue("clusterid") != "work" {
		t.Errorf("Expect query param metric is work, got %s", resp.Request.FormValue("clusterid"))
	}
}

func TestMonitorGetAlerts(t *testing.T) {
	expectResult := 200
	httpClient := http.DefaultClient
	u, _ := url.Parse(apiUrl)
	u.Path = strings.TrimRight(u.Path, "/") + "/api/v1/monitor/alerts"
	resp, _ := httpClient.Get(u.String())
	if resp.StatusCode != expectResult {
		t.Errorf("Expect query param metric is %v, got %v", expectResult, resp.StatusCode)
	}
}

func TestMonitorGetAlertsGroups(t *testing.T) {
	expectResult := 200
	httpClient := http.DefaultClient
	u, _ := url.Parse(apiUrl)
	u.Path = strings.TrimRight(u.Path, "/") + "/api/v1/monitor/alerts/groups"
	resp, _ := httpClient.Get(u.String())
	if resp.StatusCode != expectResult {
		t.Errorf("Expect query param metric is %v, got %v", expectResult, resp.StatusCode)
	}
}

func TestMonitorGetAlertsStatus(t *testing.T) {
	expectResult := 503
	httpClient := http.DefaultClient
	u, _ := url.Parse(apiUrl)
	u.Path = strings.TrimRight(u.Path, "/") + "/api/v1/monitor/alerts/status"
	resp, _ := httpClient.Get(u.String())
	if resp.StatusCode != expectResult {
		t.Errorf("Expect query param metric is %v, got %v", expectResult, resp.StatusCode)
	}

	fmt.Printf("%v\n", resp)
}
