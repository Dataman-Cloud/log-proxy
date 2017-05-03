package api

import (
	"bytes"
	"container/list"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Dataman-Cloud/log-proxy/src/models"
	mock_store "github.com/Dataman-Cloud/log-proxy/src/store/mock_datastore"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/golang/mock/gomock"
)

func TestHandleLogEvent(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockStore := mock_store.NewMockStore(mockCtrl)
	s := Search{
		Store:         mockStore,
		KeywordFilter: make(map[string]*list.List),
	}

	event := models.LogAlertEvent{
		Group:   "group",
		App:     "app",
		Source:  "source",
		Message: "test message for event",
	}
	s.HandleLogEvent(event)

	rule := models.LogAlertRule{
		Group:   "group",
		App:     "app",
		Source:  "source",
		Keyword: "test",
	}
	mockStore.EXPECT().CreateLogAlertEvent(gomock.Any()).Return(nil).Times(1)
	ruleIndex := getLogAlertRuleIndex(rule)
	s.KeywordFilter[ruleIndex] = list.New()
	s.KeywordFilter[ruleIndex].PushBack(rule)
	s.HandleLogEvent(event)

	// create alert event db error
	mockStore.EXPECT().CreateLogAlertEvent(gomock.Any()).Return(errors.New("test")).Times(1)
	s.HandleLogEvent(event)

	// no match event
	event2 := models.LogAlertEvent{
		Group:   "group",
		App:     "app",
		Source:  "source",
		Message: "",
	}
	s.HandleLogEvent(event2)

	// disable rule
	rule.Status = LogAlertRuleDisabled
	s.KeywordFilter[ruleIndex] = list.New()
	s.KeywordFilter[ruleIndex].PushBack(rule)
	s.HandleLogEvent(event)
}

var eventExample = `
{
    "DM_APP_NAME":"aaaaaaaaaaa",
    "DM_SLOT_INDEX":"0",
    "path":"stdout",
    "message":"192.168.1.146",
    "DM_TASK_ID":"0-aaaaaaaaaaa-yaoyun-datamanmesos-1a4cde96aef34c03a8ed9f07dae5ec63",
    "DM_VCLUSTER":"ymola",
    "DM_APP_ID":"aaaaaaaaaaa-yaoyun-datamanmesos",
    "logtime":"2017-05-03T02:25:46.891373648Z",
    "@version":"1",
    "@timestamp":"2017-05-03T02:25:46.892Z",
    "host":"192.168.1.102",
    "DM_CLUSTER":"datamanmesos",
    "DM_GROUP_NAME":"yaoyun",
    "DM_SLOT_ID":"0-aaaaaaaaaaa-yaoyun-datamanmesos",
    "containerid":"2101477298d0622d2fb81156ec6ad328c778035dfc528c475f8e1d6908ece819",
    "DM_USER":"yaoyun",
    "DM_USER_NAME":"yaoyun",
    "port":45316,
    "offset":1
}
`

func TestReceiveLog(t *testing.T) {
	s := Search{
		KeywordFilter: make(map[string]*list.List),
	}

	router := gin.New()
	router.POST("/rules", s.ReceiveLog)
	testServer := httptest.NewServer(router)
	assert.NotNil(t, testServer)
	defer testServer.Close()

	resp, err := http.Post(testServer.URL+"/rules", "application/json", bytes.NewReader([]byte(eventExample)))
	assert.Nil(t, err)
	assert.Equal(t, resp.StatusCode, http.StatusOK)

	resp, err = http.Post(testServer.URL+"/rules", "application/json", bytes.NewReader([]byte("xxxxxxx")))
	assert.Nil(t, err)
	assert.Equal(t, resp.StatusCode, http.StatusBadRequest)

}
