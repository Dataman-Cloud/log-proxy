package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"

	"github.com/Dataman-Cloud/log-proxy/src/config"
	"github.com/Dataman-Cloud/log-proxy/src/models"
	mock_store "github.com/Dataman-Cloud/log-proxy/src/store/mock_datastore"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gopkg.in/olivere/elastic.v3"
)

var (
	apiURL  string
	baseURL string
	server  *httptest.Server
	s       *Search
	mo      *Monitor
)

func TestGetLogAlertRuleIndex(t *testing.T) {
	rule := models.LogAlertRule{
		Group:  "g",
		User:   "u",
		App:    "a",
		Source: "s",
	}

	ruleIndex := getLogAlertRuleIndex(rule)
	assert.Equal(t, ruleIndex, "g-u-a-s")
}

func TestInitLogAlertFilter(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockStore := mock_store.NewMockStore(mockCtrl)
	s := GetSearch()
	s.Store = mockStore

	rule := models.LogAlertRule{
		App:     "app",
		Cluster: "cluster",
		Keyword: "key",
		Source:  "stdout",
		User:    "user",
		Group:   "group",
	}

	rules := []*models.LogAlertRule{&rule}
	result := map[string]interface{}{"rules": rules}
	mockStore.EXPECT().GetLogAlertRules(gomock.Any(), gomock.Any()).Return(result, nil).Times(1)
	s.InitLogKeywordFilter()

	mockStore.EXPECT().GetLogAlertRules(gomock.Any(), gomock.Any()).Return(nil, nil).Times(1)
	s.InitLogKeywordFilter()

	mockStore.EXPECT().GetLogAlertRules(gomock.Any(), gomock.Any()).Return(nil, errors.New("test")).Times(1)
	s.InitLogKeywordFilter()
}

func startAPIServer(sv *Search) *httptest.Server {
	router := gin.New()
	v1 := router.Group("/api/v1", func(ctx *gin.Context) { ctx.Set("page", models.Page{}) })
	{
		v1.GET("/ping", func(ctx *gin.Context) { sv.Ping(ctx) })
		v1.GET("/clusters", func(ctx *gin.Context) { sv.Clusters(ctx) })
		v1.GET("/clusters/:cluster/apps", func(ctx *gin.Context) { sv.Applications(ctx) })
		v1.GET("/clusters/:cluster/apps/:app/slots", func(ctx *gin.Context) { sv.Slots(ctx) })
		v1.GET("/clusters/:cluster/apps/:app/slots/:slot/tasks", func(ctx *gin.Context) { sv.Tasks(ctx) })
		v1.GET("/clusters/:cluster/apps/:app/sources", func(ctx *gin.Context) { sv.Sources(ctx) })
		v1.GET("/clusters/:cluster/apps/:app/search", func(ctx *gin.Context) { sv.Search(ctx) })
		v1.GET("/clusters/:cluster/apps/:app/context", func(ctx *gin.Context) { sv.Context(ctx) })
	}

	vr := router.Group("/v1/receive")
	{
		vr.POST("/log", receiverlog)
	}

	return httptest.NewServer(router)
}

func startErrorClient() *httptest.Server {
	router := gin.New()
	router.HEAD("/", func(ctx *gin.Context) { ctx.String(200, "") })
	router.GET("/_nodes/http", nodes)
	router.POST("/:index/_search", func(ctx *gin.Context) { ctx.String(503, "error") })
	return httptest.NewServer(router)
}

func startHTTPServer() *httptest.Server {
	router := gin.New()
	router.HEAD("/", func(ctx *gin.Context) { ctx.String(200, "") })
	router.GET("/_nodes/http", nodes)
	router.GET("/api/v1/query_rang", queryResult)
	router.GET("/api/v1/alerts", queryAlerts)
	router.GET("/api/v1/alerts/groups", queryAlertsGroups)
	router.GET("/api/v1/alerts/status", queryAlertsStatus)
	router.GET("/api/v1/silences", querySliences)

	return httptest.NewServer(router)
}

func getp(ctx *gin.Context) {
	data := `{"_index":"1","_type":"1","_id":"1","_source":{"test":"value"}}`
	var info elastic.GetResult
	json.Unmarshal([]byte(data), &info)

	ctx.JSON(http.StatusOK, info)
}

func app(ctx *gin.Context) {
	data := `
	{
    "took":137,
    "_scroll_id":"",
    "hits":{
        "total":6,
        "max_score":0,
        "hits":[

        ]
    },
    "suggest":null,
    "aggregations":{
        "apps":{
            "doc_count_error_upper_bound":0,
            "sum_other_doc_count":0,
            "buckets":[
                {
                    "key":"test-web",
                    "doc_count":6
                }
            ]
        }
    },
    "timed_out":false,
    "terminated_early":false,
    "_shards":{
        "total":5,
        "successful":5,
        "failed":0
    }
}
`
	var info elastic.SearchResult
	json.Unmarshal([]byte(data), &info)

	ctx.JSON(http.StatusOK, info)
}

