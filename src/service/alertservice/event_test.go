package alertservice

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/Dataman-Cloud/log-proxy/src/models"
	mock_store "github.com/Dataman-Cloud/log-proxy/src/store/mock_datastore"

	"github.com/golang/mock/gomock"
)

func TestReceiveAlertEvent(t *testing.T) {
	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()
	mockStore := mock_store.NewMockStore(mockCtl)
	alert := initAlert()
	alert.Store = mockStore

	event := []byte(`{"receiver":"mola_webhook","status":"firing","alerts":[{"status":"firing","labels":{"alertname":"web_zdou_datamanmesos_cpu_usage_warning","app":"web-zdou-datamanmesos","container_label_DM_APP_ID":"web-zdou-datamanmesos","container_label_DM_APP_NAME":"web","container_label_DM_CLUSTER":"datamanmesos","container_label_DM_GROUP_NAME":"dev","container_label_DM_SLOT_ID":"0-web-zdou-datamanmesos","container_label_DM_SLOT_INDEX":"0","container_label_DM_TASK_ID":"0-web-zdou-datamanmesos-c19f31f345cd4110924974df54dac683","container_label_DM_USER":"zdou","container_label_DM_USER_NAME":"zdou","container_label_DM_VCLUSTER":"mola","duration":"5m","group":"dev","id":"/docker/4828b086a265a03ce4b34ac3662c6d67d5cb1a5230e1eb77a253b3fd382c1e19","image":"nginx:1.10","indicator":"cpu_usage","instance":"192.168.56.103:5014","job":"cadvisor","judgement":"max \u003e 60%","name":"mesos-3fca175f-7593-46ad-bcc6-2875e089f8b5-S0.c5a378af-7ea7-4cc7-9f80-de748cf0bbe9","severity":"warning","value":"0.01727926502595149"},"annotations":{"description":"","summary":""},"startsAt":"2017-05-06T20:42:27.804+08:00","endsAt":"0001-01-01T00:00:00Z","generatorURL":"http://srymaster1:9090/graph#%5B%7B%22expr%22%3A%22max%28irate%28container_cpu_usage_seconds_total%7Bcontainer_label_DM_APP_ID%3D%5C%22web-zdou-datamanmesos%5C%22%2Ccontainer_label_DM_LOG_TAG%21%3D%5C%22ignore%5C%22%2Cid%3D~%5C%22%2Fdocker%2F.%2A%5C%22%2Cname%3D~%5C%22mesos.%2A%5C%22%7D%5B5m%5D%29%29%20BY%20%28container_label_DM_APP_ID%29%20KEEP_COMMON%20%2A%20100%20%3C%2060%22%2C%22tab%22%3A0%7D%5D"}],"groupLabels":{"alertname":"web_zdou_datamanmesos_cpu_usage_warning"},"commonLabels":{"alertname":"web_zdou_datamanmesos_cpu_usage_warning","app":"web-zdou-datamanmesos","container_label_DM_APP_ID":"web-zdou-datamanmesos","container_label_DM_APP_NAME":"web","container_label_DM_CLUSTER":"datamanmesos","container_label_DM_GROUP_NAME":"dev","container_label_DM_SLOT_ID":"0-web-zdou-datamanmesos","container_label_DM_SLOT_INDEX":"0","container_label_DM_TASK_ID":"0-web-zdou-datamanmesos-c19f31f345cd4110924974df54dac683","container_label_DM_USER":"zdou","container_label_DM_USER_NAME":"zdou","container_label_DM_VCLUSTER":"mola","duration":"5m","group":"dev","id":"/docker/4828b086a265a03ce4b34ac3662c6d67d5cb1a5230e1eb77a253b3fd382c1e19","image":"nginx:1.10","indicator":"cpu_usage","instance":"192.168.56.103:5014","job":"cadvisor","judgement":"max \u003e 60%","name":"mesos-3fca175f-7593-46ad-bcc6-2875e089f8b5-S0.c5a378af-7ea7-4cc7-9f80-de748cf0bbe9","severity":"warning","value":"0.01727926502595149"},"commonAnnotations":{"description":"","summary":""},"externalURL":"http://srymaster1:9093","version":"3","groupKey":9849835813238314749}`)

	var message map[string]interface{}
	json.Unmarshal(event, &message)

	mockStore.EXPECT().CreateOrIncreaseEvent(gomock.Any()).Return(nil).Times(1)

	err := alert.ReceiveAlertEvent(message)
	if err != nil {
		t.Errorf("Expect error is nil, but got %v", err)
	}

	mockStore.EXPECT().CreateOrIncreaseEvent(gomock.Any()).Return(errors.New("err")).Times(1)
	err = alert.ReceiveAlertEvent(message)
	if err == nil {
		t.Errorf("Expect error is not nil, but got %v", err)
	}
}

func TestGetAlertEvents(t *testing.T) {
	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()
	mockStore := mock_store.NewMockStore(mockCtl)
	alert := initAlert()
	alert.Store = mockStore

	page := models.Page{}
	options := make(map[string]interface{})
	data := make(map[string]interface{})
	data["app"] = "app"

	mockStore.EXPECT().ListEvents(gomock.Any(), gomock.Any()).Return(data, nil).Times(1)
	result, err := alert.GetAlertEvents(page, options)
	if result == nil {
		t.Errorf("Expect result is not nil, but got %v", err)
	}
	if err != nil {
		t.Errorf("Expect error is nil, but got %v", err)
	}

}

func TestGetAlertEventsError(t *testing.T) {
	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()
	mockStore := mock_store.NewMockStore(mockCtl)
	alert := initAlert()
	alert.Store = mockStore

	page := models.Page{}
	options := make(map[string]interface{})

	mockStore.EXPECT().ListEvents(gomock.Any(), gomock.Any()).Return(nil, errors.New("err")).Times(1)
	result, err := alert.GetAlertEvents(page, options)
	if result != nil {
		t.Errorf("Expect result is  nil, but got %v", err)
	}
	if err == nil {
		t.Errorf("Expect error is not nil, but got %v", err)
	}
}

func TestAckAlertEvent(t *testing.T) {
	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()
	mockStore := mock_store.NewMockStore(mockCtl)
	alert := initAlert()
	alert.Store = mockStore

	options := make(map[string]interface{})
	options["group"] = "dev"
	options["app"] = "web-zdou-datamanmesos"
	options["action"] = "ack"

	mockStore.EXPECT().AckEvent(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).Times(1)
	err := alert.AckAlertEvent(1, options)
	if err != nil {
		t.Errorf("Expect error is nil, but got %v", err)
	}
}

func TestAckAlertEventError(t *testing.T) {
	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()
	mockStore := mock_store.NewMockStore(mockCtl)
	alert := initAlert()
	alert.Store = mockStore

	options := make(map[string]interface{})
	options["group"] = "dev"
	options["app"] = "web-zdou-datamanmesos"
	options["action"] = "ack"

	mockStore.EXPECT().AckEvent(gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("error")).Times(1)
	err := alert.AckAlertEvent(1, options)
	if err == nil {
		t.Errorf("Expect error is not nil, but got %v", err)
	}

	options["action"] = "err"
	err = alert.AckAlertEvent(1, options)
	if err == nil {
		t.Errorf("Expect error is not nil, but got %v", err)
	}
}
