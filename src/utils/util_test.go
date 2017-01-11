package utils

import (
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	ms "github.com/Dataman-Cloud/mock-server/server"
)

func TestByte2str(t *testing.T) {
	if Byte2str([]byte("test")) == "test" {
		t.Log("success")
	} else {
		t.Error("faild")
	}
}

func TestParseDate(t *testing.T) {
	if ParseDate(nil, nil) == "*" {
		t.Log("success")
	} else {
		t.Error("faild")
	}

	if ParseDate("test", "1481596014000") == time.Now().Format("2006-01-02") {
		t.Log("success")
	} else {
		t.Error("faild")
	}

	if ParseDate("1481596014000", "test") == time.Now().Format("2006-01-02") {
		t.Log("success")
	} else {
		t.Error("faild")
	}

	if s := ParseDate("1481596014000", "1481336814000"); s == "*" {
		t.Log("success")
	} else {
		t.Error("faild", s)
	}

	if ParseDate("1481596014000", "1481596014000") == time.Now().Format("2006-01-02") {
		t.Log("success")
	} else {
		t.Error("faild")
	}
}

func TestReadRequestBody(t *testing.T) {
	req, _ := http.NewRequest("GET", "http://test.com", strings.NewReader("test"))
	if _, err := ReadRequestBody(req); err == nil {
		t.Log("success")
	} else {
		t.Error("faild")
	}
}

func TestReadResponseBody(t *testing.T) {
	mockServer := ms.NewServer()
	defer mockServer.Close()

	mockServer.AddRouter("/_ping", "get").RGroup().Reply(200).WBodyString(`{"message": "fails to get tasks"}`)

	resp, _ := http.Get(fmt.Sprintf("http://%s:%s/_ping", mockServer.Addr, mockServer.Port))
	if _, err := ReadResponseBody(resp); err == nil {
		t.Log("success")
	} else {
		t.Error("faild")
	}

}

func TestAlertNotification(t *testing.T) {
	if AlertNotification(fmt.Sprintf("http://%s:%s/notification", "localhost", "30859"), "test") != nil {
		t.Log("success")
	} else {
		t.Error("faild")
	}

	if AlertNotification(fmt.Sprintf("http://%s:%s/notification", "localhost", "30859"), `{"test":"value"}`) != nil {
		t.Log("success")
	} else {
		t.Error("faild")
	}

	mockServer := ms.NewServer()
	defer mockServer.Close()

	mockServer.AddRouter("/notification", "post").RGroup().Reply(200).WBodyString(`{"status": "success"}`)
	if AlertNotification(fmt.Sprintf("http://%s:%s/notification", mockServer.Addr, mockServer.Port), `{"test":"value"}`) == nil {
		t.Log("success")
	} else {
		t.Error("faild")
	}
}

func TestParseTask(t *testing.T) {
	if len(ParseTask("1")) == 1 {
		t.Log("success")
	} else {
		t.Error("faild")
	}

	if len(ParseTask("s-1")) == 1 {
		t.Log("success")
	} else {
		t.Error("faild")
	}

	if len(ParseTask("1-s")) == 1 {
		t.Log("success")
	} else {
		t.Error("faild")
	}

	if len(ParseTask("3-5")) == 3 {
		t.Log("success")
	} else {
		t.Error("faild")
	}
}

func TestParseMonitorTaskID(t *testing.T) {
	taskID := "1-5"
	expectString := "1|2|3|4|5"
	taskString, _ := ParseMonitorTaskID(taskID)
	if taskString != expectString {
		t.Errorf("Expect task string is %s, got %v", expectString, taskString)
	}

	taskID = ""
	expectString = ".*"
	taskString, _ = ParseMonitorTaskID(taskID)
	if taskString != expectString {
		t.Errorf("Expect task string is %s, got %v", expectString, taskString)
	}

	taskID = "1,2,3,4,5"
	expectString = "1|2|3|4|5"
	taskString, _ = ParseMonitorTaskID(taskID)
	if taskString != expectString {
		t.Errorf("Expect task string is %s, got %v", expectString, taskString)
	}

	taskID = "1"
	expectString = "1"
	taskString, _ = ParseMonitorTaskID(taskID)
	if taskString != expectString {
		t.Errorf("Expect task string is %s, got %v", expectString, taskString)
	}

}
