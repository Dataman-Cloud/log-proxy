package service

import (
	"testing"

	"github.com/Dataman-Cloud/log-proxy/src/config"
	"github.com/Dataman-Cloud/log-proxy/src/models"
)

func TestCreateFilter(t *testing.T) {
	config.InitConfig("../../env_file.template")
	service := NewEsService([]string{baseURL})
	if service.CreateFilter(new(models.KWFilter)) == nil {
		t.Log("success")
	} else {
		t.Error("faild")
	}
}

func TestGetFilters(t *testing.T) {
	config.InitConfig("../../env_file.template")
	service := NewEsService([]string{baseURL})
	if _, err := service.GetFilters(models.Page{}); err == nil {
		t.Log("success")
	} else {
		t.Error("faild")
	}
}

func TestGetFilter(t *testing.T) {
	config.InitConfig("../../env_file.template")
	service := NewEsService([]string{baseURL})
	if _, err := service.GetFilter("test"); err == nil {
		t.Log("success")
	} else {
		t.Error("faild")
	}
}

func TestDeleteFilter(t *testing.T) {
	config.InitConfig("../../env_file.template")
	service := NewEsService([]string{baseURL})
	if service.DeleteFilter("test") == nil {
		t.Log("success")
	} else {
		t.Error("faild")
	}
}

func TestUpdateFilter(t *testing.T) {
	config.InitConfig("../../env_file.template")
	service := NewEsService([]string{baseURL})
	if service.UpdateFilter(new(models.KWFilter)) == nil {
		t.Log("success")
	} else {
		t.Error("faild")
	}
}
