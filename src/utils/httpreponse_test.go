package utils

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

var (
	baseURL string
	server  *httptest.Server
)

func TestNewError(t *testing.T) {
	err := NewError("503-10000", errors.New("error"))
	fmt.Println(err.Error())

}

func startHTTPServer() *httptest.Server {
	router := gin.New()
	v1 := router.Group("v1")
	{
		v1.GET("/test/ok", OkResponse)
		v1.GET("/test/create", CreateResponse)
		v1.GET("/test/update", UpdateResponse)
		v1.DELETE("/test/delete", DeleteResponse)
		v1.GET("/test/error", ErrResponse)
		v1.GET("/test/log", LogT)
	}

	return httptest.NewServer(router)
}

func OkResponse(ctx *gin.Context) {
	Ok(ctx, "success")
}

func CreateResponse(ctx *gin.Context) {
	Create(ctx, "success")
}

func UpdateResponse(ctx *gin.Context) {
	Update(ctx, "success")
}

func DeleteResponse(ctx *gin.Context) {
	Delete(ctx, "success")
}

func ErrResponse(ctx *gin.Context) {
	ErrorResponse(ctx, errors.New("error"))
	ErrorResponse(ctx, NewError("503-10000", errors.New("error")))

}

func TestMain(m *testing.M) {
	server = startHTTPServer()
	baseURL = server.URL
	ret := m.Run()
	server.Close()
	os.Exit(ret)
}

func TestOk(t *testing.T) {
	resp, _ := http.Get(baseURL + "/v1/test/ok")

	assert.Equal(t, resp.StatusCode, http.StatusOK, "should be equal")
}

func TestCreate(t *testing.T) {
	resp, _ := http.Get(baseURL + "/v1/test/create")

	assert.Equal(t, resp.StatusCode, http.StatusCreated, "should be equal")
}

func TestUpdate(t *testing.T) {
	resp, _ := http.Get(baseURL + "/v1/test/update")

	assert.Equal(t, resp.StatusCode, http.StatusAccepted, "should be equal")
}

func TestDelete(t *testing.T) {
	req, _ := http.NewRequest("DELETE", baseURL+"/v1/test/delete", nil)
	resp, _ := http.DefaultClient.Do(req)

	assert.Equal(t, resp.StatusCode, http.StatusNoContent, "should be equal")
}

func TestErrorResponse(t *testing.T) {

	http.Get(baseURL + "/v1/test/error")

}
