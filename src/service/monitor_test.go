package service

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/Dataman-Cloud/log-proxy/src/models"
	"github.com/Dataman-Cloud/log-proxy/src/utils/prometheusexpr"
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

func TestSetQueryExprsList(t *testing.T) {
	path := "../../config/exprs/"
	prometheusexpr.Exprs(path)
	exprTempl := SetQueryExprsList()
	if exprTempl == nil {
		t.Error("Expect exprTempl is not nil, but got nil")
	}
}

func TestGetQueryItemList(t *testing.T) {
	path := "../../config/exprs/"
	prometheusexpr.Exprs(path)
	query := &Query{
		ExprTmpl: SetQueryExprsList(),
	}
	items := query.GetQueryItemList()
	if items[0] != "Cpu_Usage_Percent" {
		t.Errorf("Expect the first item is Cpu_Usage_Percent, but got %s", items[0])
	}
}

func initQuery() *Query {
	param := &models.QueryParameter{}

	query := &Query{
		ExprTmpl:       SetQueryExprsList(),
		HTTPClient:     http.DefaultClient,
		PromServer:     "http://127.0.0.1:9090",
		Path:           QUERYRANGEPATH,
		QueryParameter: param,
	}
	return query
}

func TestGetQueryMetricExpr(t *testing.T) {
	expect := "avg(irate(container_cpu_usage_seconds_total{container_label_DM_APP_ID='dataman-app', container_label_DM_SLOT_INDEX='0', id=~'/docker/.*', name=~'mesos.*'}[5m])) by (container_label_DM_APP_ID, container_label_DM_SLOT_INDEX) keep_common"
	path := "../../config/exprs/"
	prometheusexpr.Exprs(path)

	query := initQuery()
	query.QueryParameter.Metric = "Cpu_Usage_Percent"
	query.QueryParameter.App = "dataman-app"
	query.QueryParameter.Task = "0"

	exprTempl := query.getQueryMetricExpr()
	if expect != exprTempl {
		t.Errorf("Expect the string is %s, but got %s", expect, exprTempl)
	}
}

func TestQueryExpr(t *testing.T) {
	server := setup()
	defer teardown()
	query := initQuery()
	query.PromServer = server
	query.Expr = "up"
	mux.HandleFunc(query.Path, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		c := getContentOfFile("tests/query.default.response.json")
		fmt.Fprint(w, string(c))
	})
	data, err := query.QueryExpr()
	if data == nil {
		t.Errorf("Expect not nil, but got %s", data)
	}
	if err != nil {
		t.Errorf("Expect err is nil, got %v", err)
	}
}

func TestQueryExprError(t *testing.T) {
	server := setup()
	defer teardown()
	query := initQuery()
	query.PromServer = server
	query.Expr = "up"
	mux.HandleFunc(query.Path, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		c := getContentOfFile("")
		fmt.Fprint(w, string(c))
	})
	data, err := query.QueryExpr()
	if data != nil {
		t.Errorf("Expect data is nil, but got %v", data)
	}
	if err == nil {
		t.Errorf("Expect err is not nil, got %v", err)
	}
}

func TestQueryMetric(t *testing.T) {
	path := "../../config/exprs/"
	prometheusexpr.Exprs(path)

	server := setup()
	defer teardown()
	query := initQuery()
	query.PromServer = server
	query.Metric = "Cpu_Usage_Percent"
	query.QueryParameter.App = "web-zdou-datamanmesos"
	query.QueryParameter.Task = "0"
	query.QueryParameter.Start = "1493697155"
	query.QueryParameter.End = "1493697155"
	query.QueryParameter.Step = "30"
	mux.HandleFunc(query.Path, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		c := getContentOfFile("tests/query.raw.response.json")
		fmt.Fprint(w, string(c))
	})
	data, err := query.QueryMetric()
	if data == nil {
		t.Errorf("Expect not nil, but got %s", data)
	}
	if err != nil {
		t.Errorf("Expect err is nil, got %v", err)
	}
}

func TestQueryMetricError(t *testing.T) {
	path := "../../config/exprs/"
	prometheusexpr.Exprs(path)

	server := setup()
	defer teardown()
	query := initQuery()
	query.PromServer = server
	query.Metric = "Cpu_Usage_Percent"
	query.QueryParameter.App = "web-zdou-datamanmesos"
	query.QueryParameter.Task = "0"
	query.QueryParameter.Start = "1493697155"
	query.QueryParameter.End = "1493697155"
	query.QueryParameter.Step = "30"
	mux.HandleFunc(query.Path, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		c := getContentOfFile("")
		fmt.Fprint(w, string(c))
	})
	data, err := query.QueryMetric()
	if data != nil {
		t.Errorf("Expect nil, but got %s", data)
	}
	if err == nil {
		t.Errorf("Expect err is not nil, got %v", err)
	}
}

func TestQueryMetricErrorResult(t *testing.T) {
	path := "../../config/exprs/"
	prometheusexpr.Exprs(path)

	server := setup()
	defer teardown()
	query := initQuery()
	query.PromServer = server
	query.Metric = "Cpu_Usage_Percent"
	query.QueryParameter.App = "web-zdou-datamanmesos"
	query.QueryParameter.Task = "0"
	query.QueryParameter.Start = "1493697155"
	query.QueryParameter.End = "1493697155"
	query.QueryParameter.Step = "30"
	mux.HandleFunc(query.Path, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		c := getContentOfFile("tests/query.default.response.json")
		fmt.Fprint(w, string(c))
	})
	data, err := query.QueryMetric()
	if data != nil {
		t.Errorf("Expect nil, but got %s", data)
	}
	if err == nil {
		t.Errorf("Expect err is not nil, got %v", err)
	}
}

