package backends

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

var (
	mux    *http.ServeMux
	server *httptest.Server
)

func setup() string {
	mux = http.NewServeMux()
	server = httptest.NewServer(mux)

	serverURL, _ := url.Parse(server.URL)
	return serverURL.String()
}

func teardown() {
	server.Close()
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

func TestQueryExpr(t *testing.T) {
	server := setup()
	defer teardown()
	query := initQueryMetric()
	query.PromServer = server
	query.Expr = "container_tasks_state"
	mux.HandleFunc(query.Path, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		c := getContentOfFile("tests/query.node.response.json")
		fmt.Fprint(w, string(c))
	})

	data, err := query.QueryExpr()
	if data == nil {
		t.Errorf("Expect dat is not nil, got %v", data)
	}
	if err != nil {
		t.Errorf("Expect err is not nil, got %v", err)
	}

	query.PromServer = ""
	data, err = query.QueryExpr()
	if data != nil {
		t.Errorf("Expect data is nil, got %v", data)
	}
	if err == nil {
		t.Errorf("Expect err is nil, got %v", err)
	}
}

func TestQueryExprWrongJson(t *testing.T) {
	server := setup()
	defer teardown()
	query := initQueryMetric()
	query.PromServer = server
	query.Expr = "container_tasks_state"
	mux.HandleFunc(query.Path, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		c := getContentOfFile("tests/query.wrong.response.json")
		fmt.Fprint(w, string(c))
	})

	expectError := fmt.Sprintf("Failed to parse the response from %s%s?end=%s&query=%s&start=%s&step=%s", query.PromServer, query.Path, query.End, query.Expr, query.Start, query.Step)
	data, err := query.QueryExpr()
	if data != nil {
		t.Errorf("Expect dat is not nil, got %v", data)
	}
	if err.Error() != expectError {
		t.Errorf("Expect err is %s, got %v", expectError, err)
	}
}

func TestGetExprResponse(t *testing.T) {
	server := setup()
	defer teardown()

	query := initQueryMetric()
	query.PromServer = server
	path := query.Path
	mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		c := getContentOfFile("tests/query.default.response.json")
		fmt.Fprint(w, string(c))
	})
	expectRequest := fmt.Sprintf("%s%s?end=%s&query=%s&start=%s&step=%s", query.PromServer, query.Path, query.End, query.Expr, query.Start, query.Step)
	response, request, err := query.getExprResponse()
	if expectRequest != request {
		t.Errorf("Expect %v, got wrong request %v", expectRequest, request)
	}
	if response == nil {
		t.Errorf("Expect response is not nil, got wrong response %v", response)
	}
	if err != nil {
		t.Errorf("Expect err is not nil, got %v", err)
	}

	// step = 30s
	query.Step = ""
	expectRequest = fmt.Sprintf("%s%s?end=%s&query=%s&start=%s&step=30s", query.PromServer, query.Path, query.End, query.Expr, query.Start)
	response, request, err = query.getExprResponse()
	if expectRequest != request {
		t.Errorf("Expect %v, got wrong request %v", expectRequest, request)
	}
	if response == nil {
		t.Errorf("Expect response is not nil, got wrong response %v", response)
	}
	if err != nil {
		t.Errorf("Expect err is not nil, got %v", err)
	}
	// Empty URL
	query.PromServer = ""
	response, request, err = query.getExprResponse()
	expectRequest = fmt.Sprintf("%s%s?end=%s&query=%s&start=%s&step=30s", query.PromServer, query.Path, query.End, query.Expr, query.Start)
	if expectRequest != request {
		t.Errorf("Expect %v, got wrong request %v", expectRequest, request)
	}
	if response != nil {
		t.Errorf("Expect response is nil, got wrong response %v", response)
	}
	if err == nil {
		t.Errorf("Expect err is nil, got %v", err)
	}

	// QUERYPATH = query
	query.PromServer = server
	query.Path = QUERYPATH
	mux.HandleFunc(query.Path, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		c := getContentOfFile("tests/query.default.response.json")
		fmt.Fprint(w, string(c))
	})
	expectRequest = fmt.Sprintf("%s%s?query=%s", query.PromServer, query.Path, query.Expr)
	response, request, err = query.getExprResponse()
	if expectRequest != request {
		t.Errorf("Expect %v, got wrong request %v", expectRequest, request)
	}
	if response == nil {
		t.Errorf("Expect response is not nil, got wrong response %v", response)
	}
	if err != nil {
		t.Errorf("Expect err is not nil, got %v", err)
	}
}

