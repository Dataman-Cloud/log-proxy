package utils

import (
	"fmt"
	"net/http"
	"strings"
	"testing"

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

	if ParseDate("test", "1481596014000") == "2016-12-13" {
		t.Log("success")
	} else {
		t.Error("faild")
	}

	if ParseDate("1481596014000", "test") == "2016-12-13" {
		t.Log("success")
	} else {
		t.Error("faild")
	}

	if s := ParseDate("1481596014000", "1481336814000"); s == "*" {
		t.Log("success")
	} else {
		t.Error("faild", s)
	}

	if ParseDate("1481596014000", "1481596014000") == "2016-12-13" {
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