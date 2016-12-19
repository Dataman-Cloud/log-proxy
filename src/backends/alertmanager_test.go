package backends

import (
	"encoding/json"
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

func TestGetSilences(t *testing.T) {
	server := setup()
	defer teardown()

	am := initQueryAlertManager()
	am.Server = server
	am.Path = SILENCES_API

	mux.HandleFunc(am.Path, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"status":"success","data":{"silences":[{"id":1,"matchers":[{"name":"alertname","value":"LogKeyword","isRegex":true},{"name":"appid","value":"work-nginx","isRegex":false},{"name":"clusterid","value":"work","isRegex":false},{"name":"instance","value":"192.168.1.75:5098","isRegex":false},{"name":"job","value":"log-proxy","isRegex":false},{"name":"keyword","value":"GET","isRegex":false},{"name":"offset","value":"1481781258186188032","isRegex":false},{"name":"path","value":"stdout","isRegex":false},{"name":"severity","value":"Warning","isRegex":false},{"name":"taskid","value":"work-nginx.1f17a9f0-c02b-11e6-9030-024245dc84c8","isRegex":false},{"name":"userid","value":"4","isRegex":false}],"startsAt":"2016-12-16T03:06:00Z","endsAt":"2016-12-16T07:06:00Z","createdAt":"2016-12-16T13:02:23.067679215+08:00","createdBy":"zqdou@dataman-inc.com","comment":"test"}],"totalSilences":1}}`)
	})
	_, err := am.GetSilences()
	if err == nil {
		t.Log("success")
	} else {
		t.Error("faild")
	}
}

func TestCreateSilence(t *testing.T) {
	server := setup()
	defer teardown()

	am := initQueryAlertManager()
	am.Server = server
	am.Path = SILENCES_API

	mux.HandleFunc(am.Path, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		fmt.Fprint(w, `{"status":"success","data":{"silenceId":2}}`)
	})
	body := `{"matchers":[{"name":"alertname","value":"LogKeyword","isRegex":false},{"name":"appid","value":"work-nginx","isRegex":false},{"name":"clusterid","value":"work","isRegex":false},{"name":"instance","value":"192.168.1.75:5098","isRegex":false},{"name":"job","value":"log-proxy","isRegex":false},{"name":"keyword","value":"GET","isRegex":false},{"name":"offset","value":"1481781258185649664","isRegex":false},{"name":"path","value":"stdout","isRegex":false},{"name":"severity","value":"Warning","isRegex":false},{"name":"taskid","value":"work-nginx.1f17a9f0-c02b-11e6-9030-024245dc84c8","isRegex":false},{"name":"userid","value":"4","isRegex":false}],"startsAt":"2016-12-19T03:12:00.000Z","endsAt":"2016-12-19T07:12:00.000Z","createdBy":"test@123.com","comment":"asdfasdf"}`
	var m map[string]interface{}
	json.Unmarshal([]byte(body), &m)
	err := am.CreateSilence(m)
	if err == nil {
		t.Log("success")
	} else {
		t.Error("faild")
	}
}

func TestGetSilence(t *testing.T) {
	server := setup()
	defer teardown()

	am := initQueryAlertManager()
	am.Server = server
	am.Path = GET_SILENCE + "2"

	mux.HandleFunc(am.Path, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"status":"success","data":{"id":2,"matchers":[{"name":"alertname","value":"LogKeyword","isRegex":false},{"name":"appid","value":"work-nginx","isRegex":false},{"name":"clusterid","value":"work","isRegex":false},{"name":"instance","value":"192.168.1.75:5098","isRegex":false},{"name":"job","value":"log-proxy","isRegex":false},{"name":"keyword","value":"GET","isRegex":false},{"name":"offset","value":"1481781258185649664","isRegex":false},{"name":"path","value":"stdout","isRegex":false},{"name":"severity","value":"Warning","isRegex":false},{"name":"taskid","value":"work-nginx.1f17a9f0-c02b-11e6-9030-024245dc84c8","isRegex":false},{"name":"userid","value":"4","isRegex":false}],"startsAt":"2016-12-19T03:12:00Z","endsAt":"2016-12-19T07:12:00Z","createdAt":"2016-12-19T11:08:14.816592362+08:00","createdBy":"test@123.com","comment":"asdfasdf"}}`)
	})

	_, err := am.GetSilence("2")
	if err == nil {
		t.Log("success")
	} else {
		t.Error("faild")
	}
}

func TestDeleteSilence(t *testing.T) {
	server := setup()
	defer teardown()

	am := initQueryAlertManager()
	am.Server = server
	am.Path = GET_SILENCE + "2"

	mux.HandleFunc(am.Path, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		fmt.Fprint(w, `{"status":"success"}`)
	})

	err := am.DeleteSilence("2")
	if err == nil {
		t.Log("success")
	} else {
		t.Error("faild")
	}
}
