package service

import (
	"testing"

	"github.com/Dataman-Cloud/log-proxy/src/config"
	"github.com/Dataman-Cloud/log-proxy/src/models"
)

func TestCreateAlert(t *testing.T) {
	config.InitConfig("../../env_file.template")
	service := NewEsService([]string{baseURL})
	if service.CreateAlert(new(models.Alert)) == nil {
		t.Log("success")
	} else {
		t.Error("faild")
	}
}

func TestGetAlerts(t *testing.T) {
	config.InitConfig("../../env_file.template")
	service := NewEsService([]string{baseURL})
	if _, err := service.GetAlerts(models.Page{}); err == nil {
		t.Log("success")
	} else {
		t.Error("faild")
	}
}

func TestGetAlert(t *testing.T) {
	config.InitConfig("../../env_file.template")
	service := NewEsService([]string{baseURL})
	if _, err := service.GetAlert("test"); err == nil {
		t.Log("success")
	} else {
		t.Error("faild")
	}
}

func TestDeleteAlert(t *testing.T) {
	config.InitConfig("../../env_file.template")
	service := NewEsService([]string{baseURL})
	if service.DeleteAlert("test") == nil {
		t.Log("success")
	} else {
		t.Error("faild")
	}
}

func TestUpdateAlert(t *testing.T) {
	config.InitConfig("../../env_file.template")
	service := NewEsService([]string{baseURL})
	if service.UpdateAlert(new(models.Alert)) == nil {
		t.Log("success")
	} else {
		t.Error("faild")
	}
}
