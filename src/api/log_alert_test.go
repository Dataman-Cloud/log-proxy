package api

import (
	"bytes"
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
	s := Search{Store: mockStore}

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

	mockStore.EXPECT().CreateLogAlertRule(gomock.Any()).Return(nil).Times(1)
	resp, err := http.Post(testServer.URL+"/rules", "application/json", bytes.NewReader(ruleMetaData))
	assert.Nil(t, err)
	assert.Equal(t, resp.StatusCode, http.StatusOK)

	resp, err = http.Post(testServer.URL+"/rules", "application/json", bytes.NewReader([]byte("xxxxx")))
	assert.Nil(t, err)
	assert.Equal(t, resp.StatusCode, http.StatusBadRequest)

	mockStore.EXPECT().CreateLogAlertRule(gomock.Any()).Return(errors.New("test")).Times(1)
	resp, err = http.Post(testServer.URL+"/rules", "application/json", bytes.NewReader(ruleMetaData))
	assert.Nil(t, err)
	assert.Equal(t, resp.StatusCode, http.StatusServiceUnavailable)
}

func TestUpdateLogAlertRule(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockStore := mock_store.NewMockStore(mockCtrl)
	s := Search{Store: mockStore}

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

	mockStore.EXPECT().UpdateLogAlertRule(gomock.Any()).Return(nil).Times(1)
	req, err := http.NewRequest("PUT", testServer.URL+"/rules/1", bytes.NewReader(ruleMetaData))
	assert.Nil(t, err)
	resp, err := http.DefaultClient.Do(req)
	assert.Nil(t, err)
	assert.Equal(t, resp.StatusCode, http.StatusOK)

	req, err = http.NewRequest("PUT", testServer.URL+"/rules/1", bytes.NewReader([]byte("xxxxxx")))
	assert.Nil(t, err)
	resp, err = http.DefaultClient.Do(req)
	assert.Nil(t, err)
	assert.Equal(t, resp.StatusCode, http.StatusBadRequest)

	mockStore.EXPECT().UpdateLogAlertRule(gomock.Any()).Return(errors.New("test")).Times(1)
	req, err = http.NewRequest("PUT", testServer.URL+"/rules/1", bytes.NewReader(ruleMetaData))
	assert.Nil(t, err)
	resp, err = http.DefaultClient.Do(req)
	assert.Nil(t, err)
	assert.Equal(t, resp.StatusCode, http.StatusServiceUnavailable)
}
