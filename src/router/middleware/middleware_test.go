package middleware

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
)

var (
	baseURL string
	server  *httptest.Server
)

func startHTTPServer() *httptest.Server {
	router := gin.New()
	router.GET("/ping", Authenticate, CORSMiddleware(), packPage, ping)

	return httptest.NewServer(router)
}

func ping(ctx *gin.Context) {
	ctx.String(200, "success")
}

func TestMain(m *testing.M) {
	server = startHTTPServer()
	baseURL = server.URL
	ret := m.Run()
	server.Close()
	os.Exit(ret)
}

func TestAll(t *testing.T) {
	_, err := http.Get(baseURL + "/ping")
	if err == nil {
		t.Log("success")
	} else {
		t.Error("faild")
	}

	_, err = http.Get(baseURL + "/ping?to=1")
	if err == nil {
		t.Log("success")
	} else {
		t.Error("faild")
	}

	_, err = http.Get(baseURL + "/ping?to=test")
	if err == nil {
		t.Log("success")
	} else {
		t.Error("faild")
	}

	_, err = http.Get(baseURL + "/ping?from=test")
	if err == nil {
		t.Log("success")
	} else {
		t.Error("faild")
	}

	_, err = http.Get(baseURL + "/ping?from=1")
	if err == nil {
		t.Log("success")
	} else {
		t.Error("faild")
	}
}