func TestQueryMetric(t *testing.T) {
	server := setup()
	defer teardown()
	query := initQueryMetric()
	query.PromServer = server
	mux.HandleFunc(query.Path, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		c := getContentOfFile("tests/query.node.response.json")
		fmt.Fprint(w, string(c))
	})

	data, err := query.QueryMetric()
	if data.Status != "success" {
		t.Errorf("Expect status is success, got %v", data.Status)
	}
	if err != nil {
		t.Errorf("Expect err is nil, got %v", err)
	}
	query.PromServer = ""
	data, err = query.QueryMetric()
	if data != nil {
		t.Errorf("Expect data is nil, got %v", data)
	}
	if err == nil {
		t.Errorf("Expect err is not nil, got %v", err)
	}
}

func TestQueryMetricError(t *testing.T) {
	server := setup()
	defer teardown()
	query := initQueryMetric()
	query.PromServer = server
	mux.HandleFunc(query.Path, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		c := getContentOfFile("tests/query.error.response.json")
		fmt.Fprint(w, string(c))
	})

	data, err := query.QueryMetric()
	if data != nil {
		t.Errorf("Expect data is nil, got %v", data)
	}
	if err == nil {
		t.Errorf("Expect err is not nil, got %v", err)
	}
}

func TestTimeRange(t *testing.T) {
	query := initQueryMetric()
	start, end := timeRange(query.Start, query.End)
	if start != "1481612202" || end != "1481612212" {
		t.Errorf("Expect start 1481612202 and end 1481612212, got start %v, end %v", start, end)
	}

	query.Start = ""
	query.End = ""
	start, end = timeRange(query.Start, query.End)
	if start != end {
		t.Errorf("Expect start and end is equal, got start %v, end %v", start, end)
	}
}

func TestGetQueryMetricExpr(t *testing.T) {
	level := "id"
	expr := testGetTaskExpr()
	query := initQueryMetric()
	metrics := []string{"cpu", "memory", "memory_usage", "memory_total", "network_rx", "network_tx", "fs_read", "fs_write"}
	for _, metric := range metrics {
		query.Metric = metric
		newExpr := query.getQueryMetricExpr(level)

		switch query.Metric {
		case "cpu":
			if expr.CPU.Usage != newExpr {
				t.Errorf("Expect %v, got wrong expr %v", expr.CPU.Usage, newExpr)
			}
		case "memory":
			if expr.Memory.Percentage != newExpr {
				t.Errorf("Expect %v, got wrong expr %v", expr.Memory.Percentage, newExpr)
			}
		case "memory_usage":
			if expr.Memory.Usage != newExpr {
				t.Errorf("Expect %v, got wrong expr %v", expr.Memory.Usage, newExpr)
			}
		case "memory_total":
			if expr.Memory.Total != newExpr {
				t.Errorf("Expect %v, got wrong expr %v", expr.Memory.Total, newExpr)
			}
		case "network_rx":
			if expr.Network.Receive != newExpr {
				t.Errorf("Expect %v, got wrong expr %v", expr.Network.Receive, newExpr)
			}
		case "network_tx":
			if expr.Network.Transmit != newExpr {
				t.Errorf("Expect %v, got wrong expr %v", expr.Network.Transmit, newExpr)
			}
		case "fs_read":
			if expr.Filesystem.Read != newExpr {
				t.Errorf("Expect %v, got wrong expr %v", expr.Filesystem.Read, newExpr)
			}
		case "fs_write":
			if expr.Filesystem.Write != newExpr {
				t.Errorf("Expect %v, got wrong expr %v", expr.Filesystem.Write, newExpr)
			}
		}
	}
	level = "cluster"
	expr = testGetClusterExpr(query.NodeID)
	query = initQueryMetric()
	metrics = []string{"fs_usage", "fs_limit"}
	for _, metric := range metrics {
		query.Metric = metric
		newExpr := query.getQueryMetricExpr(level)

		switch query.Metric {

		case "fs_usage":
			if expr.Filesystem.Usage != newExpr {
				t.Errorf("Expect %v, got wrong expr %v", expr.Filesystem.Usage, newExpr)
			}
		case "fs_limit":
			if expr.Filesystem.Limit != newExpr {
				t.Errorf("Expect %v, got wrong expr %v", expr.Filesystem.Limit, newExpr)
			}
		}
	}
}

func TestQueryInfo(t *testing.T) {
	server := setup()
	defer teardown()
	query := initQueryMetric()
	query.PromServer = server
	mux.HandleFunc(query.Path, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		c := getContentOfFile("tests/query.node.response.json")
		fmt.Fprint(w, string(c))
	})

	data, err := query.QueryInfo()
	if data.Status != "success" {
		t.Errorf("Expect status is success, got %v", data.Status)
	}
	if err != nil {
		t.Errorf("Expect err is nil, got %v", err)
	}

	query.ClusterID = ""
	data, err = query.QueryInfo()
	if data.Status != "success" {
		t.Errorf("Expect status is success, got %v", data.Status)
	}
	if err != nil {
		t.Errorf("Expect err is nil, got %v", err)
	}

	query.ClusterID = ""
	query.AppID = ""
	data, err = query.QueryInfo()
	if data.Status != "success" {
		t.Errorf("Expect status is success, got %v", data.Status)
	}
	if err != nil {
		t.Errorf("Expect err is nil, got %v", err)
	}

}

