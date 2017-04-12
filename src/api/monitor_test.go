package api

import (
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/Dataman-Cloud/log-proxy/src/utils/prometheusexpr"
	"github.com/gin-gonic/gin"
)

func getQueryItems(ctx *gin.Context) {
	mo.GetQueryItems(ctx)
}

func getUsers(ctx *gin.Context) {
	mo.GetUsers(ctx)
}

func TestQueryItems(t *testing.T) {
	path := "../../config/exprs/"
	prometheusexpr.Exprs(path)

	httpClient := http.DefaultClient
	u, _ := url.Parse(apiURL)
	u.Path = strings.TrimRight(u.Path, "/") + "/api/v1/monitor/query/items"
	_, err := httpClient.Get(u.String())
	if err == nil {
		t.Log("success")
	} else {
		t.Error("faild")
	}
}

func TestGetUsers(t *testing.T) {
	httpClient := http.DefaultClient
	u, _ := url.Parse(apiURL)
	u.Path = strings.TrimRight(u.Path, "/") + "/api/v1/monitor/users"
	q := u.Query()
	q.Set("start", "1481853425")
	q.Set("end", "1481853425")
	q.Set("step", "30s")
	u.RawQuery = q.Encode()
	_, err := httpClient.Get(u.String())
	if err == nil {
		t.Log("success")
	} else {
		t.Error("faild")
	}
}
