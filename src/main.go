package main

import (
	"net/http"
	"time"

	"github.com/Dataman-Cloud/log-proxy/src/config"
	"github.com/Dataman-Cloud/log-proxy/src/router"
	"github.com/Dataman-Cloud/log-proxy/src/router/middleware"

	log "github.com/Sirupsen/logrus"
)

func main() {
	log.Infof("http server: %s start...", config.GetConfig().ADDR)

	server := &http.Server{
		Addr:           config.GetConfig().ADDR,
		Handler:        router.Router(middleware.Authenticate),
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("http listen server error: %v", err)
	}
}