func TestQueryMetricErrorServer(t *testing.T) {
	path := "../../config/exprs/"
	prometheusexpr.Exprs(path)

	server := setup()
	defer teardown()
	query := initQuery()
	query.PromServer = server + "err"
	query.Metric = "Cpu_Usage_Percent"
	query.QueryParameter.App = "web-zdou-datamanmesos"
	query.QueryParameter.Task = "0"
	query.QueryParameter.Start = "1493697155"
	query.QueryParameter.End = "1493697155"
	query.QueryParameter.Step = "30"
	mux.HandleFunc(query.Path, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		c := getContentOfFile("tests/query.default.response.json")
		fmt.Fprint(w, string(c))
	})
	data, err := query.QueryMetric()
	if data != nil {
		t.Errorf("Expect nil, but got %s", data)
	}
	if err == nil {
		t.Errorf("Expect err is not nil, got %v", err)
	}
}

func TestQueryApps(t *testing.T) {
	path := "../../config/exprs/"
	prometheusexpr.Exprs(path)

	server := setup()
	defer teardown()
	query := initQuery()
	query.PromServer = server
	mux.HandleFunc(query.Path, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		c := getContentOfFile("tests/query.apps.response.json")
		fmt.Fprint(w, string(c))
	})
	expect := "web-zdou-datamanmesos"
	data, err := query.GetQueryApps()
	if data[0] != expect {
		t.Errorf("Expect got %s, but got %s", expect, data[0])
	}
	if err != nil {
		t.Errorf("Expect err is nil, got %v", err)
	}
}

func TestQueryAppsError(t *testing.T) {
	path := "../../config/exprs/"
	prometheusexpr.Exprs(path)

	server := setup()
	defer teardown()
	query := initQuery()
	query.PromServer = server + "err"
	mux.HandleFunc(query.Path, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		c := getContentOfFile("tests/query.apps.response.json")
		fmt.Fprint(w, string(c))
	})
	data, err := query.GetQueryApps()
	if data != nil {
		t.Errorf("Expect nil, but got %s", data)
	}
	if err == nil {
		t.Errorf("Expect err is not nil, got %v", err)
	}
}

func TestQueryAppsErrorResponse(t *testing.T) {
	path := "../../config/exprs/"
	prometheusexpr.Exprs(path)

	server := setup()
	defer teardown()
	query := initQuery()
	query.PromServer = server
	mux.HandleFunc(query.Path, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		c := getContentOfFile("")
		fmt.Fprint(w, string(c))
	})
	data, err := query.GetQueryApps()
	if data != nil {
		t.Errorf("Expect nil, but got %s", data)
	}
	if err == nil {
		t.Errorf("Expect err is not nil, got %v", err)
	}
}

func TestQueryAppsErrorResult(t *testing.T) {
	path := "../../config/exprs/"
	prometheusexpr.Exprs(path)

	server := setup()
	defer teardown()
	query := initQuery()
	query.PromServer = server
	mux.HandleFunc(query.Path, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		c := getContentOfFile("tests/query.default.response.json")
		fmt.Fprint(w, string(c))
	})
	data, err := query.GetQueryApps()
	if data != nil {
		t.Errorf("Expect nil, but got %s", data)
	}
	if err == nil {
		t.Errorf("Expect err is not nil, got %v", err)
	}
}

func TestQueryAppTasks(t *testing.T) {
	path := "../../config/exprs/"
	prometheusexpr.Exprs(path)

	server := setup()
	defer teardown()
	query := initQuery()
	query.PromServer = server
	mux.HandleFunc(query.Path, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		c := getContentOfFile("tests/query.tasks.response.json")
		fmt.Fprint(w, string(c))
	})
	expect := "0"
	data, err := query.GetQueryAppTasks()
	if data[0] != expect {
		t.Errorf("Expect got %s, but got %s", expect, data[0])
	}
	if err != nil {
		t.Errorf("Expect err is nil, got %v", err)
	}
}

func TestQueryAppTasksError(t *testing.T) {
	path := "../../config/exprs/"
	prometheusexpr.Exprs(path)

	server := setup()
	defer teardown()
	query := initQuery()
	query.PromServer = server + "err"
	mux.HandleFunc(query.Path, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		c := getContentOfFile("tests/query.apps.response.json")
		fmt.Fprint(w, string(c))
	})
	data, err := query.GetQueryAppTasks()
	if data != nil {
		t.Errorf("Expect nil, but got %s", data)
	}
	if err == nil {
		t.Errorf("Expect err is not nil, got %v", err)
	}
}

func TestQueryAppTasksErrorResponse(t *testing.T) {
	path := "../../config/exprs/"
	prometheusexpr.Exprs(path)

	server := setup()
	defer teardown()
	query := initQuery()
	query.PromServer = server
	mux.HandleFunc(query.Path, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		c := getContentOfFile("")
		fmt.Fprint(w, string(c))
	})
	data, err := query.GetQueryAppTasks()
	if data != nil {
		t.Errorf("Expect nil, but got %s", data)
	}
	if err == nil {
		t.Errorf("Expect err is not nil, got %v", err)
	}
}

func TestQueryAppTasksErrorResult(t *testing.T) {
	path := "../../config/exprs/"
	prometheusexpr.Exprs(path)

	server := setup()
	defer teardown()
	query := initQuery()
	query.PromServer = server
	mux.HandleFunc(query.Path, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		c := getContentOfFile("tests/query.default.response.json")
		fmt.Fprint(w, string(c))
	})
	data, err := query.GetQueryAppTasks()
	if data != nil {
		t.Errorf("Expect nil, but got %s", data)
	}
	if err == nil {
		t.Errorf("Expect err is not nil, got %v", err)
	}
}
