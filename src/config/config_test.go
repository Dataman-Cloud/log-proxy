package config

import (
	"testing"
)

func TestInitConfig(t *testing.T) {
	InitConfig("env_file")
	InitConfig("../../env_file.template")
	GetConfig()
}
