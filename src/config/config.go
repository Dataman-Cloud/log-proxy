package config

import (
	"bufio"
	"flag"
	"io"
	"log"
	"os"
	"reflect"
	"strconv"
	"strings"
)

// Config defines the conf info
type Config struct {
	ADDR             string `require:"true"`
	ES_URL           string `require:"true"`
	SEARCH_DEBUG     bool
	PROMETHEUS_URL   string `require:"true"`
	ALERTMANAGER_URL string `require:"true"`
	FRONTEND_PATH    string
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
		rb, err := strconv.ParseBool(robj.Type().Field(i).Tag.Get("require"))
		if err == nil {
			if rb && os.Getenv(robj.Type().Field(i).Name) == "" {
				log.Fatalf("config field %s not setting", robj.Type().Field(i).Name)
			}
		}
		switch robj.Type().Field(i).Type.String() {
		case "string":
			robj.Field(i).Set(reflect.ValueOf(os.Getenv(robj.Type().Field(i).Name)))
		case "bool":
			if b, err := strconv.ParseBool(os.Getenv(robj.Type().Field(i).Name)); err == nil {
				robj.Field(i).Set(reflect.ValueOf(b))
			} else if rb {
				log.Fatalf("config %s value invalid", robj.Type().Field(i).Name)
			}
		case "int":
			if b, err := strconv.Atoi(os.Getenv(robj.Type().Field(i).Name)); err == nil {
				robj.Field(i).Set(reflect.ValueOf(b))
			} else if rb {
				log.Fatalf("config %s value invalid", robj.Type().Field(i).Name)
			}
		case "uint64":
			if b, err := strconv.ParseUint(os.Getenv(robj.Type().Field(i).Name), 10, 64); err == nil {
				robj.Field(i).Set(reflect.ValueOf(b))
			} else if rb {
				log.Fatalf("config %s value invalid", robj.Type().Field(i).Name)
			}
		case "int64":
			if b, err := strconv.ParseInt(os.Getenv(robj.Type().Field(i).Name), 10, 64); err == nil {
				robj.Field(i).Set(reflect.ValueOf(b))
			} else if rb {
				log.Fatalf("config %s value invalid", robj.Type().Field(i).Name)
			}
		}
	}
}
