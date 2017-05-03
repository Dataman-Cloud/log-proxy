package api

import (
	"bytes"
	"container/list"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Dataman-Cloud/log-proxy/src/models"
	mock_store "github.com/Dataman-Cloud/log-proxy/src/store/mock_datastore"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestCreateLogAlertRule(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockStore := mock_store.NewMockStore(mockCtrl)
	s := Search{
		Store:         mockStore,
		KeywordFilter: make(map[string]*list.List),
	}

	router := gin.New()
	router.POST("/rules", s.CreateLogAlertRule)
	testServer := httptest.NewServer(router)
	assert.NotNil(t, testServer)
	defer testServer.Close()

	rule := models.LogAlertRule{
		App:     "app",
		Cluster: "cluster",
		Keyword: "key",
		Source:  "stdout",
		User:    "user",
		Group:   "group",
	}

	ruleMetaData, err := json.Marshal(rule)
	assert.Nil(t, err)

	// test success condition
	mockStore.EXPECT().CreateLogAlertRule(gomock.Any()).Return(nil).Times(1)
	resp, err := http.Post(testServer.URL+"/rules", "application/json", bytes.NewReader(ruleMetaData))
	assert.Nil(t, err)
	assert.Equal(t, resp.StatusCode, http.StatusOK)

	// test binJSON error
	resp, err = http.Post(testServer.URL+"/rules", "application/json", bytes.NewReader([]byte("xxxxx")))
	assert.Nil(t, err)
	assert.Equal(t, resp.StatusCode, http.StatusBadRequest)

	// test for duplicate keyword
	resp, err = http.Post(testServer.URL+"/rules", "application/json", bytes.NewReader(ruleMetaData))
	assert.Nil(t, err)
	assert.Equal(t, resp.StatusCode, http.StatusServiceUnavailable)

	// test db return error
	rule.Group = "group1"
	ruleMetaData, err = json.Marshal(rule)
	assert.Nil(t, err)
	mockStore.EXPECT().CreateLogAlertRule(gomock.Any()).Return(errors.New("test")).Times(1)
	resp, err = http.Post(testServer.URL+"/rules", "application/json", bytes.NewReader(ruleMetaData))
	assert.Nil(t, err)
	assert.Equal(t, resp.StatusCode, http.StatusServiceUnavailable)
}

func TestUpdateLogAlertRule(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockStore := mock_store.NewMockStore(mockCtrl)
	s := Search{
		Store:         mockStore,
		KeywordFilter: make(map[string]*list.List),
	}

	router := gin.New()
	router.PUT("/rules/:id", s.UpdateLogAlertRule)
	testServer := httptest.NewServer(router)
	assert.NotNil(t, testServer)
	defer testServer.Close()

	rule := models.LogAlertRule{
		ID:      1,
		App:     "app",
		Cluster: "cluster",
		Keyword: "key",
		Source:  "stdout",
		User:    "user",
		Group:   "group",
	}

	ruleMetaData, err := json.Marshal(rule)
	assert.Nil(t, err)

	// test success condition 1: keyword filter was empty
	mockStore.EXPECT().UpdateLogAlertRule(gomock.Any()).Return(nil).Times(1)
	req, err := http.NewRequest("PUT", testServer.URL+"/rules/1", bytes.NewReader(ruleMetaData))
	assert.Nil(t, err)
	resp, err := http.DefaultClient.Do(req)
	assert.Nil(t, err)
	assert.Equal(t, resp.StatusCode, http.StatusOK)

	// test success condition 2: keyword filter was not empty
	ruleIndex := getLogAlertRuleIndex(rule)
	s.KeywordFilter[ruleIndex] = list.New()
	s.KeywordFilter[ruleIndex].PushBack(rule)
	mockStore.EXPECT().UpdateLogAlertRule(gomock.Any()).Return(nil).Times(1)
	req, err = http.NewRequest("PUT", testServer.URL+"/rules/1", bytes.NewReader(ruleMetaData))
	assert.Nil(t, err)
	resp, err = http.DefaultClient.Do(req)
	assert.Nil(t, err)
	assert.Equal(t, resp.StatusCode, http.StatusOK)

	// test error with bindJSON error
	req, err = http.NewRequest("PUT", testServer.URL+"/rules/1", bytes.NewReader([]byte("xxxxxx")))
	assert.Nil(t, err)
	resp, err = http.DefaultClient.Do(req)
	assert.Nil(t, err)
	assert.Equal(t, resp.StatusCode, http.StatusBadRequest)

	// test error with db update error
	mockStore.EXPECT().UpdateLogAlertRule(gomock.Any()).Return(errors.New("test")).Times(1)
	req, err = http.NewRequest("PUT", testServer.URL+"/rules/1", bytes.NewReader(ruleMetaData))
	assert.Nil(t, err)
	resp, err = http.DefaultClient.Do(req)
	assert.Nil(t, err)
	assert.Equal(t, resp.StatusCode, http.StatusServiceUnavailable)
}

func TestGetLogAlertRule(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockStore := mock_store.NewMockStore(mockCtrl)
	s := Search{Store: mockStore}

	router := gin.New()
	router.GET("/rules/:id", s.GetLogAlertRule)
	testServer := httptest.NewServer(router)
	assert.NotNil(t, testServer)
	defer testServer.Close()

	rule := models.LogAlertRule{
		ID:      1,
		App:     "app",
		Cluster: "cluster",
		Keyword: "key",
		Source:  "stdout",
		User:    "user",
		Group:   "group",
	}

	mockStore.EXPECT().GetLogAlertRule(gomock.Any()).Return(rule, nil).Times(1)
	resp, err := http.Get(testServer.URL + "/rules/test")
	assert.Nil(t, err)
	assert.Equal(t, resp.StatusCode, http.StatusOK)

	mockStore.EXPECT().GetLogAlertRule(gomock.Any()).Return(rule, errors.New("test")).Times(1)
	resp, err = http.Get(testServer.URL + "/rules/test")
	assert.Nil(t, err)
	assert.Equal(t, resp.StatusCode, http.StatusServiceUnavailable)
}

func TestGetLogAlertRules(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockStore := mock_store.NewMockStore(mockCtrl)
	s := Search{Store: mockStore}

	router := gin.New()
	router.GET("/rules", func(ctx *gin.Context) { ctx.Set("page", models.Page{}) }, s.GetLogAlertRules)
	testServer := httptest.NewServer(router)
	assert.NotNil(t, testServer)
	defer testServer.Close()

	rules := map[string]interface{}{"test": 0}
	mockStore.EXPECT().GetLogAlertRules(gomock.Any(), gomock.Any()).Return(rules, nil).Times(1)
	resp, err := http.Get(testServer.URL + "/rules?group=test")
	assert.Nil(t, err)
	assert.Equal(t, resp.StatusCode, http.StatusOK)

	mockStore.EXPECT().GetLogAlertRules(gomock.Any(), gomock.Any()).Return(rules, errors.New("test")).Times(1)
	resp, err = http.Get(testServer.URL + "/rules?group=test")
	assert.Nil(t, err)
	assert.Equal(t, resp.StatusCode, http.StatusServiceUnavailable)
}

func TestDeleteLogAlertRule(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockStore := mock_store.NewMockStore(mockCtrl)
	s := Search{
		Store:         mockStore,
		KeywordFilter: make(map[string]*list.List),
	}

	router := gin.New()
	router.DELETE("/rules/:id", s.DeleteLogAlertRule)
	testServer := httptest.NewServer(router)
	assert.NotNil(t, testServer)
	defer testServer.Close()

	rule := models.LogAlertRule{
		App:     "app",
		Cluster: "cluster",
		Keyword: "key",
		Source:  "stdout",
		User:    "user",
		Group:   "group",
	}

	// test success condition 1: keyword map was empty
	mockStore.EXPECT().GetLogAlertRule(gomock.Any()).Return(rule, nil).Times(1)
	mockStore.EXPECT().DeleteLogAlertRule(gomock.Any()).Return(nil).Times(1)
	req, err := http.NewRequest("DELETE", testServer.URL+"/rules/1", nil)
	assert.Nil(t, err)
	resp, err := http.DefaultClient.Do(req)
	assert.Nil(t, err)
	assert.Equal(t, resp.StatusCode, http.StatusOK)

	// test succsess condition 2: keyword map was not empty
	ruleIndex := getLogAlertRuleIndex(rule)
	s.KeywordFilter[ruleIndex] = list.New()
	s.KeywordFilter[ruleIndex].PushBack(rule)
	mockStore.EXPECT().GetLogAlertRule(gomock.Any()).Return(rule, nil).Times(1)
	mockStore.EXPECT().DeleteLogAlertRule(gomock.Any()).Return(nil).Times(1)
	req, err = http.NewRequest("DELETE", testServer.URL+"/rules/1", nil)
	assert.Nil(t, err)
	resp, err = http.DefaultClient.Do(req)
	assert.Nil(t, err)
	assert.Equal(t, resp.StatusCode, http.StatusOK)

	// test get rule error
	mockStore.EXPECT().GetLogAlertRule(gomock.Any()).Return(rule, errors.New("test")).Times(1)
	req, err = http.NewRequest("DELETE", testServer.URL+"/rules/1", nil)
	assert.Nil(t, err)
	resp, err = http.DefaultClient.Do(req)
	assert.Nil(t, err)
	assert.Equal(t, resp.StatusCode, http.StatusServiceUnavailable)

	// test delete error
	mockStore.EXPECT().GetLogAlertRule(gomock.Any()).Return(rule, nil).Times(1)
	mockStore.EXPECT().DeleteLogAlertRule(gomock.Any()).Return(errors.New("test")).Times(1)
	req, err = http.NewRequest("DELETE", testServer.URL+"/rules/1", nil)
	assert.Nil(t, err)
	resp, err = http.DefaultClient.Do(req)
	assert.Nil(t, err)
	assert.Equal(t, resp.StatusCode, http.StatusServiceUnavailable)
}

func TestGetLogAlertEvents(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockStore := mock_store.NewMockStore(mockCtrl)
	s := Search{Store: mockStore}

	router := gin.New()
	router.GET("/events", func(ctx *gin.Context) { ctx.Set("page", models.Page{}) }, s.GetLogAlertEvents)
	testServer := httptest.NewServer(router)
	assert.NotNil(t, testServer)
	defer testServer.Close()

	results := map[string]interface{}{"test": 1}
	mockStore.EXPECT().GetLogAlertEvents(gomock.Any(), gomock.Any()).Return(results, nil).Times(1)
	resp, err := http.Get(testServer.URL + "/events?group=test&app=test")
	assert.Nil(t, err)
	assert.Equal(t, resp.StatusCode, http.StatusOK)

	mockStore.EXPECT().GetLogAlertEvents(gomock.Any(), gomock.Any()).Return(nil, errors.New("test")).Times(1)
	resp, err = http.Get(testServer.URL + "/events")
	assert.Nil(t, err)
	assert.Equal(t, resp.StatusCode, http.StatusServiceUnavailable)
}

func TestGetLogAlertApps(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockStore := mock_store.NewMockStore(mockCtrl)
	s := Search{Store: mockStore}

	router := gin.New()
	router.GET("/apps", func(ctx *gin.Context) { ctx.Set("page", models.Page{}) }, s.GetLogAlertApps)
	testServer := httptest.NewServer(router)
	assert.NotNil(t, testServer)
	defer testServer.Close()

	mockStore.EXPECT().GetLogAlertApps(gomock.Any(), gomock.Any()).Return([]*models.LogAlertApps{}, nil).Times(1)
	resp, err := http.Get(testServer.URL + "/apps?group=test")
	assert.Nil(t, err)
	assert.Equal(t, resp.StatusCode, http.StatusOK)

	mockStore.EXPECT().GetLogAlertApps(gomock.Any(), gomock.Any()).Return(nil, errors.New("test")).Times(1)
	resp, err = http.Get(testServer.URL + "/apps")
	assert.Nil(t, err)
	assert.Equal(t, resp.StatusCode, http.StatusServiceUnavailable)
}

func TestLogAlertEventAction(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockStore := mock_store.NewMockStore(mockCtrl)
	s := Search{Store: mockStore}

	router := gin.New()
	router.PATCH("/events/:id", s.LogAlertEventAction)
	testServer := httptest.NewServer(router)
	assert.NotNil(t, testServer)
	defer testServer.Close()

	// test success condition
	mockStore.EXPECT().AckLogAlertEvent(gomock.Any()).Return(nil).Times(1)
	req, err := http.NewRequest("PATCH", testServer.URL+"/events/1", bytes.NewReader([]byte(`{"action":"ack"}`)))
	assert.Nil(t, err)
	resp, err := http.DefaultClient.Do(req)
	assert.Nil(t, err)
	assert.Equal(t, resp.StatusCode, http.StatusOK)

	// test db error
	mockStore.EXPECT().AckLogAlertEvent(gomock.Any()).Return(errors.New("test")).Times(1)
	req, err = http.NewRequest("PATCH", testServer.URL+"/events/1", bytes.NewReader([]byte(`{"action":"ack"}`)))
	assert.Nil(t, err)
	resp, err = http.DefaultClient.Do(req)
	assert.Nil(t, err)
	assert.Equal(t, resp.StatusCode, http.StatusBadRequest)

	// test invalid body param
	req, err = http.NewRequest("PATCH", testServer.URL+"/events/1", bytes.NewReader([]byte("ack")))
	assert.Nil(t, err)
	resp, err = http.DefaultClient.Do(req)
	assert.Nil(t, err)
	assert.Equal(t, resp.StatusCode, http.StatusBadRequest)

	// test invalid action
	req, err = http.NewRequest("PATCH", testServer.URL+"/events/1", bytes.NewReader([]byte(`{"action":"test"}`)))
	assert.Nil(t, err)
	resp, err = http.DefaultClient.Do(req)
	assert.Nil(t, err)
	assert.Equal(t, resp.StatusCode, http.StatusBadRequest)

}
