package api

import (
	"bytes"
	"container/list"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/Dataman-Cloud/log-proxy/src/config"
	"github.com/Dataman-Cloud/log-proxy/src/models"
)

func TestCreateAlert(t *testing.T) {
	sr := startErrorClient()
	config.GetConfig().EsURL = sr.URL
	baseURL = sr.URL
	s = GetSearch()
	se := startAPIServer(s)
	//s.KeywordFilter = map[string][]string{}
	s.KeywordFilter = make(map[string]*list.List)
	alert := models.Alert{
		AppID:   "test",
		Keyword: "test1",
	}
	data, _ := json.Marshal(alert)
	resp, _ := http.NewRequest("POST", se.URL+"/api/v1/monitor/alert", bytes.NewBuffer(data))
	req, err := http.DefaultClient.Do(resp)
	if err == nil && req.StatusCode == 503 {
		t.Log("success")
	} else {
		t.Error("faild")
	}

	sr = startHTTPServer()
	config.GetConfig().EsURL = sr.URL
	baseURL = sr.URL
	s = GetSearch()
	se = startAPIServer(s)
	resp, _ = http.NewRequest("POST", se.URL+"/api/v1/monitor/alert", nil)
	req, err = http.DefaultClient.Do(resp)
	if err == nil && req.StatusCode == 400 {
		t.Log("success")
	} else {
		t.Error("faild")
	}

	alert = models.Alert{}
	data, _ = json.Marshal(alert)
	resp, _ = http.NewRequest("POST", se.URL+"/api/v1/monitor/alert", bytes.NewBuffer(data))
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
	resp, _ = http.NewRequest("POST", se.URL+"/api/v1/monitor/alert", bytes.NewBuffer(data))
	req, err = http.DefaultClient.Do(resp)
	if err == nil && req.StatusCode == 400 {
		t.Log("success")
	} else {
		t.Error("faild")
	}

	l := list.New()
	l.PushFront("test")
	s.KeywordFilter = map[string]*list.List{
		"test": l,
	}
	alert = models.Alert{
		AppID:   "test",
		Keyword: "test",
	}
	data, _ = json.Marshal(alert)
	resp, _ = http.NewRequest("POST", se.URL+"/api/v1/monitor/alert", bytes.NewBuffer(data))
	req, err = http.DefaultClient.Do(resp)
	if err == nil && req.StatusCode == 503 {
		t.Log("success")
	} else {
		t.Error("faild")
	}

	//s.KeywordFilter = map[string][]string{}
	s.KeywordFilter = make(map[string]*list.List)
	s.KeywordFilter["test"] = list.New()
	alert = models.Alert{
		AppID:   "test",
		Keyword: "test",
	}
	data, _ = json.Marshal(alert)
	resp, _ = http.NewRequest("POST", se.URL+"/api/v1/monitor/alert", bytes.NewBuffer(data))
	req, err = http.DefaultClient.Do(resp)
	if err == nil && req.StatusCode == 200 {
		t.Log("success")
	} else {
		t.Error("faild")
	}
}

func TestDeleteAlert(t *testing.T) {
	sr := startHTTPServer()
	config.GetConfig().EsURL = sr.URL
	baseURL = sr.URL
	s = GetSearch()
	se := startAPIServer(s)
	l := list.New()
	l.PushFront("test")
	l.PushFront("")
	s.KeywordFilter[""] = l
	req, _ := http.NewRequest("DELETE", se.URL+"/api/v1/monitor/alert/test", nil)
	resp, err := http.DefaultClient.Do(req)

	if err == nil && resp.StatusCode == 200 {
		t.Log("success")
	} else {
		t.Error("faild")
	}

}

func TestGetAlerts(t *testing.T) {
	sr := startHTTPServer()
	config.GetConfig().EsURL = sr.URL
	baseURL = sr.URL
	s = GetSearch()
	se := startAPIServer(s)

	resp, err := http.Get(se.URL + "/api/v1/monitor/alert")
	if err == nil && resp.StatusCode == 200 {
		t.Log("success")
	} else {
		t.Error("faild")
	}
}

func TestGetAlert(t *testing.T) {
	sr := startHTTPServer()
	config.GetConfig().EsURL = sr.URL
	baseURL = sr.URL
	s = GetSearch()
	se := startAPIServer(s)

	resp, err := http.Get(se.URL + "/api/v1/monitor/alert/test")
	if err == nil && resp.StatusCode == 200 {
		t.Log("success")
	} else {
		t.Error("faild")
	}
}

func TestUpdateAlert(t *testing.T) {
	sr := startHTTPServer()
	config.GetConfig().EsURL = sr.URL
	baseURL = sr.URL
	s = GetSearch()
	se := startAPIServer(s)

	req, _ := http.NewRequest("PUT", se.URL+"/api/v1/monitor/alert", nil)
	resp, err := http.DefaultClient.Do(req)
	if err == nil && resp.StatusCode == 400 {
		t.Log("success")
	} else {
		t.Error("faild")
	}

	alert := models.Alert{}
	data, _ := json.Marshal(alert)
	req, _ = http.NewRequest("PUT", se.URL+"/api/v1/monitor/alert", bytes.NewBuffer(data))
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
	req, _ = http.NewRequest("PUT", se.URL+"/api/v1/monitor/alert", bytes.NewBuffer(data))
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
	req, _ = http.NewRequest("PUT", se.URL+"/api/v1/monitor/alert", bytes.NewBuffer(data))
	resp, err = http.DefaultClient.Do(req)
	if err == nil && resp.StatusCode == 400 {
		t.Log("success")
	} else {
		t.Error("faild")
	}

	s.KeywordFilter = make(map[string]*list.List)
	s.KeywordFilter[""] = list.New()
	alert = models.Alert{
		ID:      "test",
		AppID:   "appid",
		Keyword: "test",
	}
	data, _ = json.Marshal(alert)
	req, _ = http.NewRequest("PUT", se.URL+"/api/v1/monitor/alert", bytes.NewBuffer(data))
	resp, err = http.DefaultClient.Do(req)
	if err == nil && resp.StatusCode == 200 {
		t.Log("success")
	} else {
		t.Error("faild")
	}

	l := list.New()
	l.PushFront("test")
	s.KeywordFilter[""] = l
	req, _ = http.NewRequest("PUT", se.URL+"/api/v1/monitor/alert", bytes.NewBuffer(data))
	resp, err = http.DefaultClient.Do(req)
	if err == nil && resp.StatusCode == 200 {
		t.Log("success")
	} else {
		t.Error("faild")
	}
}
