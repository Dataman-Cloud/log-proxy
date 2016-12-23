package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func getMonitorSilences(ctx *gin.Context) {
	mo.GetSilences(ctx)
}

func getMonitorSilence(ctx *gin.Context) {
	mo.GetSilence(ctx)
}

func createMonitorSilence(ctx *gin.Context) {
	mo.CreateSilence(ctx)
}

func updateMonitorSilence(ctx *gin.Context) {
	mo.UpdateSilence(ctx)
}

func deleteMonitorSilence(ctx *gin.Context) {
	mo.DeleteSilence(ctx)
}

func TestMonitorGetSilences(t *testing.T) {
	expectResult := 200
	httpClient := http.DefaultClient
	u, _ := url.Parse(apiURL)
	u.Path = strings.TrimRight(u.Path, "/") + "/api/v1/monitor/silences"
	resp, _ := httpClient.Get(u.String())
	if resp.StatusCode != expectResult {
		t.Errorf("Expect query param metric is %v, got %v", expectResult, resp.StatusCode)
	}
}

func TestMonitorGetSilence(t *testing.T) {
	expectResult := 200
	httpClient := http.DefaultClient
	u, _ := url.Parse(apiURL)
	u.Path = strings.TrimRight(u.Path, "/") + "/api/v1/monitor/silence"
	resp, _ := httpClient.Get(u.String())
	if resp.StatusCode != expectResult {
		t.Errorf("Expect query param metric is %v, got %v", expectResult, resp.StatusCode)
	}
}

type Silence struct {
	Matchers  []*Matcher `json:"matchers"`
	StartsAt  string     `json:"startsAt"`
	EndsAt    string     `json:"endsAt"`
	CreatedBy string     `json:"createdBy"`
	Comment   string     `json:"comment"`
}

type Matcher struct {
	Name    string `json:"name"`
	Value   string `json:"value"`
	IsRegex bool   `json:"isRegex"`
}

func initSlience() *Silence {
	matcher := &Matcher{
		Name:    "alertname",
		Value:   "cpu_usage",
		IsRegex: false,
	}
	matchers := []*Matcher{matcher}
	silence := &Silence{
		Matchers:  matchers,
		StartsAt:  "2016-12-05T11:08:00.000Z",
		EndsAt:    "2016-12-05T15:08:00.000Z",
		CreatedBy: "yqguo@dataman-inc.com",
		Comment:   "this is a test",
	}
	return silence
}

func TestCreateSilence(t *testing.T) {
	path := apiURL + "/api/v1/monitor/silence"
	resp, _ := http.NewRequest("POST", path, nil)
	req, err := http.DefaultClient.Do(resp)
	if err == nil && req.StatusCode == 400 {
		t.Log("success")
	} else {
		t.Errorf("Expect get status code 400, got %v", req.StatusCode)
	}

	//silence := initSlience()
	//data, _ := json.Marshal(silence)
	str := `{"createdBy":"12@123.com","comment":"asdfasdf","endsAt":"2016-12-27T08:21:42.000Z","startsAt":"2016-12-27T06:21:42.000Z","matchers":[{"name":"asf","value":"asdf"}]}`
	resp, _ = http.NewRequest("POST", path, bytes.NewBuffer([]byte(str)))
	req, err = http.DefaultClient.Do(resp)
	if err == nil && req.StatusCode == 200 {
		t.Log("success")
	} else {
		t.Errorf("Expect get status code 200, got %v", req.StatusCode)
	}

	str = `{"comment":"asdfasdf","endsAt":"2016-12-27T08:21:42.000Z","startsAt":"2016-12-27T06:21:42.000Z","matchers":[{"name":"asf","value":"asdf"}]}`
	resp, _ = http.NewRequest("POST", path, bytes.NewBuffer([]byte(str)))
	req, err = http.DefaultClient.Do(resp)
	if err == nil && req.StatusCode == 503 {
		t.Log("success")
	} else {
		t.Errorf("Expect get status code 200, got %v", req.StatusCode)
	}

	str = `{"createdBy":"12@123.com","endsAt":"2016-12-27T08:21:42.000Z","startsAt":"2016-12-27T06:21:42.000Z","matchers":[{"name":"asf","value":"asdf"}]}`
	resp, _ = http.NewRequest("POST", path, bytes.NewBuffer([]byte(str)))
	req, err = http.DefaultClient.Do(resp)
	if err == nil && req.StatusCode == 503 {
		t.Log("success")
	} else {
		t.Errorf("Expect get status code 200, got %v", req.StatusCode)
	}

	str = `{"createdBy":"12@123.com","comment":"asdfasdf","startsAt":"2016-12-27T06:21:42.000Z","matchers":[{"name":"asf","value":"asdf"}]}`
	resp, _ = http.NewRequest("POST", path, bytes.NewBuffer([]byte(str)))
	req, err = http.DefaultClient.Do(resp)
	if err == nil && req.StatusCode == 503 {
		t.Log("success")
	} else {
		t.Errorf("Expect get status code 200, got %v", req.StatusCode)
	}

	str = `{"createdBy":"12@123.com","comment":"asdfasdf","endsAt":"2016-12-27T08:21:42.000Z","matchers":[{"name":"asf","value":"asdf"}]}`
	resp, _ = http.NewRequest("POST", path, bytes.NewBuffer([]byte(str)))
	req, err = http.DefaultClient.Do(resp)
	if err == nil && req.StatusCode == 503 {
		t.Log("success")
	} else {
		t.Errorf("Expect get status code 200, got %v", req.StatusCode)
	}

	str = `{"createdBy":"12@123.com","comment":"asdfasdf","endsAt":"2016-12-27T08:21:42.000Z","startsAt":"2016-12-27T06:21:42.000Z"}`
	resp, _ = http.NewRequest("POST", path, bytes.NewBuffer([]byte(str)))
	req, err = http.DefaultClient.Do(resp)
	if err == nil && req.StatusCode == 503 {
		t.Log("success")
	} else {
		t.Errorf("Expect get status code 200, got %v", req.StatusCode)
	}
}

