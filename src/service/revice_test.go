package service

import (
	"testing"

	"github.com/Dataman-Cloud/log-proxy/src/config"
	"github.com/Dataman-Cloud/log-proxy/src/models"
)

func TestGetPrometheus(t *testing.T) {
	config.InitConfig("../../env_file.template")
	service := NewEsService([]string{baseURL})
	if _, err := service.GetPrometheus(models.Page{}); err != nil {
		t.Error("faild")
	} else {
		t.Log("success")
	}
}

func TestGetPrometheu(t *testing.T) {
	config.InitConfig("../../env_file.template")
	service := NewEsService([]string{baseURL})
	if _, err := service.GetPrometheu("AVj3kWyMIIGpJqE63T3m"); err != nil {
		t.Error("faild")
	} else {
		t.Log("success")
	}
}

func TestSavePrometheus(t *testing.T) {
	config.InitConfig("../../env_file.template")
	service := NewEsService([]string{baseURL})
	if service.SavePrometheus(map[string]interface{}{"test": "value"}) != nil {
		t.Error("faild")
	} else {
		t.Log("success")
	}
}
