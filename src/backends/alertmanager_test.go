package backends

import (
	"fmt"
	"net/http"
	"testing"
)

func initQueryAlertManager() *AlertManager {
	return &AlertManager{
		HttpClient: http.DefaultClient,
		Server:     "http://127.0.0.1:9093",
		Path:       ALERTSPATH,
	}
}

func TestGetResponse(t *testing.T) {
	server := setup()
	defer teardown()

	am := initQueryAlertManager()
	am.Server = server

	mux.HandleFunc(am.Path, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		c := getContentOfFile("tests/alertmanager.default.response.json")
		fmt.Fprint(w, string(c))
	})

	response, err := am.getResponse()
	if response == nil {
		t.Errorf("Expect response is not nil, got wrong response %v", response)
	}
	if err != nil {
		t.Errorf("Expect err is nil, got %v", err)
	}

	am.Server = ""
	response, err = am.getResponse()
	if response != nil {
		t.Errorf("Expect response is nil, got wrong response %v", response)
	}
	if err == nil {
		t.Errorf("Expect err is not nil, got %v", err)
	}
}

func TestGetAlertManagerResponse(t *testing.T) {
	server := setup()
	defer teardown()

	am := initQueryAlertManager()
	am.Server = server

	mux.HandleFunc(am.Path, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		c := getContentOfFile("tests/alertmanager.default.response.json")
		fmt.Fprint(w, string(c))
	})

	data, err := am.GetAlertManagerResponse()
	if data["code"] != float64(0) {
		t.Errorf("Expect code is 0, got %v", data["code"])
	}
	if err != nil {
		t.Errorf("Expect err is nil, got %v", err)
	}
}

func TestGetAlertManagerResponseFailed(t *testing.T) {
	server := setup()
	defer teardown()

	am := initQueryAlertManager()
	am.Server = server

	mux.HandleFunc(am.Path, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		c := getContentOfFile("tests/alertmanager.error.response.json")
		fmt.Fprint(w, string(c))
	})

	data, err := am.GetAlertManagerResponse()
	if data != nil {
		t.Errorf("Expect data is nil, got %v", data)
	}
	if err == nil {
		t.Errorf("Expect err is not nil, got %v", err)
	}
}
