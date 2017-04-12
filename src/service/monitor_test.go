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
	mux           *http.ServeMux
	backendServer *httptest.Server
)

func setup() string {
	mux = http.NewServeMux()
	backendServer = httptest.NewServer(mux)
	serverURL, _ := url.Parse(backendServer.URL)
	return serverURL.String()
}

func teardown() {
	backendServer.Close()
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
	if items[0] != "CPU使用率" {
		t.Errorf("Expect the first item is CPU使用率, but got %s", items[0])
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

func TestGetQueryUsers(t *testing.T) {
	server := setup()
	defer teardown()
	query := initQuery()
	query.PromServer = server
	mux.HandleFunc(query.Path, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		c := getContentOfFile("tests/query.users.response.json")
		fmt.Fprint(w, string(c))
	})
	expectUser := "zdou"
	data, err := query.GetQueryUsers()
	if data[0] != expectUser {
		t.Errorf("Expect get %s, but got %s", expectUser, data[0])
	}
	if err != nil {
		t.Errorf("Expect err is nil, got %v", err)
	}
}

func TestGetQueryUsersError(t *testing.T) {
	server := setup()
	defer teardown()
	query := initQuery()
	query.PromServer = server
	mux.HandleFunc(query.Path, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		c := getContentOfFile("")
		fmt.Fprint(w, string(c))
	})
	data, err := query.GetQueryUsers()
	if data != nil {
		t.Errorf("Expect data is nil, but got %v", data)
	}
	if err == nil {
		t.Errorf("Expect err is not nil, got %v", err)
	}
}

func TestGetQueryUsersWrongURL(t *testing.T) {
	server := setup()
	defer teardown()
	query := initQuery()
	query.PromServer = server + "wrong"
	mux.HandleFunc(query.Path, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		c := getContentOfFile("")
		fmt.Fprint(w, string(c))
	})
	data, err := query.GetQueryUsers()
	if data != nil {
		t.Errorf("Expect data is nil, but got %v", data)
	}
	if err == nil {
		t.Errorf("Expect err is not nil, got %v", err)
	}
}
