package main

import (
	"flag"
	"net/http"
	"path"
	"time"

	"github.com/Dataman-Cloud/log-proxy/src/config"
	"github.com/Dataman-Cloud/log-proxy/src/router"
	"github.com/Dataman-Cloud/log-proxy/src/router/middleware"
	"github.com/Dataman-Cloud/log-proxy/src/store/datastore"
	expr "github.com/Dataman-Cloud/log-proxy/src/utils/prometheusexpr"

	log "github.com/Sirupsen/logrus"
	"github.com/gin-gonic/contrib/static"
)

func main() {
	configFile := flag.String("config", "env_file", "config file path")
	flag.Parse()
	config.InitConfig(*configFile)
	config.LoadLogOptionalLabels()
	datastore.InitDB(config.GetConfig().DbDriver, config.GetConfig().DbDSN)
	err := expr.Exprs(config.GetConfig().QueryExprPATH)
	if err != nil {
		log.Fatal(err)
	}

	log.Infof("http server: %s start...", config.GetConfig().Addr)

	server := &http.Server{
		Addr: config.GetConfig().Addr,
		Handler: router.Router(middleware.Authenticate,
			static.Serve("/",
				static.LocalFile(path.Join(config.GetConfig().FrontendPath, "frontend"), true))),
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("http listen server error: %v", err)
	}
}
