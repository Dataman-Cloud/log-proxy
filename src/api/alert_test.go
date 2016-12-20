package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/Dataman-Cloud/log-proxy/src/models"

	"github.com/gin-gonic/gin"
)

func createAlert(ctx *gin.Context) {
	s.CreateAlert(ctx)
}

func deleteAlert(ctx *gin.Context) {
	s.DeleteAlert(ctx)
}

func getAlerts(ctx *gin.Context) {
	s.GetAlerts(ctx)
}

func getAlert(ctx *gin.Context) {
	s.GetAlert(ctx)
}

func updateAlert(ctx *gin.Context) {
	s.UpdateAlert(ctx)
}

func TestCreateAlert(t *testing.T) {
	if s == nil {
		s = GetSearch()
	}
	resp, _ := http.NewRequest("POST", apiURL+"/api/v1/monitor/alert", nil)
	req, err := http.DefaultClient.Do(resp)
	if err == nil && req.StatusCode == 400 {
		t.Log("success")
	} else {
		t.Error("faild")
	}

	alert := models.Alert{}
	data, _ := json.Marshal(alert)
	resp, _ = http.NewRequest("POST", apiURL+"/api/v1/monitor/alert", bytes.NewBuffer(data))
	req, err = http.DefaultClient.Do(resp)
	if err == nil && req.StatusCode == 400 {
		t.Log("success")
	} else {
		t.Error("faild")
	}

	alert = models.Alert{
		AppID: "test",
	}
	data, _ = json.Marshal(alert)
	resp, _ = http.NewRequest("POST", apiURL+"/api/v1/monitor/alert", bytes.NewBuffer(data))
	req, err = http.DefaultClient.Do(resp)
	if err == nil && req.StatusCode == 400 {
		t.Log("success")
	} else {
		t.Error("faild")
	}

	alert = models.Alert{
		AppID:   "test",
		Keyword: "test",
	}
	data, _ = json.Marshal(alert)
	resp, _ = http.NewRequest("POST", apiURL+"/api/v1/monitor/alert", bytes.NewBuffer(data))
	req, err = http.DefaultClient.Do(resp)
	if err == nil && req.StatusCode == 200 {
		t.Log("success")
	} else {
		t.Error("faild")
	}
}

func TestDeleteAlert(t *testing.T) {
	if s == nil {
		s = GetSearch()
	}
	s.KeywordFilter[""] = []string{"test", ""}
	req, _ := http.NewRequest("DELETE", apiURL+"/api/v1/monitor/alert/test", nil)
	resp, err := http.DefaultClient.Do(req)

	if err == nil && resp.StatusCode == 200 {
		t.Log("success")
	} else {
		t.Error("faild")
	}

}

func TestGetAlerts(t *testing.T) {
	if s == nil {
		s = GetSearch()
	}

	resp, err := http.Get(apiURL + "/api/v1/monitor/alert")
	if err == nil && resp.StatusCode == 200 {
		t.Log("success")
	} else {
		t.Error("faild")
	}
}

func TestGetAlert(t *testing.T) {

	if s == nil {
		s = GetSearch()
	}

	resp, err := http.Get(apiURL + "/api/v1/monitor/alert/test")
	if err == nil && resp.StatusCode == 200 {
		t.Log("success")
	} else {
		t.Error("faild")
	}
}

func TestUpdateAlert(t *testing.T) {
	if s == nil {
		s = GetSearch()
	}

	req, _ := http.NewRequest("PUT", apiURL+"/api/v1/monitor/alert", nil)
	resp, err := http.DefaultClient.Do(req)
	if err == nil && resp.StatusCode == 400 {
		t.Log("success")
	} else {
		t.Error("faild")
	}

	alert := models.Alert{}
	data, _ := json.Marshal(alert)
	req, _ = http.NewRequest("PUT", apiURL+"/api/v1/monitor/alert", bytes.NewBuffer(data))
	resp, err = http.DefaultClient.Do(req)
	if err == nil && resp.StatusCode == 400 {
		t.Log("success")
	} else {
		t.Error("faild")
	}

	alert = models.Alert{
		ID: "test",
	}
	data, _ = json.Marshal(alert)
	req, _ = http.NewRequest("PUT", apiURL+"/api/v1/monitor/alert", bytes.NewBuffer(data))
	resp, err = http.DefaultClient.Do(req)
	if err == nil && resp.StatusCode == 400 {
		t.Log("success")
	} else {
		t.Error("faild")
	}

	alert = models.Alert{
		ID:    "test",
		AppID: "appid",
	}
	data, _ = json.Marshal(alert)
	req, _ = http.NewRequest("PUT", apiURL+"/api/v1/monitor/alert", bytes.NewBuffer(data))
	resp, err = http.DefaultClient.Do(req)
	if err == nil && resp.StatusCode == 400 {
		t.Log("success")
	} else {
		t.Error("faild")
	}

	alert = models.Alert{
		ID:      "test",
		AppID:   "appid",
		Keyword: "test",
	}
	data, _ = json.Marshal(alert)
	req, _ = http.NewRequest("PUT", apiURL+"/api/v1/monitor/alert", bytes.NewBuffer(data))
	resp, err = http.DefaultClient.Do(req)
	if err == nil && resp.StatusCode == 200 {
		t.Log("success")
	} else {
		t.Error("faild")
	}
}
