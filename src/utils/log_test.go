package utils

import (
	"golang.org/x/net/context"
	"net/http"
	"testing"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
)

func TestGetLogger(t *testing.T) {
	GetLogger(context.Background())
}

func TestWithLogger(t *testing.T) {
	c := context.Background()
	WithLogger(c, GetLogger(c))
}

func LogT(ctx *gin.Context) {
	f := Ginrus(log.StandardLogger(), time.RFC3339Nano, false)
	f(ctx)
}

func TestGinrus(t *testing.T) {
	http.Get(baseURL + "/v1/test/log")
}
