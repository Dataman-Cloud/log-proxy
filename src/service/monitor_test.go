package service

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/Dataman-Cloud/log-proxy/src/backends"
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

func initQueryMetric() *backends.Query {
	param := &backends.QueryParameter{
		Metric:    "cpu",
		ClusterID: "",
		AppID:     "",
		SlotID:    "",
		UserID:    "user1",
		Start:     "1481767298",
		End:       "1481767298",
		Step:      "100s",
	}

	query := &backends.Query{
		HTTPClient:     http.DefaultClient,
		PromServer:     "http://127.0.0.1:9090",
		Path:           backends.QUERYRANGEPATH,
		QueryParameter: param,
	}
	return query
}

func TestGetQueryMetric(t *testing.T) {
	server := setup()
	defer teardown()
	query := initQueryMetric()
	query.PromServer = server
	mux.HandleFunc(query.Path, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		c := getContentOfFile("tests/query.default.response.json")
		fmt.Fprint(w, string(c))
	})
	metrics := []string{"cpu", "memory", "memory_usage", "memory_total", "network_rx", "network_tx", "fs_read", "fs_write"}
	for _, metric := range metrics {
		query.Metric = metric
		data := NewMetric()
		err := data.GetQueryMetric(query)
		switch query.Metric {
		case "cpu":
			if data.CPU.Usage == nil {
				t.Errorf("Expect data.CPU.Usage is not nil, got %v", data.CPU.Usage)
			}
		case "memory":
			if data.Memory.Percentage == nil {
				t.Errorf("Expect data.Memory.Percentage is not nil, got %v", data.Memory.Percentage)
			}
		case "memory_usage":
			if data.Memory.Usage == nil {
				t.Errorf("Expect data.Memory.Usage is not nil, got %v", data.Memory.Usage)
			}
		case "memory_total":
			if data.Memory.Total == nil {
				t.Errorf("Expect data.Memory.Total is not nil, got %v", data.Memory.Total)
			}
		case "network_rx":
			if data.Network.Receive == nil {
				t.Errorf("Expect data.Network.Receive is not nil, got %v", data.Network.Receive)
			}
		case "network_tx":
			if data.Network.Transmit == nil {
				t.Errorf("Expect data.Network.Transmit is not nil, got %v", data.Network.Transmit)
			}
		case "fs_read":
			if data.Filesystem.Read == nil {
				t.Errorf("Expect data.Filesystem.Read is not nil, got %v", data.Filesystem.Read)
			}
		case "fs_write":
			if data.Filesystem.Write == nil {
				t.Errorf("Expect data.Filesystem.Write is not nil, got %v", data.Filesystem.Write)
			}
		}
		if err != nil {
			t.Errorf("Expect err is not nil, got %v", err)
		}
	}
	query.Metric = "metric"
	data := NewMetric()
	err := data.GetQueryMetric(query)
	if data.CPU.Usage != nil {
		t.Errorf("Expect data.CPU.Usage is not nil, got %v", data.CPU.Usage)
	}
	if err == nil {
		t.Errorf("Expect err is not nil, got %v", err)
	}
}

func TestGetQueryMetricError(t *testing.T) {
	server := setup()
	defer teardown()
	query := initQueryMetric()
	query.PromServer = server + "wrongURL"
	mux.HandleFunc(query.Path, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		c := getContentOfFile("tests/query.default.response.json")
		fmt.Fprint(w, string(c))
	})
	data := NewMetric()
	err := data.GetQueryMetric(query)
	if data.CPU.Usage != nil {
		t.Errorf("Expect data.CPU.Usage is not nil, got %v", data.CPU.Usage)
	}
	if err == nil {
		t.Errorf("Expect err is not nil, got %v", err)
	}
}

func TestGetQueryInfo(t *testing.T) {
	server := setup()
	defer teardown()
	query := initQueryMetric()
	query.PromServer = server
	mux.HandleFunc(query.Path, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		c := getContentOfFile("tests/query.default.response.json")
		fmt.Fprint(w, string(c))
	})
	data := NewInfo()
	err := data.GetQueryInfo(query)
	for ClusterName, ClusterValue := range data.Clusters {
		if ClusterName == "work" {
			for name := range ClusterValue.Users {
				if name == "user1" {
					break
				}
				t.Errorf("Expect user user1 in list, but missing")
			}
			break
		}
		t.Errorf("Expect cluster work and user user1 in list, but missing")
	}
	if err != nil {
		t.Errorf("Expect err is not nil, got %v", err)
	}
}

func TestGetQueryAppInfo(t *testing.T) {
	server := setup()
	defer teardown()
	query := initQueryMetric()
	query.PromServer = server
	mux.HandleFunc(query.Path, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		c := getContentOfFile("tests/query.default.response.json")
		fmt.Fprint(w, string(c))
	})
	query.ClusterID = "work"
	query.AppID = "work-nginx"
	query.UserID = "user1"
	data := NewInfo()
	err := data.GetQueryInfo(query)
	for ClusterName, ClusterValue := range data.Clusters {
		if ClusterName == "work" {
			for userName, userValue := range ClusterValue.Users {
				if userName == "user1" {
					for _, value := range userValue.Applications {
						if value.CPU.Usage == nil {
							t.Errorf("Expect value.CPU.Usage is not nil, got %v", value.CPU.Usage)
						}
					}
				}
			}
		}
	}
	if err != nil {
		t.Errorf("Expect err is nil, got %v", err)
	}
}

func TestGetQueryInfoError(t *testing.T) {
	server := setup()
	defer teardown()
	query := initQueryMetric()
	query.PromServer = server + "wrongURL"
	mux.HandleFunc(query.Path, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		c := getContentOfFile("tests/query.default.response.json")
		fmt.Fprint(w, string(c))
	})
	data := NewInfo()
	err := data.GetQueryInfo(query)

	if len(data.Clusters) != 0 {
		t.Errorf("Expect len of data.Clusters is 0, got %v", len(data.Clusters))
	}
	if err == nil {
		t.Errorf("Expect err is not nil, got %v", err)
	}
}

func TestGetQueryNodesInfo(t *testing.T) {
	server := setup()
	defer teardown()
	query := initQueryMetric()
	query.PromServer = server
	mux.HandleFunc(query.Path, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		c := getContentOfFile("tests/query.node.response.json")
		fmt.Fprint(w, string(c))
	})
	data := NewNodesInfo()
	err := data.GetQueryNodesInfo(query)
	if len(data.Nodes) == 0 {
		t.Errorf("Expect len of data.Nodes is not 0, got %v", len(data.Nodes))
	}
	if err != nil {
		t.Errorf("Expect err is nil, got %v", err)
	}
}

func TestGetQueryNodesInfoError(t *testing.T) {
	server := setup()
	defer teardown()
	query := initQueryMetric()
	query.PromServer = server + "wrongURL"
	mux.HandleFunc(query.Path, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		c := getContentOfFile("tests/query.node.response.json")
		fmt.Fprint(w, string(c))
	})
	data := NewNodesInfo()
	err := data.GetQueryNodesInfo(query)
	if len(data.Nodes) != 0 {
		t.Errorf("Expect len of data.Nodes is 0, got %v", len(data.Nodes))
	}
	if err == nil {
		t.Errorf("Expect err is not nil, got %v", err)
	}
}
