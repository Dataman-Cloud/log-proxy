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

func TestReceiver(t *testing.T) {
	if s == nil {
		s = GetSearch()
	}
	req, _ := http.NewRequest("POST", apiURL+"/v1/receive/prometheus", nil)
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
					"task":      "value",
				},
			},
		},
	}

	data, _ := json.Marshal(m)

	req, _ = http.NewRequest("POST", apiURL+"/v1/receive/prometheus", bytes.NewBuffer(data))
	resp, err = http.DefaultClient.Do(req)
	if err == nil && resp.StatusCode == 200 {
		t.Log("success")
	} else {
		t.Error("faild")
	}
}