func TestQueryInfoError(t *testing.T) {
	server := setup()
	defer teardown()
	query := initQueryMetric()
	query.PromServer = server
	mux.HandleFunc(query.Path, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		c := getContentOfFile("tests/query.error.response.json")
		fmt.Fprint(w, string(c))
	})

	data, err := query.QueryInfo()
	if data != nil {
		t.Errorf("Expect data is nil, got %v", data)
	}
	if err == nil {
		t.Errorf("Expect err is nil, got %v", err)
	}
}

func TestQueryInfoFailed(t *testing.T) {
	server := setup()
	defer teardown()
	query := initQueryMetric()
	query.PromServer = server
	mux.HandleFunc(query.Path, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		c := getContentOfFile("tests/query.error.response.json")
		fmt.Fprint(w, string(c))
	})

	data, err := query.QueryInfo()
	if data != nil {
		t.Errorf("Expect data is nil, got %v", data)
	}
	if err == nil {
		t.Errorf("Expect err is nil, got %v", err)
	}
}

func TestQueryAppMetric(t *testing.T) {
	server := setup()
	defer teardown()
	query := initQueryMetric()
	query.PromServer = server
	mux.HandleFunc(query.Path, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		c := getContentOfFile("tests/query.node.response.json")
		fmt.Fprint(w, string(c))
	})

	data, err := query.QueryAppMetric()
	if data.Status != "success" {
		t.Errorf("Expect status is success, got %v", data.Status)
	}
	if err != nil {
		t.Errorf("Expect err is nil, got %v", err)
	}
}

func TestQueryAppMetricError(t *testing.T) {
	server := setup()
	defer teardown()
	query := initQueryMetric()
	query.PromServer = server
	mux.HandleFunc(query.Path, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		c := getContentOfFile("tests/query.error.response.json")
		fmt.Fprint(w, string(c))
	})

	data, err := query.QueryAppMetric()
	if data != nil {
		t.Errorf("Expect data is nil, got %v", data)
	}
	if err == nil {
		t.Errorf("Expect err is nil, got %v", err)
	}
}

func TestQueryAppMetricFailed(t *testing.T) {
	server := setup()
	defer teardown()
	query := initQueryMetric()
	query.PromServer = server
	mux.HandleFunc(query.Path, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		c := getContentOfFile("tests/query.error.response.json")
		fmt.Fprint(w, string(c))
	})

	data, err := query.QueryAppMetric()
	if data != nil {
		t.Errorf("Expect data is nil, got %v", data)
	}
	if err == nil {
		t.Errorf("Expect err is nil, got %v", err)
	}
}

func TestQueryNodeMetric(t *testing.T) {
	server := setup()
	defer teardown()
	query := initQueryMetric()
	query.PromServer = server
	mux.HandleFunc(query.Path, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		c := getContentOfFile("tests/query.node.response.json")
		fmt.Fprint(w, string(c))
	})

	data, err := query.QueryNodeMetric()
	if data.Status != "success" {
		t.Errorf("Expect status is success, got %v", data.Status)
	}
	if err != nil {
		t.Errorf("Expect err is nil, got %v", err)
	}
}

func TestQueryNodeMetricError(t *testing.T) {
	server := setup()
	defer teardown()
	query := initQueryMetric()
	query.PromServer = server
	mux.HandleFunc(query.Path, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		c := getContentOfFile("tests/query.error.response.json")
		fmt.Fprint(w, string(c))
	})

	data, err := query.QueryNodeMetric()
	if data != nil {
		t.Errorf("Expect data is nil, got %v", data)
	}
	if err == nil {
		t.Errorf("Expect err is nil, got %v", err)
	}
}

func TestQueryNodeMetricFailed(t *testing.T) {
	server := setup()
	defer teardown()
	query := initQueryMetric()
	query.PromServer = server
	mux.HandleFunc(query.Path, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		c := getContentOfFile("tests/query.error.response.json")
		fmt.Fprint(w, string(c))
	})

	data, err := query.QueryNodeMetric()
	if data != nil {
		t.Errorf("Expect data is nil, got %v", data)
	}
	if err == nil {
		t.Errorf("Expect err is nil, got %v", err)
	}
}
