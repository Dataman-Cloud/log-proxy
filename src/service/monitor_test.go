package service

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

var (
	mux        *http.ServeMux
	testServer *httptest.Server
)

func setup() string {
	mux = http.NewServeMux()
	testServer = httptest.NewServer(mux)

	serverURL, _ := url.Parse(testServer.URL)
	return serverURL.String()
}

func teardown() {
	testServer.Close()
}

func testMethod(t *testing.T, r *http.Request, want string) {
	if got := r.Method; got != want {
		t.Errorf("Request method: %v, want %v", got, want)
	}
}

func getContentOfFile(fileName string) []byte {
	content, err := ioutil.ReadFile(fileName)
	if err != nil {
		return []byte{}
	}

	return content
}

func TestIsInArray(t *testing.T) {
	value := "test"
	values := make([]string, 0)
	values = append(values, "test")
	if !isInArray(values, value) {
		t.Error("Expect get true, but get false")
	}
	value = "test1"
	if isInArray(values, value) {
		t.Error("Expect get false, but get true")
	}
}
