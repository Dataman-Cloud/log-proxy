package api

import (
	"bytes"
	"container/list"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
)

func receiverlog(ctx *gin.Context) {
	s.ReceiverLog(ctx)
}

func TestReciverLog(t *testing.T) {
	if s == nil {
		s = GetSearch()
	}

	m := make(map[string]interface{})
	data, _ := json.Marshal(m)

	req, _ := http.NewRequest("POST", apiURL+"/v1/receive/log", bytes.NewBuffer(data))
	resp, err := http.DefaultClient.Do(req)
	if err == nil && resp.StatusCode == 503 {
		t.Log("success")
	} else {
		t.Error("faild")
	}

	m = map[string]interface{}{
		"app": "test",
	}

	data, _ = json.Marshal(m)
	req, _ = http.NewRequest("POST", apiURL+"/v1/receive/log", bytes.NewBuffer(data))
	resp, err = http.DefaultClient.Do(req)
	if err == nil && resp.StatusCode == 503 {
		t.Log("success")
	} else {
		t.Error("faild")
	}

	m = map[string]interface{}{
		"app":  "test",
		"task": "test",
	}

	data, _ = json.Marshal(m)
	req, _ = http.NewRequest("POST", apiURL+"/v1/receive/log", bytes.NewBuffer(data))
	resp, err = http.DefaultClient.Do(req)
	if err == nil && resp.StatusCode == 503 {
		t.Log("success")
	} else {
		t.Error("faild")
	}

	m = map[string]interface{}{
		"app":  "test",
		"task": "test",
		"path": "test",
	}

	data, _ = json.Marshal(m)
	req, _ = http.NewRequest("POST", apiURL+"/v1/receive/log", bytes.NewBuffer(data))
	resp, err = http.DefaultClient.Do(req)
	if err == nil && resp.StatusCode == 503 {
		t.Log("success")
	} else {
		t.Error("faild")
	}

	m = map[string]interface{}{
		"app":  "test",
		"task": "test",
		"path": "test",
		"user": "test",
	}

	data, _ = json.Marshal(m)
	req, _ = http.NewRequest("POST", apiURL+"/v1/receive/log", bytes.NewBuffer(data))
	resp, err = http.DefaultClient.Do(req)
	if err == nil && resp.StatusCode == 503 {
		t.Log("success")
	} else {
		t.Error("faild")
	}

	m = map[string]interface{}{
		"app":     "test",
		"task":    "test",
		"path":    "test",
		"user":    "test",
		"cluster": "test",
	}

	data, _ = json.Marshal(m)
	req, _ = http.NewRequest("POST", apiURL+"/v1/receive/log", bytes.NewBuffer(data))
	resp, err = http.DefaultClient.Do(req)
	if err == nil && resp.StatusCode == 503 {
		t.Log("success")
	} else {
		t.Error("faild")
	}

	m = map[string]interface{}{
		"app":     "test",
		"task":    "test",
		"path":    "test",
		"user":    "test",
		"cluster": "test",
		"offset":  111,
	}

	data, _ = json.Marshal(m)
	req, _ = http.NewRequest("POST", apiURL+"/v1/receive/log", bytes.NewBuffer(data))
	resp, err = http.DefaultClient.Do(req)
	if err == nil && resp.StatusCode == 503 {
		t.Log("success")
	} else {
		t.Error("faild")
	}

	m = map[string]interface{}{
		"app":     "test",
		"task":    "test",
		"path":    "test",
		"user":    "test",
		"cluster": "test",
		"offset":  111,
		"message": "get",
	}

	data, _ = json.Marshal(m)
	req, _ = http.NewRequest("POST", apiURL+"/v1/receive/log", bytes.NewBuffer(data))
	resp, err = http.DefaultClient.Do(req)
	if err == nil && resp.StatusCode == 200 {
		t.Log("success")
	} else {
		t.Error("faild")
	}

	m = map[string]interface{}{
		"app":     "test",
		"task":    "test",
		"path":    "test",
		"user":    "test",
		"cluster": "test",
		"offset":  111,
		"message": "get",
	}
	//s.KeywordFilter["testtest"] = []string{"get"}

	l := list.New()
	l.PushFront("get")
	s.KeywordFilter["testtest"] = l

	data, _ = json.Marshal(m)
	req, _ = http.NewRequest("POST", apiURL+"/v1/receive/log", bytes.NewBuffer(data))
	resp, err = http.DefaultClient.Do(req)
	if err == nil && resp.StatusCode == 200 {
		t.Log("success")
	} else {
		t.Error("faild")
	}
}
