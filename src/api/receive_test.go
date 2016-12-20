package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
)

func receiver(ctx *gin.Context) {
	s.Receiver(ctx)
}

func receiverlog(ctx *gin.Context) {
	s.ReceiverLog(ctx)
}

func getprometheus(ctx *gin.Context) {
	s.GetPrometheus(ctx)
}

func getprometheu(ctx *gin.Context) {
	s.GetPrometheu(ctx)
}

func TestReciver(t *testing.T) {
	if s == nil {
		s = GetSearch()
	}
	req, _ := http.NewRequest("POST", apiURL+"/v1/recive/prometheus", nil)
	resp, err := http.DefaultClient.Do(req)
	if err == nil && resp.StatusCode == 503 {
		t.Log("success")
	} else {
		t.Error("faild")
	}

	m := map[string][]interface{}{
		"alerts": []interface{}{
			map[string]interface{}{
				"labels": map[string]interface{}{
					"alertname": "test",
					"taskid":    "value",
				},
			},
		},
	}

	data, _ := json.Marshal(m)

	req, _ = http.NewRequest("POST", apiURL+"/v1/recive/prometheus", bytes.NewBuffer(data))
	resp, err = http.DefaultClient.Do(req)
	if err == nil && resp.StatusCode == 200 {
		t.Log("success")
	} else {
		t.Error("faild")
	}
}

func TestReciverLog(t *testing.T) {
	if s == nil {
		s = GetSearch()
	}

	m := make(map[string]interface{})
	data, _ := json.Marshal(m)

	req, _ := http.NewRequest("POST", apiURL+"/v1/recive/log", bytes.NewBuffer(data))
	resp, err := http.DefaultClient.Do(req)
	if err == nil && resp.StatusCode == 503 {
		t.Log("success")
	} else {
		t.Error("faild")
	}

	m = map[string]interface{}{
		"appid": "test",
	}

	data, _ = json.Marshal(m)
	req, _ = http.NewRequest("POST", apiURL+"/v1/recive/log", bytes.NewBuffer(data))
	resp, err = http.DefaultClient.Do(req)
	if err == nil && resp.StatusCode == 503 {
		t.Log("success")
	} else {
		t.Error("faild")
	}

	m = map[string]interface{}{
		"appid":  "test",
		"taskid": "test",
	}

	data, _ = json.Marshal(m)
	req, _ = http.NewRequest("POST", apiURL+"/v1/recive/log", bytes.NewBuffer(data))
	resp, err = http.DefaultClient.Do(req)
	if err == nil && resp.StatusCode == 503 {
		t.Log("success")
	} else {
		t.Error("faild")
	}

	m = map[string]interface{}{
		"appid":  "test",
		"taskid": "test",
		"path":   "test",
	}

	data, _ = json.Marshal(m)
	req, _ = http.NewRequest("POST", apiURL+"/v1/recive/log", bytes.NewBuffer(data))
	resp, err = http.DefaultClient.Do(req)
	if err == nil && resp.StatusCode == 503 {
		t.Log("success")
	} else {
		t.Error("faild")
	}

	m = map[string]interface{}{
		"appid":  "test",
		"taskid": "test",
		"path":   "test",
		"userid": "test",
	}

	data, _ = json.Marshal(m)
	req, _ = http.NewRequest("POST", apiURL+"/v1/recive/log", bytes.NewBuffer(data))
	resp, err = http.DefaultClient.Do(req)
	if err == nil && resp.StatusCode == 503 {
		t.Log("success")
	} else {
		t.Error("faild")
	}

	m = map[string]interface{}{
		"appid":     "test",
		"taskid":    "test",
		"path":      "test",
		"userid":    "test",
		"clusterid": "test",
	}

	data, _ = json.Marshal(m)
	req, _ = http.NewRequest("POST", apiURL+"/v1/recive/log", bytes.NewBuffer(data))
	resp, err = http.DefaultClient.Do(req)
	if err == nil && resp.StatusCode == 503 {
		t.Log("success")
	} else {
		t.Error("faild")
	}

	m = map[string]interface{}{
		"appid":     "test",
		"taskid":    "test",
		"path":      "test",
		"userid":    "test",
		"clusterid": "test",
		"offset":    111,
	}

	data, _ = json.Marshal(m)
	req, _ = http.NewRequest("POST", apiURL+"/v1/recive/log", bytes.NewBuffer(data))
	resp, err = http.DefaultClient.Do(req)
	if err == nil && resp.StatusCode == 503 {
		t.Log("success")
	} else {
		t.Error("faild")
	}

	m = map[string]interface{}{
		"appid":     "test",
		"taskid":    "test",
		"path":      "test",
		"userid":    "test",
		"clusterid": "test",
		"offset":    111,
		"message":   "get",
	}

	data, _ = json.Marshal(m)
	req, _ = http.NewRequest("POST", apiURL+"/v1/recive/log", bytes.NewBuffer(data))
	resp, err = http.DefaultClient.Do(req)
	if err == nil && resp.StatusCode == 200 {
		t.Log("success")
	} else {
		t.Error("faild")
	}

	m = map[string]interface{}{
		"appid":     "test",
		"taskid":    "test",
		"path":      "test",
		"userid":    "test",
		"clusterid": "test",
		"offset":    111,
		"message":   "get",
	}
	s.KeywordFilter["testtest"] = []string{"get"}

	data, _ = json.Marshal(m)
	req, _ = http.NewRequest("POST", apiURL+"/v1/recive/log", bytes.NewBuffer(data))
	resp, err = http.DefaultClient.Do(req)
	if err == nil && resp.StatusCode == 200 {
		t.Log("success")
	} else {
		t.Error("faild")
	}
}

func TestGetPrometheus(t *testing.T) {
	if s == nil {
		s = GetSearch()
	}
	req, _ := http.NewRequest("GET", apiURL+"/api/v1/monitor/prometheus", nil)
	resp, err := http.DefaultClient.Do(req)
	if err == nil && resp.StatusCode == 200 {
		t.Log("success")
	} else {
		t.Error("faild")
	}
}

func TestGetPrometheu(t *testing.T) {
	if s == nil {
		s = GetSearch()
	}
	req, _ := http.NewRequest("GET", apiURL+"/api/v1/monitor/prometheus/test", nil)
	resp, err := http.DefaultClient.Do(req)
	if err == nil && resp.StatusCode == 200 {
		t.Log("success")
	} else {
		t.Error("faild")
	}
}
