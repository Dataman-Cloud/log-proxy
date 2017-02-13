package config

import (
	"bufio"
	"io"
	"log"
	"os"
	"reflect"
	"strconv"
	"strings"
)

// Config defines the conf info
type Config struct {
	Addr            string `require:"true" alias:"ADDR"`
	EsURL           string `require:"true" alias:"ES_URL"`
	SearchDebug     bool   `alias:"SEARCH_DEBUG"`
	PrometheusURL   string `require:"true" alias:"PROMETHEUS_URL"`
	AlertManagerURL string `require:"true" alias:"ALERTMANAGER_URL"`
	FrontendPath    string `alias:"FRONTEND_PATH"`
	MarathonURL     string `alias:"MARATHON_URL"`
	NotificationURL string `alias:"NOTIFICATION_URL"`
	DbDSN           string `alias:"DB_DSN"`
	DbDriver        string `alias:"DB_DRIVER"`
}

var c *Config

// GetConfig get config data
func GetConfig() *Config {
	return c
}

// InitConfig init config
func InitConfig(file string) {
	c = new(Config)

	f, err := os.Open(file)
	if err != nil {
		log.Printf("open config file %s error: %v", file, err)
		return
	}
	defer func() {
		if err != nil {
			f.Close()
		}
	}()

	rd := bufio.NewReader(f)
	for {
		line, err := rd.ReadString('\n')
		if err != nil || io.EOF == err {
			break
		}
		kv := strings.SplitN(string(line), "=", 2)
		if len(kv) == 2 && os.Getenv(kv[0]) == "" {
			os.Setenv(kv[0], strings.TrimRight(kv[1], "\n"))
		}
	}
	LoadConfig()
}

// LoadConfig load config data
func LoadConfig() {
	robj := reflect.ValueOf(c).Elem()
	for i := 0; i < robj.NumField(); i++ {
		rb, err := strconv.ParseBool(robj.Type().Field(i).Tag.Get("require"))
		if err == nil {
			if rb && os.Getenv(robj.Type().Field(i).Tag.Get("alias")) == "" {
				log.Fatalf("config field %s not setting", robj.Type().Field(i).Tag.Get("alias"))
			}
		}
		switch robj.Type().Field(i).Type.String() {
		case "string":
			robj.Field(i).Set(reflect.ValueOf(os.Getenv(robj.Type().Field(i).Tag.Get("alias"))))
		case "bool":
			if b, err := strconv.ParseBool(os.Getenv(robj.Type().Field(i).Tag.Get("alias"))); err == nil {
				robj.Field(i).Set(reflect.ValueOf(b))
			} else if rb {
				log.Fatalf("config %s value invalid", robj.Type().Field(i).Tag.Get("alias"))
			}
		}
	}
}