func TestUpdateSilence(t *testing.T) {
	path := apiURL + "/api/v1/monitor/silence"
	resp, _ := http.NewRequest("PUT", path, nil)
	req, err := http.DefaultClient.Do(resp)
	if err == nil && req.StatusCode == 400 {
		t.Log("success")
	} else {
		t.Errorf("Expect get status code 400, got %v", req.StatusCode)
	}

	silence := initSlience()
	data, _ := json.Marshal(silence)
	resp, _ = http.NewRequest("PUT", path, bytes.NewBuffer(data))
	req, err = http.DefaultClient.Do(resp)
	if err == nil && req.StatusCode == 200 {
		t.Log("success")
	} else {
		t.Errorf("Expect get status code 200, got %v", req.StatusCode)
	}

	str := `{"comment":"asdfasdf","endsAt":"2016-12-27T08:21:42.000Z","startsAt":"2016-12-27T06:21:42.000Z","matchers":[{"name":"asf","value":"asdf"}]}`
	resp, _ = http.NewRequest("PUT", path, bytes.NewBuffer([]byte(str)))
	req, err = http.DefaultClient.Do(resp)
	if err == nil && req.StatusCode == 503 {
		t.Log("success")
	} else {
		t.Errorf("Expect get status code 200, got %v", req.StatusCode)
	}

	str = `{"createdBy":"12@123.com","endsAt":"2016-12-27T08:21:42.000Z","startsAt":"2016-12-27T06:21:42.000Z","matchers":[{"name":"asf","value":"asdf"}]}`
	resp, _ = http.NewRequest("PUT", path, bytes.NewBuffer([]byte(str)))
	req, err = http.DefaultClient.Do(resp)
	if err == nil && req.StatusCode == 503 {
		t.Log("success")
	} else {
		t.Errorf("Expect get status code 200, got %v", req.StatusCode)
	}

	str = `{"createdBy":"12@123.com","comment":"asdfasdf","startsAt":"2016-12-27T06:21:42.000Z","matchers":[{"name":"asf","value":"asdf"}]}`
	resp, _ = http.NewRequest("PUT", path, bytes.NewBuffer([]byte(str)))
	req, err = http.DefaultClient.Do(resp)
	if err == nil && req.StatusCode == 503 {
		t.Log("success")
	} else {
		t.Errorf("Expect get status code 200, got %v", req.StatusCode)
	}

	str = `{"createdBy":"12@123.com","comment":"asdfasdf","endsAt":"2016-12-27T08:21:42.000Z","matchers":[{"name":"asf","value":"asdf"}]}`
	resp, _ = http.NewRequest("PUT", path, bytes.NewBuffer([]byte(str)))
	req, err = http.DefaultClient.Do(resp)
	if err == nil && req.StatusCode == 503 {
		t.Log("success")
	} else {
		t.Errorf("Expect get status code 200, got %v", req.StatusCode)
	}

	str = `{"createdBy":"12@123.com","comment":"asdfasdf","endsAt":"2016-12-27T08:21:42.000Z","startsAt":"2016-12-27T06:21:42.000Z"}`
	resp, _ = http.NewRequest("PUT", path, bytes.NewBuffer([]byte(str)))
	req, err = http.DefaultClient.Do(resp)
	if err == nil && req.StatusCode == 503 {
		t.Log("success")
	} else {
		t.Errorf("Expect get status code 200, got %v", req.StatusCode)
	}
}

func TestDeleteSilence(t *testing.T) {
	path := apiURL + "/api/v1/monitor/silence/test"
	resp, _ := http.NewRequest("DELETE", path, nil)
	req, err := http.DefaultClient.Do(resp)
	if err == nil && req.StatusCode == 200 {
		t.Log("success")
	} else {
		t.Errorf("Expect get status code 400, got %v", req.StatusCode)
	}

}
