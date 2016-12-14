package middleware

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
)

var (
	baseUrl string
	server  *httptest.Server
)

func startHttpServer() *httptest.Server {
	router := gin.New()
	router.GET("/ping", Authenticate, CORSMiddleware(), packPage, ping)

	return httptest.NewServer(router)
}

func ping(ctx *gin.Context) {
	ctx.String(200, "success")
}

func TestMain(m *testing.M) {
	server = startHttpServer()
	baseUrl = server.URL
	ret := m.Run()
	server.Close()
	os.Exit(ret)
}

func TestAll(t *testing.T) {
	_, err := http.Get(baseUrl + "/ping")
	if err == nil {
		t.Log("success")
	} else {
		t.Error("faild")
	}

	_, err = http.Get(baseUrl + "/ping?to=1")
	if err == nil {
		t.Log("success")
	} else {
		t.Error("faild")
	}

	_, err = http.Get(baseUrl + "/ping?to=test")
	if err == nil {
		t.Log("success")
	} else {
		t.Error("faild")
	}

	_, err = http.Get(baseUrl + "/ping?from=test")
	if err == nil {
		t.Log("success")
	} else {
		t.Error("faild")
	}

	_, err = http.Get(baseUrl + "/ping?from=1")
	if err == nil {
		t.Log("success")
	} else {
		t.Error("faild")
	}
}