func nodes(ctx *gin.Context) {
	u, _ := url.Parse(baseURL)
	data := `{"cluster_name":"elasticsearch","nodes":{"Ijb_-48ZQYmEnQ0a5BnXAw":{"name":"Choice","transport_address":"172.17.0.5:9300","host":"172.17.0.5","ip":"172.17.0.5","version":"2.4.1","build":"c67dc32","http_address":"` + u.Host + `","http":{"bound_address":["[::]:9200"],"publish_address":"172.17.0.5:9200","max_content_length_in_bytes":104857600}}}}`
	var nodes elastic.NodesInfoResponse
	json.Unmarshal([]byte(data), &nodes)

	ctx.JSON(http.StatusOK, nodes)
}

func queryResult(ctx *gin.Context) {
	data := `{"code":0,"data":{"cpu":{"usage":null},"memory":{"usage":[{"metric":{"container_label_APP_ID":"work-web","container_label_VCLUSTER":"work","id":"/docker/4f84929cb252ed0c0f2d987f2f29f133395c1b26166f99562d76e64e0af6c80c","image":"192.168.1.75/library/nginx-stress:1.10","instance":"192.168.1.102:5014","job":"cadvisor","name":"mesos-05da0395-c3c4-4a76-bd6c-bfe8454f7244-S2.104a009e-c354-4bd3-8497-8e4d3212e3cf"},"values":[[1.481853425e+09,"0.02099609375"]]}],"usage_bytes":null,"total_bytes":null},"network":{"receive":null,"transmit":null},"filesystem":{"read":null,"write":null}}}`
	var result map[string]interface{}
	json.Unmarshal([]byte(data), &result)
	ctx.JSON(http.StatusOK, result)
}

func queryAlerts(ctx *gin.Context) {
	data := `{"status": "success","data":[]}`
	var result map[string]interface{}
	json.Unmarshal([]byte(data), &result)
	ctx.JSON(http.StatusOK, result)
}

func queryAlertsGroups(ctx *gin.Context) {
	data := `{"status": "success","data":[]}`
	var result map[string]interface{}
	json.Unmarshal([]byte(data), &result)
	ctx.JSON(http.StatusOK, result)
}

func queryAlertsStatus(ctx *gin.Context) {
	data := `{"status": "success","data":{}}`
	var result map[string]interface{}
	json.Unmarshal([]byte(data), &result)
	ctx.JSON(http.StatusOK, result)
}

func querySliences(ctx *gin.Context) {
	data := `{"status":"success","data":{"silences":[],"totalSilences":0}}`
	var result map[string]interface{}
	json.Unmarshal([]byte(data), &result)
	ctx.JSON(http.StatusOK, result)
}

func querySlience(ctx *gin.Context) {
	data := `{"status":"success","data":{"silences":[],"totalSilences":0}}`
	var result map[string]interface{}
	json.Unmarshal([]byte(data), &result)
	ctx.JSON(http.StatusOK, result)
}

func TestMain(m *testing.M) {
	config.InitConfig("../../env_file.template")
	server = startHTTPServer()
	baseURL = server.URL
	config.GetConfig().EsURL = baseURL
	config.GetConfig().PrometheusURL = baseURL
	config.GetConfig().AlertManagerURL = baseURL
	apiserver := startAPIServer(GetSearch())
	apiURL = apiserver.URL
	defer apiserver.Close()
	mo = NewMonitor()
	ret := m.Run()
	server.Close()
	os.Exit(ret)
}

func TestPing(t *testing.T) {
	sr := startHTTPServer()
	s = GetSearch()
	config.GetConfig().EsURL = sr.URL
	baseURL = sr.URL
	se := startAPIServer(GetSearch())
	resp, err := http.Get(se.URL + "/api/v1/ping")
	if err == nil && resp.StatusCode == 200 {
		t.Log("success")
	} else {
		t.Error("faild")
	}
}

func TestClusters(t *testing.T) {
	sr := startErrorClient()
	config.GetConfig().EsURL = sr.URL
	baseURL = sr.URL
	s = GetSearch()
	se := startAPIServer(s)
	resp, err := http.Get(se.URL + "/api/v1/clusters")
	if err == nil && resp.StatusCode == 503 {
		t.Log("success")
	} else {
		t.Error("faild")
	}

	sr = startHTTPServer()
	config.GetConfig().EsURL = sr.URL
	baseURL = sr.URL
	s = GetSearch()
	se = startAPIServer(s)
	resp, err = http.Get(se.URL + "/api/v1/clusters")
	if err == nil && resp.StatusCode == 200 {
		t.Log("success")
	} else {
		t.Error("faild")
	}
}

