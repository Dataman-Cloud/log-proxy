package config

import (
	"bufio"
	"flag"
	"io"
	"log"
	"os"
	"reflect"
	"strings"
)

// Config defines the conf info
type Config struct {
	ADDR           string
	ES_URL         string
	PROMETHEUS_URL string
}

var c *Config

func GetConfig() *Config {
	return c
}

func init() {
	c = new(Config)
	configFile := flag.String("config", "env_file", "config file path")
	flag.Parse()

	f, err := os.Open(*configFile)
	if err != nil {
		log.Printf("open config file %s error: %v", *configFile, err)
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

func LoadConfig() {
	robj := reflect.ValueOf(c).Elem()
	for i := 0; i < robj.NumField(); i++ {
		robj.Field(i).Set(reflect.ValueOf(os.Getenv(robj.Type().Field(i).Name)))
	}
}
