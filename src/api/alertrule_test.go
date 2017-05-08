package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Dataman-Cloud/log-proxy/src/models"
	"github.com/Dataman-Cloud/log-proxy/src/router/middleware"
	mock_alerter "github.com/Dataman-Cloud/log-proxy/src/service/mock_alerter"
	"github.com/golang/mock/gomock"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGetAlertIndicators(t *testing.T) {
	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()
	mockAlerter := mock_alerter.NewMockAlerter(mockCtl)

	m := NewMonitor()
	m.Alert = mockAlerter

	router := gin.New()
	router.Use(middleware.CORSMiddleware())
	router.GET("/indicators", m.GetAlertIndicators)
	testServer := httptest.NewServer(router)
	assert.NotNil(t, testServer)
	defer testServer.Close()

	mockAlerter.EXPECT().GetAlertIndicators().Return(nil).Times(1)

	resp, err := http.Get(testServer.URL + "/indicators")
	if err != nil {
		return
	}
	defer resp.Body.Close()
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestCreateAlertRule(t *testing.T) {
	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()
	mockAlerter := mock_alerter.NewMockAlerter(mockCtl)

	m := NewMonitor()
	m.Alert = mockAlerter

	router := gin.New()
	router.Use(middleware.CORSMiddleware())
	router.POST("/alert/rules", m.CreateAlertRule)
	testServer := httptest.NewServer(router)
	assert.NotNil(t, testServer)
	defer testServer.Close()

	rule := models.NewRule()
	rule.App = "work-nginx"
	rule.Pending = "5s"
	rule.Severity = "warning"
	rule.Indicator = "CPU使用百分比"
	rule.Aggregation = "max"
	rule.Comparison = ">"
	rule.Threshold = 60

	body, err := json.Marshal(rule)
	assert.Nil(t, err, "invalid param")
	req, err := http.NewRequest("POST", testServer.URL+"/alert/rules", strings.NewReader(string(body)))
	if err != nil {
		return
	}
	mockAlerter.EXPECT().CreateAlertRule(gomock.Any()).Return(nil, errors.New("err")).Times(1)
	//mockAlerter.EXPECT().GetIndicatorName(gomock.Any()).Return("", "", errors.New("err")).Times(1)

	req.Header.Set("Content-Type", "application/json")
	httpClient := http.DefaultClient
	resp, err := httpClient.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	assert.Nil(t, err)
	assert.Equal(t, http.StatusServiceUnavailable, resp.StatusCode)
}

func TestCreateAlertRuleErrorBody(t *testing.T) {
	m := NewMonitor()
	router := gin.New()
	router.Use(middleware.CORSMiddleware())
	router.POST("/alert/rules", m.CreateAlertRule)
	testServer := httptest.NewServer(router)
	assert.NotNil(t, testServer)
	defer testServer.Close()

	rule := models.NewRule()
	rule.App = "work-nginx"
	rule.Pending = "5s"
	rule.Severity = "warning"
	rule.Indicator = "CPU使用百分比"
	rule.Aggregation = "max"
	rule.Comparison = ">"
	rule.Threshold = 60

	//body, err := json.Marshal(rule)
	//assert.Nil(t, err, "invalid param")

	req, err := http.NewRequest("POST", testServer.URL+"/alert/rules", strings.NewReader(""))
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")
	httpClient := http.DefaultClient
	resp, err := httpClient.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	assert.Nil(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestDeleteAlertRule(t *testing.T) {
	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()
	mockAlerter := mock_alerter.NewMockAlerter(mockCtl)

	m := NewMonitor()
	m.Alert = mockAlerter

	router := gin.New()
	router.Use(middleware.CORSMiddleware())
	router.DELETE("/alert/rules/:id", m.DeleteAlertRule)
	testServer := httptest.NewServer(router)
	assert.NotNil(t, testServer)
	defer testServer.Close()

	rule := models.NewRule()
	rule.App = "work-nginx"
	rule.Pending = "5s"
	rule.Severity = "warning"
	rule.Indicator = "CPU使用百分比"
	rule.Aggregation = "max"
	rule.Comparison = ">"
	rule.Threshold = 60

	body, err := json.Marshal(rule)
	assert.Nil(t, err, "invalid param")

	mockAlerter.EXPECT().DeleteAlertRule(gomock.Any(), gomock.Any()).Return(nil).Times(1)

	req, err := http.NewRequest("DELETE", testServer.URL+"/alert/rules/1", strings.NewReader(string(body)))
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")
	httpClient := http.DefaultClient
	resp, err := httpClient.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestDeleteAlertRuleError(t *testing.T) {
	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()
	mockAlerter := mock_alerter.NewMockAlerter(mockCtl)

	m := NewMonitor()
	m.Alert = mockAlerter

	router := gin.New()
	router.Use(middleware.CORSMiddleware())
	router.DELETE("/alert/rules/:id", m.DeleteAlertRule)
	testServer := httptest.NewServer(router)
	assert.NotNil(t, testServer)
	defer testServer.Close()

	rule := models.NewRule()
	rule.App = "work-nginx"
	rule.Pending = "5s"
	rule.Severity = "warning"
	rule.Indicator = "CPU使用百分比"
	rule.Aggregation = "max"
	rule.Comparison = ">"
	rule.Threshold = 60

	body, err := json.Marshal(rule)
	assert.Nil(t, err, "invalid param")

	mockAlerter.EXPECT().DeleteAlertRule(gomock.Any(), gomock.Any()).Return(errors.New("err")).Times(1)

	req, err := http.NewRequest("DELETE", testServer.URL+"/alert/rules/abc", strings.NewReader(string(body)))
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")
	httpClient := http.DefaultClient
	resp, err := httpClient.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	assert.Nil(t, err)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusServiceUnavailable, resp.StatusCode)

	req, err = http.NewRequest("DELETE", testServer.URL+"/alert/rules/1", strings.NewReader(string(body)))
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")
	httpClient = http.DefaultClient
	resp, err = httpClient.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	assert.Nil(t, err)
	assert.Equal(t, http.StatusServiceUnavailable, resp.StatusCode)
}

func TestListAlertRulesError(t *testing.T) {
	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()
	mockAlerter := mock_alerter.NewMockAlerter(mockCtl)

	m := NewMonitor()
	m.Alert = mockAlerter

	router := gin.New()
	router.Use(middleware.CORSMiddleware())
	router.GET("/alert/rules", m.ListAlertRules)
	testServer := httptest.NewServer(router)
	assert.NotNil(t, testServer)
	defer testServer.Close()

	mockAlerter.EXPECT().ListAlertRules(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.New("err")).Times(1)

	resp, err := http.Get(testServer.URL + "/alert/rules")
	if err != nil {
		return
	}
	defer resp.Body.Close()
	assert.Nil(t, err)
	assert.Equal(t, http.StatusServiceUnavailable, resp.StatusCode)
}

func TestListAlertRules(t *testing.T) {
	m := NewMonitor()
	router := gin.New()
	router.Use(middleware.CORSMiddleware())
	router.GET("/alert/rules", m.ListAlertRules)
	testServer := httptest.NewServer(router)
	assert.NotNil(t, testServer)
	defer testServer.Close()

	resp, err := http.Get(testServer.URL + "/alert/rules")
	if err != nil {
		return
	}
	defer resp.Body.Close()
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestGetAlertRule(t *testing.T) {
	m := NewMonitor()
	router := gin.New()
	router.Use(middleware.CORSMiddleware())
	router.GET("/alert/rules/:id", m.GetAlertRule)
	testServer := httptest.NewServer(router)
	assert.NotNil(t, testServer)
	defer testServer.Close()

	resp, err := http.Get(testServer.URL + "/alert/rules/1")
	if err != nil {
		return
	}
	defer resp.Body.Close()
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestGetAlertRuleError(t *testing.T) {
	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()
	mockAlerter := mock_alerter.NewMockAlerter(mockCtl)

	m := NewMonitor()
	m.Alert = mockAlerter

	router := gin.New()
	router.Use(middleware.CORSMiddleware())
	router.GET("/alert/rules/:id", m.GetAlertRule)
	testServer := httptest.NewServer(router)
	assert.NotNil(t, testServer)
	defer testServer.Close()

	mockAlerter.EXPECT().GetAlertRule(gomock.Any()).Return(nil, errors.New("err")).Times(1)

	resp, err := http.Get(testServer.URL + "/alert/rules/abc")
	if err != nil {
		return
	}
	defer resp.Body.Close()
	assert.Nil(t, err)
	assert.Equal(t, http.StatusServiceUnavailable, resp.StatusCode)

	resp, err = http.Get(testServer.URL + "/alert/rules/1")
	if err != nil {
		return
	}
	defer resp.Body.Close()
	assert.Nil(t, err)
	assert.Equal(t, http.StatusServiceUnavailable, resp.StatusCode)
}

func TestUpdateAlertRule(t *testing.T) {
	m := NewMonitor()
	router := gin.New()
	router.Use(middleware.CORSMiddleware())
	router.PUT("/alert/rules/:id", m.UpdateAlertRule)
	testServer := httptest.NewServer(router)
	assert.NotNil(t, testServer)
	defer testServer.Close()

	rule := models.NewRule()
	rule.Pending = "5s"
	rule.Aggregation = "max"
	rule.Comparison = ">"
	rule.Threshold = 60

	body, err := json.Marshal(rule)
	assert.Nil(t, err, "invalid param")

	req, err := http.NewRequest("PUT", testServer.URL+"/alert/rules/1", strings.NewReader(string(body)))
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")
	httpClient := http.DefaultClient
	resp, err := httpClient.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestUpdateAlertRuleError(t *testing.T) {
	m := NewMonitor()
	router := gin.New()
	router.Use(middleware.CORSMiddleware())
	router.PUT("/alert/rules/:id", m.UpdateAlertRule)
	testServer := httptest.NewServer(router)
	assert.NotNil(t, testServer)
	defer testServer.Close()

	rule := models.NewRule()
	rule.Pending = "5s"
	rule.Aggregation = "max"
	rule.Comparison = ">"
	rule.Threshold = 60

	body, err := json.Marshal(rule)
	assert.Nil(t, err, "invalid param")

	req, err := http.NewRequest("PUT", testServer.URL+"/alert/rules/1", strings.NewReader(""))
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")
	httpClient := http.DefaultClient
	resp, err := httpClient.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	assert.Nil(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	req, err = http.NewRequest("PUT", testServer.URL+"/alert/rules/abc", strings.NewReader(string(body)))
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")
	httpClient = http.DefaultClient
	resp, err = httpClient.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	assert.Nil(t, err)
	assert.Equal(t, http.StatusServiceUnavailable, resp.StatusCode)
}

func TestReceiveAlertEvent(t *testing.T) {
	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()
	mockAlerter := mock_alerter.NewMockAlerter(mockCtl)

	m := NewMonitor()
	m.Alert = mockAlerter

	router := gin.New()
	router.Use(middleware.CORSMiddleware())
	router.POST("/alert/receiver", m.ReceiveAlertEvent)
	testServer := httptest.NewServer(router)
	assert.NotNil(t, testServer)
	defer testServer.Close()

	body := []byte(`{"receiver":"mola_webhook","status":"firing","alerts":[{"status":"firing","labels":{"alertname":"web_zdou_datamanmesos_cpu_usage_warning","app":"web-zdou-datamanmesos","container_label_DM_APP_ID":"web-zdou-datamanmesos","container_label_DM_APP_NAME":"web","container_label_DM_CLUSTER":"datamanmesos","container_label_DM_GROUP_NAME":"dev","container_label_DM_SLOT_ID":"0-web-zdou-datamanmesos","container_label_DM_SLOT_INDEX":"0","container_label_DM_TASK_ID":"0-web-zdou-datamanmesos-c19f31f345cd4110924974df54dac683","container_label_DM_USER":"zdou","container_label_DM_USER_NAME":"zdou","container_label_DM_VCLUSTER":"mola","duration":"5m","group":"dev","id":"/docker/4828b086a265a03ce4b34ac3662c6d67d5cb1a5230e1eb77a253b3fd382c1e19","image":"nginx:1.10","indicator":"cpu_usage","instance":"192.168.56.103:5014","job":"cadvisor","judgement":"max \u003e 60%","name":"mesos-3fca175f-7593-46ad-bcc6-2875e089f8b5-S0.c5a378af-7ea7-4cc7-9f80-de748cf0bbe9","severity":"warning","value":"0.01727926502595149"},"annotations":{"description":"","summary":""},"startsAt":"2017-05-06T20:42:27.804+08:00","endsAt":"0001-01-01T00:00:00Z","generatorURL":"http://srymaster1:9090/graph#%5B%7B%22expr%22%3A%22max%28irate%28container_cpu_usage_seconds_total%7Bcontainer_label_DM_APP_ID%3D%5C%22web-zdou-datamanmesos%5C%22%2Ccontainer_label_DM_LOG_TAG%21%3D%5C%22ignore%5C%22%2Cid%3D~%5C%22%2Fdocker%2F.%2A%5C%22%2Cname%3D~%5C%22mesos.%2A%5C%22%7D%5B5m%5D%29%29%20BY%20%28container_label_DM_APP_ID%29%20KEEP_COMMON%20%2A%20100%20%3C%2060%22%2C%22tab%22%3A0%7D%5D"}],"groupLabels":{"alertname":"web_zdou_datamanmesos_cpu_usage_warning"},"commonLabels":{"alertname":"web_zdou_datamanmesos_cpu_usage_warning","app":"web-zdou-datamanmesos","container_label_DM_APP_ID":"web-zdou-datamanmesos","container_label_DM_APP_NAME":"web","container_label_DM_CLUSTER":"datamanmesos","container_label_DM_GROUP_NAME":"dev","container_label_DM_SLOT_ID":"0-web-zdou-datamanmesos","container_label_DM_SLOT_INDEX":"0","container_label_DM_TASK_ID":"0-web-zdou-datamanmesos-c19f31f345cd4110924974df54dac683","container_label_DM_USER":"zdou","container_label_DM_USER_NAME":"zdou","container_label_DM_VCLUSTER":"mola","duration":"5m","group":"dev","id":"/docker/4828b086a265a03ce4b34ac3662c6d67d5cb1a5230e1eb77a253b3fd382c1e19","image":"nginx:1.10","indicator":"cpu_usage","instance":"192.168.56.103:5014","job":"cadvisor","judgement":"max \u003e 60%","name":"mesos-3fca175f-7593-46ad-bcc6-2875e089f8b5-S0.c5a378af-7ea7-4cc7-9f80-de748cf0bbe9","severity":"warning","value":"0.01727926502595149"},"commonAnnotations":{"description":"","summary":""},"externalURL":"http://srymaster1:9093","version":"3","groupKey":9849835813238314749}`)

	req, err := http.NewRequest("POST", testServer.URL+"/alert/receiver", strings.NewReader(string(body)))
	if err != nil {
		return
	}
	mockAlerter.EXPECT().ReceiveAlertEvent(gomock.Any()).Return(nil).Times(1)
	//mockAlerter.EXPECT().GetIndicatorName(gomock.Any()).Return("", "", errors.New("err")).Times(1)

	req.Header.Set("Content-Type", "application/json")
	httpClient := http.DefaultClient
	resp, err := httpClient.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestReceiveAlertEventError(t *testing.T) {
	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()
	mockAlerter := mock_alerter.NewMockAlerter(mockCtl)

	m := NewMonitor()
	m.Alert = mockAlerter

	router := gin.New()
	router.Use(middleware.CORSMiddleware())
	router.POST("/alert/receiver", m.ReceiveAlertEvent)
	testServer := httptest.NewServer(router)
	assert.NotNil(t, testServer)
	defer testServer.Close()

	body := []byte(`{"receiver":"mola_webhook","status":"firing","alerts":[{"status":"firing","labels":{"alertname":"web_zdou_datamanmesos_cpu_usage_warning","app":"web-zdou-datamanmesos","container_label_DM_APP_ID":"web-zdou-datamanmesos","container_label_DM_APP_NAME":"web","container_label_DM_CLUSTER":"datamanmesos","container_label_DM_GROUP_NAME":"dev","container_label_DM_SLOT_ID":"0-web-zdou-datamanmesos","container_label_DM_SLOT_INDEX":"0","container_label_DM_TASK_ID":"0-web-zdou-datamanmesos-c19f31f345cd4110924974df54dac683","container_label_DM_USER":"zdou","container_label_DM_USER_NAME":"zdou","container_label_DM_VCLUSTER":"mola","duration":"5m","group":"dev","id":"/docker/4828b086a265a03ce4b34ac3662c6d67d5cb1a5230e1eb77a253b3fd382c1e19","image":"nginx:1.10","indicator":"cpu_usage","instance":"192.168.56.103:5014","job":"cadvisor","judgement":"max \u003e 60%","name":"mesos-3fca175f-7593-46ad-bcc6-2875e089f8b5-S0.c5a378af-7ea7-4cc7-9f80-de748cf0bbe9","severity":"warning","value":"0.01727926502595149"},"annotations":{"description":"","summary":""},"startsAt":"2017-05-06T20:42:27.804+08:00","endsAt":"0001-01-01T00:00:00Z","generatorURL":"http://srymaster1:9090/graph#%5B%7B%22expr%22%3A%22max%28irate%28container_cpu_usage_seconds_total%7Bcontainer_label_DM_APP_ID%3D%5C%22web-zdou-datamanmesos%5C%22%2Ccontainer_label_DM_LOG_TAG%21%3D%5C%22ignore%5C%22%2Cid%3D~%5C%22%2Fdocker%2F.%2A%5C%22%2Cname%3D~%5C%22mesos.%2A%5C%22%7D%5B5m%5D%29%29%20BY%20%28container_label_DM_APP_ID%29%20KEEP_COMMON%20%2A%20100%20%3C%2060%22%2C%22tab%22%3A0%7D%5D"}],"groupLabels":{"alertname":"web_zdou_datamanmesos_cpu_usage_warning"},"commonLabels":{"alertname":"web_zdou_datamanmesos_cpu_usage_warning","app":"web-zdou-datamanmesos","container_label_DM_APP_ID":"web-zdou-datamanmesos","container_label_DM_APP_NAME":"web","container_label_DM_CLUSTER":"datamanmesos","container_label_DM_GROUP_NAME":"dev","container_label_DM_SLOT_ID":"0-web-zdou-datamanmesos","container_label_DM_SLOT_INDEX":"0","container_label_DM_TASK_ID":"0-web-zdou-datamanmesos-c19f31f345cd4110924974df54dac683","container_label_DM_USER":"zdou","container_label_DM_USER_NAME":"zdou","container_label_DM_VCLUSTER":"mola","duration":"5m","group":"dev","id":"/docker/4828b086a265a03ce4b34ac3662c6d67d5cb1a5230e1eb77a253b3fd382c1e19","image":"nginx:1.10","indicator":"cpu_usage","instance":"192.168.56.103:5014","job":"cadvisor","judgement":"max \u003e 60%","name":"mesos-3fca175f-7593-46ad-bcc6-2875e089f8b5-S0.c5a378af-7ea7-4cc7-9f80-de748cf0bbe9","severity":"warning","value":"0.01727926502595149"},"commonAnnotations":{"description":"","summary":""},"externalURL":"http://srymaster1:9093","version":"3","groupKey":9849835813238314749}`)

	req, err := http.NewRequest("POST", testServer.URL+"/alert/receiver", strings.NewReader("error"))
	if err != nil {
		return
	}
	//mockAlerter.EXPECT().GetIndicatorName(gomock.Any()).Return("", "", errors.New("err")).Times(1)

	req.Header.Set("Content-Type", "application/json")
	httpClient := http.DefaultClient
	resp, err := httpClient.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	assert.Nil(t, err)
	assert.Equal(t, http.StatusServiceUnavailable, resp.StatusCode)

	req, err = http.NewRequest("POST", testServer.URL+"/alert/receiver", strings.NewReader(string(body)))
	if err != nil {
		return
	}
	mockAlerter.EXPECT().ReceiveAlertEvent(gomock.Any()).Return(errors.New("Error")).Times(1)

	req.Header.Set("Content-Type", "application/json")
	httpClient = http.DefaultClient
	resp, err = httpClient.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	assert.Nil(t, err)
	assert.Equal(t, http.StatusServiceUnavailable, resp.StatusCode)
}

func TestGetAlertEvents(t *testing.T) {
	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()
	mockAlerter := mock_alerter.NewMockAlerter(mockCtl)

	m := NewMonitor()
	m.Alert = mockAlerter

	result := make(map[string]interface{})
	result["count"] = 1
	result["events"] = make([]interface{}, 0)

	router := gin.New()
	router.Use(middleware.CORSMiddleware())
	router.GET("/alert/events", m.GetAlertEvents)
	testServer := httptest.NewServer(router)
	assert.NotNil(t, testServer)
	defer testServer.Close()

	mockAlerter.EXPECT().GetAlertEvents(gomock.Any(), gomock.Any()).Return(result, nil).Times(1)

	resp, err := http.Get(testServer.URL + "/alert/events?ack=false&group=dev&app=app1")
	if err != nil {
		return
	}
	defer resp.Body.Close()
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestGetAlertEventsError(t *testing.T) {
	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()
	mockAlerter := mock_alerter.NewMockAlerter(mockCtl)

	m := NewMonitor()
	m.Alert = mockAlerter

	result := make(map[string]interface{})
	result["count"] = 1
	result["events"] = make([]interface{}, 0)

	router := gin.New()
	router.Use(middleware.CORSMiddleware())
	router.GET("/alert/events", m.GetAlertEvents)
	testServer := httptest.NewServer(router)
	assert.NotNil(t, testServer)
	defer testServer.Close()

	mockAlerter.EXPECT().GetAlertEvents(gomock.Any(), gomock.Any()).Return(result, errors.New("err")).Times(1)

	resp, err := http.Get(testServer.URL + "/alert/events?group=dev&app=app1&start=1234567&end=1234567")
	if err != nil {
		return
	}
	defer resp.Body.Close()
	assert.Nil(t, err)
	assert.Equal(t, http.StatusServiceUnavailable, resp.StatusCode)

	resp, err = http.Get(testServer.URL + "/alert/events?ack=123")
	if err != nil {
		return
	}
	defer resp.Body.Close()
	assert.Nil(t, err)
	assert.Equal(t, http.StatusServiceUnavailable, resp.StatusCode)
}

func TestAckAlertEvent(t *testing.T) {
	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()
	mockAlerter := mock_alerter.NewMockAlerter(mockCtl)

	m := NewMonitor()
	m.Alert = mockAlerter

	result := make(map[string]interface{})
	result["count"] = 1
	result["events"] = make([]interface{}, 0)

	router := gin.New()
	router.Use(middleware.CORSMiddleware())
	router.PUT("/alert/events/:id", m.AckAlertEvent)
	testServer := httptest.NewServer(router)
	assert.NotNil(t, testServer)
	defer testServer.Close()

	body := []byte(`{
  "action":"ack",
	"group": "dev",
	"app": "web-zdou-datamanmesos"}`)
	req, err := http.NewRequest("PUT", testServer.URL+"/alert/events/1", strings.NewReader(string(body)))
	if err != nil {
		return
	}
	mockAlerter.EXPECT().AckAlertEvent(gomock.Any(), gomock.Any()).Return(nil).Times(1)

	req.Header.Set("Content-Type", "application/json")
	httpClient := http.DefaultClient
	resp, err := httpClient.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestAckAlertEventError(t *testing.T) {
	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()
	mockAlerter := mock_alerter.NewMockAlerter(mockCtl)

	m := NewMonitor()
	m.Alert = mockAlerter

	result := make(map[string]interface{})
	result["count"] = 1
	result["events"] = make([]interface{}, 0)

	router := gin.New()
	router.Use(middleware.CORSMiddleware())
	router.PUT("/alert/events/:id", m.AckAlertEvent)
	testServer := httptest.NewServer(router)
	assert.NotNil(t, testServer)
	defer testServer.Close()

	body := []byte(`{
  "action":"ack",
	"group": "dev",
	"app": "web-zdou-datamanmesos"}`)
	req, err := http.NewRequest("PUT", testServer.URL+"/alert/events/abc", strings.NewReader(string(body)))
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")
	httpClient := http.DefaultClient
	resp, err := httpClient.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	assert.Nil(t, err)
	assert.Equal(t, http.StatusServiceUnavailable, resp.StatusCode)

	req, err = http.NewRequest("PUT", testServer.URL+"/alert/events/1", strings.NewReader("err"))
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")
	httpClient = http.DefaultClient
	resp, err = httpClient.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	assert.Nil(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	req, err = http.NewRequest("PUT", testServer.URL+"/alert/events/1", strings.NewReader(string(body)))
	if err != nil {
		return
	}
	mockAlerter.EXPECT().AckAlertEvent(gomock.Any(), gomock.Any()).Return(errors.New("err")).Times(1)

	req.Header.Set("Content-Type", "application/json")
	httpClient = http.DefaultClient
	resp, err = httpClient.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	assert.Nil(t, err)
	assert.Equal(t, http.StatusServiceUnavailable, resp.StatusCode)
}

/*

func TestUpdateAlertRule(t *testing.T) {
	router := gin.New()
	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()
	mockStore := mock_store.NewMockStore(mockCtl)
	alert := NewAlert()
	alert.Store = mockStore
	alert.RulesPath = os.TempDir()
	os.Open(alert.RulesPath)
	defer os.Remove(alert.RulesPath)

	router.Use(middleware.CORSMiddleware())
	router.PUT("/alert/rules", alert.UpdateAlertRule)
	router.POST("/-/reload", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "")
	})
	testServer := httptest.NewServer(router)
	assert.NotNil(t, testServer)
	defer testServer.Close()
	alert.PromServer = testServer.URL

	var rule = &models.Rule{
		Name: "user1",
	}
	var result = models.Rule{
		ID:    1,
		Name:  "user1",
		Alert: "alert",
	}
	body, err := json.Marshal(rule)
	assert.Nil(t, err, "invalid param")

	mockStore.EXPECT().UpdateAlertRule(gomock.Any()).Return(nil).Times(1)
	mockStore.EXPECT().GetAlertRuleByName(gomock.Any(), gomock.Any()).Return(result, nil).Times(1)
	req, err := http.NewRequest("PUT", testServer.URL+"/alert/rules", strings.NewReader(string(body)))
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := alert.HTTPClient.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Test rule form error
	req, err = http.NewRequest("PUT", testServer.URL+"/alert/rules", strings.NewReader("err"))
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err = alert.HTTPClient.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	assert.Nil(t, err)
	assert.Equal(t, 400, resp.StatusCode)

	// test Store.UpdateAlertRule error
	mockStore.EXPECT().UpdateAlertRule(gomock.Any()).Return(errors.New("err")).Times(1)
	req, err = http.NewRequest("PUT", testServer.URL+"/alert/rules", strings.NewReader(string(body)))
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err = alert.HTTPClient.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	assert.Nil(t, err)
	assert.Equal(t, 503, resp.StatusCode)

	// test Store.GetAlertRuleByName error
	mockStore.EXPECT().UpdateAlertRule(gomock.Any()).Return(nil).Times(1)
	mockStore.EXPECT().GetAlertRuleByName(gomock.Any(), gomock.Any()).Return(result, errors.New("err")).Times(1)
	req, err = http.NewRequest("PUT", testServer.URL+"/alert/rules", strings.NewReader(string(body)))
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err = alert.HTTPClient.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	assert.Nil(t, err)
	assert.Equal(t, 503, resp.StatusCode)

	// test WriteAlertFile error
	alert.RulesPath = "test"
	mockStore.EXPECT().UpdateAlertRule(gomock.Any()).Return(nil).Times(1)
	mockStore.EXPECT().GetAlertRuleByName(gomock.Any(), gomock.Any()).Return(result, nil).Times(1)
	req, err = http.NewRequest("PUT", testServer.URL+"/alert/rules", strings.NewReader(string(body)))
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err = alert.HTTPClient.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	assert.Nil(t, err)
	assert.Equal(t, 503, resp.StatusCode)
}

func TestListAlertRules(t *testing.T) {
	router := gin.New()
	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()
	mockStore := mock_store.NewMockStore(mockCtl)

	alert := NewAlert()
	alert.Store = mockStore
	router.Use(middleware.CORSMiddleware())
	router.GET("/alert/rules", alert.ListAlertRules)
	testServer := httptest.NewServer(router)
	assert.NotNil(t, testServer)
	defer testServer.Close()

	var rule = &models.Rule{
		ID:    1,
		Name:  "user1",
		Alert: "alert",
	}
	var rules []*models.Rule
	rules = append(rules, rule)
	data := map[string]interface{}{"count": 1, "rules": rules}
	mockStore.EXPECT().ListAlertRules(gomock.Any(), gomock.Any()).Return(data, nil).Times(1)
	resp, err := http.Get(testServer.URL + "/alert/rules")
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// test error
	mockStore.EXPECT().ListAlertRules(gomock.Any(), gomock.Any()).Return(data, errors.New("test")).Times(1)
	resp, err = http.Get(testServer.URL + "/alert/rules")
	assert.Nil(t, err)
	assert.Equal(t, http.StatusServiceUnavailable, resp.StatusCode)
}

func TestGetAlertRule(t *testing.T) {
	router := gin.New()
	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()
	mockStore := mock_store.NewMockStore(mockCtl)

	alert := NewAlert()
	alert.Store = mockStore
	router.GET("/alert/rules/:id", alert.GetAlertRule)
	testServer := httptest.NewServer(router)
	assert.NotNil(t, testServer)
	defer testServer.Close()

	var rule = models.Rule{
		ID:    1,
		Name:  "user1",
		Alert: "alert",
	}

	mockStore.EXPECT().GetAlertRule(gomock.Any()).Return(rule, nil).Times(1)
	resp, err := http.Get(testServer.URL + "/alert/rules/1")
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// test error id
	resp, err = http.Get(testServer.URL + "/alert/rules/abc")
	assert.Nil(t, err)
	assert.Equal(t, http.StatusServiceUnavailable, resp.StatusCode)

	//test error with get rule error
	mockStore.EXPECT().GetAlertRule(gomock.Any()).Return(rule, errors.New("test")).Times(1)
	resp, err = http.Get(testServer.URL + "/alert/rules/1")
	assert.Nil(t, err)
	assert.Equal(t, http.StatusServiceUnavailable, resp.StatusCode)
}

func TestReloadAlertRuleConf(t *testing.T) {
	router := gin.New()
	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()
	mockStore := mock_store.NewMockStore(mockCtl)
	alert := NewAlert()
	alert.Store = mockStore

	router.Use(middleware.CORSMiddleware())
	router.POST("/alert/rules/conf", alert.ReloadAlertRuleConf)
	router.POST("/-/reload", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "")
	})
	testServer := httptest.NewServer(router)
	assert.NotNil(t, testServer)
	defer testServer.Close()
	alert.PromServer = testServer.URL

	req, err := http.NewRequest("POST", testServer.URL+"/alert/rules/conf", nil)
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := alert.HTTPClient.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	assert.Nil(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}

func TestReloadAlertRuleConfError(t *testing.T) {
	router := gin.New()
	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()
	mockStore := mock_store.NewMockStore(mockCtl)
	alert := NewAlert()
	alert.Store = mockStore

	router.Use(middleware.CORSMiddleware())
	router.POST("/alert/rules/conf", alert.ReloadAlertRuleConf)
	router.POST("/-/reload", func(ctx *gin.Context) {
		ctx.String(http.StatusForbidden, "")
	})
	testServer := httptest.NewServer(router)
	assert.NotNil(t, testServer)
	defer testServer.Close()
	alert.PromServer = testServer.URL

	req, err := http.NewRequest("POST", testServer.URL+"/alert/rules/conf", nil)
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := alert.HTTPClient.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	assert.Nil(t, err)
	assert.Equal(t, 503, resp.StatusCode)
}

func TestGetAlertEventsAck(t *testing.T) {
	router := gin.New()
	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()
	mockStore := mock_store.NewMockStore(mockCtl)
	alert := NewAlert()
	alert.Store = mockStore

	router.Use(middleware.CORSMiddleware())
	router.GET("/alert/events", alert.GetAlertEvents)
	testServer := httptest.NewServer(router)
	assert.NotNil(t, testServer)
	defer testServer.Close()

	var event = &models.Event{
		Count:       1,
		Severity:    "warning",
		VCluster:    "cluster1",
		App:         "app1",
		Slot:        "0",
		UserName:    "user1",
		GroupName:   "group1",
		ContainerID: "container1",
		AlertName:   "alert1",
		Ack:         true,
	}
	var result []*models.Event
	result = append(result, event)
	var data = map[string]interface{}{
		"count":  1,
		"events": result,
	}
	mockStore.EXPECT().ListAckedEvent(gomock.Any(), gomock.Any(), gomock.Any()).Return(data).Times(1)
	u, _ := url.Parse(testServer.URL)
	u.Path = strings.TrimRight(u.Path, "/") + "/alert/events"
	q := u.Query()
	q.Set("ack", "true")
	q.Set("user_name", "user1")
	q.Set("group_name", "group1")
	u.RawQuery = q.Encode()
	resp, err := alert.HTTPClient.Get(u.String())
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestGetAlertEventsUnack(t *testing.T) {
	router := gin.New()
	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()
	mockStore := mock_store.NewMockStore(mockCtl)
	alert := NewAlert()
	alert.Store = mockStore

	router.Use(middleware.CORSMiddleware())
	router.GET("/alert/events", alert.GetAlertEvents)
	testServer := httptest.NewServer(router)
	assert.NotNil(t, testServer)
	defer testServer.Close()

	var event = &models.Event{
		Count:       1,
		Severity:    "warning",
		VCluster:    "cluster1",
		App:         "app1",
		Slot:        "0",
		UserName:    "user1",
		GroupName:   "group1",
		ContainerID: "container1",
		AlertName:   "alert1",
		Ack:         false,
	}
	var result []*models.Event
	result = append(result, event)
	var data = map[string]interface{}{
		"count":  1,
		"events": result,
	}
	mockStore.EXPECT().ListUnackedEvent(gomock.Any(), gomock.Any(), gomock.Any()).Return(data).Times(1)
	u, _ := url.Parse(testServer.URL)
	u.Path = strings.TrimRight(u.Path, "/") + "/alert/events"
	q := u.Query()
	q.Set("ack", "false")
	q.Set("user_name", "user1")
	q.Set("group_name", "group1")
	u.RawQuery = q.Encode()
	resp, err := alert.HTTPClient.Get(u.String())
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestAckAlertEvent(t *testing.T) {
	router := gin.New()
	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()
	mockStore := mock_store.NewMockStore(mockCtl)
	alert := NewAlert()
	alert.Store = mockStore

	router.Use(middleware.CORSMiddleware())
	router.PUT("/alert/events/:id", alert.AckAlertEvent)
	testServer := httptest.NewServer(router)
	assert.NotNil(t, testServer)
	defer testServer.Close()

	var event = map[string]interface{}{
		"action": "ack",
	}
	body, err := json.Marshal(event)
	assert.Nil(t, err, "invalid param")

	mockStore.EXPECT().AckEvent(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).Times(1)
	req, err := http.NewRequest("PUT", testServer.URL+"/alert/events/1", strings.NewReader(string(body)))
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := alert.HTTPClient.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	req, err = http.NewRequest("PUT", testServer.URL+"/alert/events/abc", strings.NewReader(string(body)))
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err = alert.HTTPClient.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	assert.Nil(t, err)
	assert.Equal(t, 503, resp.StatusCode)

	req, err = http.NewRequest("PUT", testServer.URL+"/alert/events/1", strings.NewReader("err"))
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err = alert.HTTPClient.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	assert.Nil(t, err)
	assert.Equal(t, 400, resp.StatusCode)

	event = map[string]interface{}{
		"action":     "ack",
		"user_name":  "user1",
		"group_name": "group1",
	}
	body, err = json.Marshal(event)
	assert.Nil(t, err, "invalid param")
	mockStore.EXPECT().AckEvent(gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("err")).Times(1)
	req, err = http.NewRequest("PUT", testServer.URL+"/alert/events/1", strings.NewReader(string(body)))
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err = alert.HTTPClient.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	assert.Nil(t, err)
	assert.Equal(t, 503, resp.StatusCode)
}

func TestReceiveAlertEvent(t *testing.T) {
	router := gin.New()
	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()
	mockStore := mock_store.NewMockStore(mockCtl)
	alert := NewAlert()
	alert.Store = mockStore

	router.Use(middleware.CORSMiddleware())
	router.POST("/receiver", alert.ReceiveAlertEvent)
	testServer := httptest.NewServer(router)
	assert.NotNil(t, testServer)
	defer testServer.Close()

	labels := map[string]interface{}{
		"alertname":                  "alert1",
		"severity":                   "Warning",
		"container_label_VCLUSTER":   "wtzhou-VCluster",
		"container_label_APP":        "app-3",
		"container_label_SLOT":       "slot-2",
		"container_label_GROUP_NAME": "group-1",
		"container_label_USER_NAME":  "wtzhou",
		"id": "/docker/aaaxefgh32e2e23rfsda",
	}
	annotations := map[string]interface{}{
		"description": "High Mem usage on instance: test-1",
		"summary":     "Mem Usage on instance: test-1",
	}
	event := map[string]interface{}{
		"labels":      labels,
		"annotations": annotations,
	}
	var events []map[string]interface{}
	events = append(events, event)

	var data = map[string]interface{}{
		"alerts": events,
	}
	body, err := json.Marshal(data)
	assert.Nil(t, err, "invalid param")
	mockStore.EXPECT().CreateOrIncreaseEvent(gomock.Any()).Return(nil).Times(1)

	req, err := http.NewRequest("POST", testServer.URL+"/receiver", strings.NewReader(string(body)))
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := alert.HTTPClient.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	req, err = http.NewRequest("POST", testServer.URL+"/receiver", strings.NewReader("err"))
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err = alert.HTTPClient.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	assert.Nil(t, err)
	assert.Equal(t, 503, resp.StatusCode)

	mockStore.EXPECT().CreateOrIncreaseEvent(gomock.Any()).Return(errors.New("err")).Times(1)
	req, err = http.NewRequest("POST", testServer.URL+"/receiver", strings.NewReader(string(body)))
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err = alert.HTTPClient.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	assert.Nil(t, err)
	assert.Equal(t, 503, resp.StatusCode)
}

func TestUpdateAlertRuleFiles(t *testing.T) {
	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()
	mockStore := mock_store.NewMockStore(mockCtl)
	alert := NewAlert()
	alert.Store = mockStore
	alert.RulesPath = os.TempDir() + "rules"
	os.Mkdir(alert.RulesPath, 0755)
	defer os.Remove(alert.RulesPath)

	var rules []*models.Rule

	var rule = &models.Rule{
		ID:       1,
		Name:     "user1",
		Alert:    "alert",
		Expr:     "expr1",
		Duration: "duration",
		Labels:   "labels",
	}
	rule.Description = "desciption"
	rule.Summary = "summary"
	rules = append(rules, rule)
	alert.WriteAlertFile(rule)

	mockStore.EXPECT().GetAlertRules().Return(rules, nil).Times(1)
	alert.UpdateAlertRuleFiles()
}

func TestUpdateAlertRuleFilesCreate(t *testing.T) {
	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()
	mockStore := mock_store.NewMockStore(mockCtl)
	alert := NewAlert()
	alert.Store = mockStore
	alert.RulesPath = os.TempDir() + "rules"
	os.Mkdir(alert.RulesPath, 0755)
	defer os.Remove(alert.RulesPath)

	var rules []*models.Rule

	var rule = &models.Rule{
		ID:       1,
		Name:     "user1",
		Alert:    "alert",
		Expr:     "expr1",
		Duration: "duration",
		Labels:   "labels",
	}
	rule.Description = "desciption"
	rule.Summary = "summary"

	var rule2 = &models.Rule{
		ID:       2,
		Name:     "user2",
		Alert:    "alert",
		Expr:     "expr1",
		Duration: "duration",
		Labels:   "labels",
	}
	rule2.Description = "desciption"
	rule2.Summary = "summary"

	rules = append(rules, rule)
	rules = append(rules, rule2)
	alert.WriteAlertFile(rule)

	mockStore.EXPECT().GetAlertRules().Return(rules, nil).Times(1)
	alert.UpdateAlertRuleFiles()
}
*/
