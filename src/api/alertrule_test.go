package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/Dataman-Cloud/log-proxy/src/models"
	"github.com/Dataman-Cloud/log-proxy/src/router/middleware"
	mock_store "github.com/Dataman-Cloud/log-proxy/src/store/mock_datastore"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestNewAlert(t *testing.T) {
	alert := NewAlert()
	assert.Equal(t, "1m", alert.Interval)
}

func TestCreateAlertRule(t *testing.T) {
	router := gin.New()
	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()
	mockStore := mock_store.NewMockStore(mockCtl)
	alert := NewAlert()
	alert.Store = mockStore
	alert.RulesPath = os.TempDir()
	_, err := os.Open(alert.RulesPath)
	defer os.Remove(alert.RulesPath)

	router.Use(middleware.CORSMiddleware())
	router.POST("/alert/rules", alert.CreateAlertRule)
	router.POST("/-/reload", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "")
	})
	testServer := httptest.NewServer(router)
	assert.NotNil(t, testServer)
	defer testServer.Close()
	alert.PromServer = testServer.URL

	var rule = &models.Rule{
		Name:  "user1",
		Alert: "alert",
	}
	var result = models.Rule{
		ID:    1,
		Name:  "user1",
		Alert: "alert",
	}
	body, err := json.Marshal(rule)
	assert.Nil(t, err, "invalid param")

	mockStore.EXPECT().CreateAlertRule(gomock.Any()).Return(nil).Times(1)
	mockStore.EXPECT().GetAlertRuleByName(gomock.Any(), gomock.Any()).Return(result, nil).Times(1)

	req, err := http.NewRequest("POST", testServer.URL+"/alert/rules", strings.NewReader(string(body)))
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

	// test error store return
	mockStore.EXPECT().CreateAlertRule(gomock.Any()).Return(errors.New("err")).Times(1)
	req, err = http.NewRequest("POST", testServer.URL+"/alert/rules", strings.NewReader(string(body)))
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

	// test error rule input
	req, err = http.NewRequest("POST", testServer.URL+"/alert/rules", strings.NewReader("err"))
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

	// test error rule
	rule = &models.Rule{
		Alert: "alert",
	}
	body, err = json.Marshal(rule)
	assert.Nil(t, err, "invalid param")
	req, err = http.NewRequest("POST", testServer.URL+"/alert/rules", strings.NewReader(string(body)))
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

func TestDeleteAlertRule(t *testing.T) {
	router := gin.New()
	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()
	mockStore := mock_store.NewMockStore(mockCtl)
	alert := NewAlert()
	alert.Store = mockStore
	alert.RulesPath = os.TempDir()
	os.Open(alert.RulesPath + "user1.rule")
	defer os.Remove(alert.RulesPath)

	router.Use(middleware.CORSMiddleware())
	router.DELETE("/alert/rules/:id", alert.DeleteAlertRule)
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

	mockStore.EXPECT().GetAlertRule(gomock.Any()).Return(result, nil).Times(1)
	mockStore.EXPECT().DeleteAlertRuleByIDName(gomock.Any(), gomock.Any()).Return(int64(1), nil).Times(1)
	req, err := http.NewRequest("DELETE", testServer.URL+"/alert/rules/1", strings.NewReader(string(body)))
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

	// test rule form error
	req, err = http.NewRequest("DELETE", testServer.URL+"/alert/rules/1", strings.NewReader("err"))
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

	// test error id
	req, err = http.NewRequest("DELETE", testServer.URL+"/alert/rules/abc", strings.NewReader(string(body)))
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

	// test store.GetAlertRule error
	mockStore.EXPECT().GetAlertRule(gomock.Any()).Return(result, errors.New("err")).Times(1)
	req, err = http.NewRequest("DELETE", testServer.URL+"/alert/rules/1", strings.NewReader(string(body)))
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

	// test store.DeleteAlertRuleByIDName error
	mockStore.EXPECT().GetAlertRule(gomock.Any()).Return(result, nil).Times(1)
	mockStore.EXPECT().DeleteAlertRuleByIDName(gomock.Any(), gomock.Any()).Return(int64(1), errors.New("err")).Times(1)
	req, err = http.NewRequest("DELETE", testServer.URL+"/alert/rules/1", strings.NewReader(string(body)))
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

	// test rowsAffected == 0
	mockStore.EXPECT().GetAlertRule(gomock.Any()).Return(result, nil).Times(1)
	mockStore.EXPECT().DeleteAlertRuleByIDName(gomock.Any(), gomock.Any()).Return(int64(0), nil).Times(1)
	req, err = http.NewRequest("DELETE", testServer.URL+"/alert/rules/1", strings.NewReader(string(body)))
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

	// test remove file error
	alert.RulesPath = "file"
	mockStore.EXPECT().GetAlertRule(gomock.Any()).Return(result, nil).Times(1)
	mockStore.EXPECT().DeleteAlertRuleByIDName(gomock.Any(), gomock.Any()).Return(int64(1), nil).Times(1)
	req, err = http.NewRequest("DELETE", testServer.URL+"/alert/rules/1", strings.NewReader(string(body)))
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