func TestApplications(t *testing.T) {
	sr := startErrorClient()
	config.GetConfig().EsURL = sr.URL
	baseURL = sr.URL
	s = GetSearch()
	se := startAPIServer(s)
	resp, err := http.Get(se.URL + "/api/v1/clusters/test/apps")
	if err == nil && resp.StatusCode == 503 {
		t.Log("success")
	} else {
		t.Error("faild")
	}

	sr = startHTTPServer()
	config.GetConfig().EsURL = sr.URL
	baseURL = sr.URL
	s = GetSearch()
	se = startAPIServer(s)
	resp, err = http.Get(se.URL + "/api/v1/clusters/test/apps")
	if err == nil && resp.StatusCode == 200 {
		t.Log("success")
	} else {
		t.Error("faild")
	}
}

func TestSlots(t *testing.T) {
	sr := startErrorClient()
	config.GetConfig().EsURL = sr.URL
	baseURL = sr.URL
	s = GetSearch()
	se := startAPIServer(s)
	resp, err := http.Get(se.URL + "/api/v1/clusters/test/apps/test/slots")
	if err == nil && resp.StatusCode == 503 {
		t.Log("success")
	} else {
		t.Error("faild")
	}

	sr = startHTTPServer()
	config.GetConfig().EsURL = sr.URL
	baseURL = sr.URL
	s = GetSearch()
	se = startAPIServer(s)
	resp, err = http.Get(se.URL + "/api/v1/clusters/test/apps/test/slots")
	if err == nil && resp.StatusCode == 200 {
		t.Log("success")
	} else {
		t.Error("faild")
	}
}

func TestTasks(t *testing.T) {
	sr := startErrorClient()
	config.GetConfig().EsURL = sr.URL
	baseURL = sr.URL
	s = GetSearch()
	se := startAPIServer(s)
	resp, err := http.Get(se.URL + "/api/v1/clusters/test/apps/test/slots/0/tasks")
	if err == nil && resp.StatusCode == 503 {
		t.Log("success")
	} else {
		t.Error("faild")
	}

	sr = startHTTPServer()
	config.GetConfig().EsURL = sr.URL
	baseURL = sr.URL
	s = GetSearch()
	se = startAPIServer(s)
	resp, err = http.Get(se.URL + "/api/v1/clusters/test/apps/test/slots/0/tasks")
	if err == nil && resp.StatusCode == 200 {
		t.Log("success")
	} else {
		t.Error("faild")
	}
}

func TestSources(t *testing.T) {
	sr := startErrorClient()
	config.GetConfig().EsURL = sr.URL
	baseURL = sr.URL
	s = GetSearch()
	se := startAPIServer(s)
	resp, err := http.Get(se.URL + "/api/v1/clusters/test/apps/test/sources?slot=0&task=test")
	if err == nil && resp.StatusCode == 503 {
		t.Log("success")
	} else {
		t.Error("faild")
	}

	sr = startHTTPServer()
	config.GetConfig().EsURL = sr.URL
	baseURL = sr.URL
	s = GetSearch()
	se = startAPIServer(s)
	resp, err = http.Get(se.URL + "/api/v1/clusters/test/apps/test/sources?slot=0&task=test")
	if err == nil && resp.StatusCode == 200 {
		t.Log("success")
	} else {
		t.Error("faild")
	}
}

func TestSearch(t *testing.T) {
	sr := startErrorClient()
	config.GetConfig().EsURL = sr.URL
	baseURL = sr.URL
	s = GetSearch()
	se := startAPIServer(s)
	resp, err := http.Get(se.URL + "/api/v1/clusters/test/apps/test/search?slot=0&task=test")
	if err == nil && resp.StatusCode == 503 {
		t.Log("success")
	} else {
		t.Error("faild")
	}

	sr = startHTTPServer()
	config.GetConfig().EsURL = sr.URL
	baseURL = sr.URL
	s = GetSearch()
	se = startAPIServer(s)
	resp, err = http.Get(se.URL +
		"/api/v1/clusters/mola/apps/test/search?slot=0&keyword=GET&source=stderr&conj=or&task=test-task")
	if err == nil && resp.StatusCode == 200 {
		t.Log("success")
	} else {
		t.Log(err)
		t.Error("faild")
	}
}

func TestContext(t *testing.T) {
	sr := startErrorClient()
	config.GetConfig().EsURL = sr.URL
	baseURL = sr.URL
	s = GetSearch()
	se := startAPIServer(s)
	resp, err := http.Get(se.URL + "/api/v1/clusters/test/apps/test/context?slot=0&task=test")
	if err == nil && resp.StatusCode == 503 {
		t.Log("success")
	} else {
		t.Error("faild")
	}

	sr = startHTTPServer()
	config.GetConfig().EsURL = sr.URL
	baseURL = sr.URL
	s = GetSearch()
	se = startAPIServer(s)
	resp, err = http.Get(se.URL +
		"/api/v1/clusters/mola/apps/test/context?slot=0&keyword=GET&source=stderr&conj=or&task=test-task&offset=12313131310000")
	if err == nil && resp.StatusCode == 200 {
		t.Log("success")
	} else {
		t.Log(err, resp.StatusCode)
		t.Error("faild")
	}
}
