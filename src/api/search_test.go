package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"

	"github.com/Dataman-Cloud/log-proxy/src/config"
	"github.com/Dataman-Cloud/log-proxy/src/models"
	"github.com/gin-gonic/gin"
	"gopkg.in/olivere/elastic.v3"
)

var (
	apiUrl  string
	baseUrl string
	server  *httptest.Server
	s       *search
	mo      *monitor
)

func startApiServer() *httptest.Server {
	router := gin.New()
	v1 := router.Group("/api/v1", func(ctx *gin.Context) { ctx.Set("page", models.Page{}) })
	{
		v1.GET("/ping", func(ctx *gin.Context) { s.Ping(ctx) })
		v1.GET("/applications", func(ctx *gin.Context) {
			s.Applications(ctx)
		})
		v1.GET("/tasks/:appid", func(ctx *gin.Context) { s.Tasks(ctx) })
		v1.GET("/paths/:appid", func(ctx *gin.Context) { s.Paths(ctx) })
		v1.GET("/index", func(ctx *gin.Context) { s.Index(ctx) })
		v1.GET("/context", func(ctx *gin.Context) { s.Context(ctx) })
	}

	vr := router.Group("/v1/recive")
	{
		vr.POST("/prometheus", receiver)
		vr.POST("/log", receiverlog)
	}

	v1m := router.Group("/api/v1/monitor", func(ctx *gin.Context) { ctx.Set("page", models.Page{}) })
	{
		v1m.POST("/alert", createAlert)
		v1m.GET("/alert", getAlerts)
		v1m.GET("/alert/:id", getAlert)
		v1m.PUT("/alert", updateAlert)
		v1m.DELETE("/alert/:id", deleteAlert)
		v1m.GET("/prometheus", getprometheus)
		v1m.GET("/prometheus/:id", getprometheu)
		v1m.GET("/query", getQuery)
		v1m.GET("/info", getQueryInfo)
		v1m.GET("/nodes", getQueryNodes)
		v1m.GET("/alerts", getMonitorAlerts)
		v1m.GET("/alerts/groups", getMonitorAlertsGroups)
		v1m.GET("/alerts/status", getMonitorAlertsStatus)
		v1m.GET("/silences", getMonitorSilences)
		v1m.GET("/silence", getMonitorSilence)
		v1m.POST("/silence", createMonitorSilence)
		v1m.PUT("/silence", updateMonitorSilence)
		v1m.DELETE("/silence/:id", deleteMonitorSilence)
	}
	return httptest.NewServer(router)
}

func startHttpServer() *httptest.Server {
	router := gin.New()
	router.HEAD("/", func(ctx *gin.Context) { ctx.String(200, "") })
	router.GET("/_nodes/http", nodes)
	router.POST("/:index/dataman-test-web", webs)
	router.POST("/:index/dataman-test-web/_search", task)
	router.POST("/:index/dataman-prometheus/*action", pro)
	router.POST("/:index/dataman-alerts/*action", alerts)
	router.POST("/:index/_search", app)
	router.GET("/.dataman-prometheus/dataman-prometheus/AVj3kWyMIIGpJqE63T3m", getp)
	router.GET("/.dataman-alerts/dataman-alerts/test", getp)
	router.DELETE("/.dataman-alerts/dataman-alerts/test", getp)
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
	data := `{"took":137,"_scroll_id":"","hits":{"total":6,"max_score":0,"hits":[]},"suggest":null,"aggregations":{"apps":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"test-web","doc_count":6}]}},"timed_out":false,"terminated_early":false,"_shards":{"total":5,"successful":5,"failed":0}}`
	var info elastic.SearchResult
	json.Unmarshal([]byte(data), &info)

	ctx.JSON(http.StatusOK, info)
}

func webs(ctx *gin.Context) {
	data := `{"took":15,"_scroll_id":"","hits":{"total":6,"max_score":0,"hits":[{"_score":0.1825316,"_index":"dataman-test-2016-12-13","_type":"dataman-test-web","_id":"AVj3kWyMIIGpJqE63T3m","_uid":"","_timestamp":0,"_ttl":600680748,"_routing":"","_parent":"","_version":null,"sort":null,"highlight":{"message":["192.168.1.98 - - [13/Dec/2016:17:33:31 +0000] \"@dataman-highlighted-field@GET@/dataman-highlighted-field@ / HTTP/1.1\" 304 0 \"-\" \"Mozilla/5.0 (Macintosh"]},"_source":{"message":"192.168.1.98 - - [13/Dec/2016:17:33:31 +0000] \"GET / HTTP/1.1\" 304 0 \"-\" \"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/47.0.2526.80 Safari/537.36\" \"-\"\n","@version":"1","@timestamp":"2016-12-13T09:44:07.282Z","host":"192.168.1.63","port":33762,"containerid":"55f3c919563f276b0566a8a2bb01167d24a7498a18a72d556fa8f630c5956958","logtime":"2016-12-14T01:33:35.421815898+08:00","path":"stdout","offset":1481650415421815800,"appid":"test-web","userid":"4","taskid":"test-web.ac4616e4-c02b-11e6-9030-024245dc84c8","groupid":"1","clusterid":"test"},"fields":null,"_explanation":null,"matched_queries":null,"inner_hits":null}]},"suggest":null,"aggregations":{"tasks":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"test-web.ac4616e4-c02b-11e6-9030-024245dc84c8","doc_count":6}]},"paths":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"test-web.ac4616e4-c02b-11e6-9030-024245dc84c8","doc_count":6}]}},"timed_out":false,"terminated_early":false,"_shards":{"total":5,"successful":5,"failed":0}}`
	var info elastic.SearchResult
	json.Unmarshal([]byte(data), &info)

	ctx.JSON(http.StatusOK, info)
}

func alerts(ctx *gin.Context) {
	data := `{"took":15,"_scroll_id":"","hits":{"total":6,"max_score":0,"hits":[{"_score":0.1825316,"_index":"dataman-test-2016-12-13","_type":"dataman-test-web","_id":"AVj3kWyMIIGpJqE63T3m","_uid":"","_timestamp":0,"_ttl":600680748,"_routing":"","_parent":"","_version":null,"sort":null,"highlight":{"message":["192.168.1.98 - - [13/Dec/2016:17:33:31 +0000] \"@dataman-highlighted-field@GET@/dataman-highlighted-field@ / HTTP/1.1\" 304 0 \"-\" \"Mozilla/5.0 (Macintosh"]},"_source":{"message":"192.168.1.98 - - [13/Dec/2016:17:33:31 +0000] \"GET / HTTP/1.1\" 304 0 \"-\" \"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/47.0.2526.80 Safari/537.36\" \"-\"\n","@version":"1","@timestamp":"2016-12-13T09:44:07.282Z","host":"192.168.1.63","port":33762,"containerid":"55f3c919563f276b0566a8a2bb01167d24a7498a18a72d556fa8f630c5956958","logtime":"2016-12-14T01:33:35.421815898+08:00","path":"stdout","offset":1481650415421815800,"appid":"test-web","userid":"4","taskid":"test-web.ac4616e4-c02b-11e6-9030-024245dc84c8","groupid":"1","clusterid":"test"},"fields":null,"_explanation":null,"matched_queries":null,"inner_hits":null}]},"suggest":null,"aggregations":{"tasks":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"test-web.ac4616e4-c02b-11e6-9030-024245dc84c8","doc_count":6}]},"paths":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"test-web.ac4616e4-c02b-11e6-9030-024245dc84c8","doc_count":6}]}},"timed_out":false,"terminated_early":false,"_shards":{"total":5,"successful":5,"failed":0}}`
	var info elastic.SearchResult
	json.Unmarshal([]byte(data), &info)

	ctx.JSON(http.StatusOK, info)
}

func pro(ctx *gin.Context) {
	data := `{"took":15,"_scroll_id":"","hits":{"total":6,"max_score":0,"hits":[{"_score":0.1825316,"_index":"dataman-test-2016-12-13","_type":"dataman-test-web","_id":"AVj3kWyMIIGpJqE63T3m","_uid":"","_timestamp":0,"_ttl":600680748,"_routing":"","_parent":"","_version":null,"sort":null,"highlight":{"message":["192.168.1.98 - - [13/Dec/2016:17:33:31 +0000] \"@dataman-highlighted-field@GET@/dataman-highlighted-field@ / HTTP/1.1\" 304 0 \"-\" \"Mozilla/5.0 (Macintosh"]},"_source":{"message":"192.168.1.98 - - [13/Dec/2016:17:33:31 +0000] \"GET / HTTP/1.1\" 304 0 \"-\" \"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/47.0.2526.80 Safari/537.36\" \"-\"\n","@version":"1","@timestamp":"2016-12-13T09:44:07.282Z","host":"192.168.1.63","port":33762,"containerid":"55f3c919563f276b0566a8a2bb01167d24a7498a18a72d556fa8f630c5956958","logtime":"2016-12-14T01:33:35.421815898+08:00","path":"stdout","offset":1481650415421815800,"appid":"test-web","userid":"4","taskid":"test-web.ac4616e4-c02b-11e6-9030-024245dc84c8","groupid":"1","clusterid":"test"},"fields":null,"_explanation":null,"matched_queries":null,"inner_hits":null}]},"suggest":null,"aggregations":{"tasks":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"test-web.ac4616e4-c02b-11e6-9030-024245dc84c8","doc_count":6}]},"paths":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"test-web.ac4616e4-c02b-11e6-9030-024245dc84c8","doc_count":6}]}},"timed_out":false,"terminated_early":false,"_shards":{"total":5,"successful":5,"failed":0}}`
	var info elastic.SearchResult
	json.Unmarshal([]byte(data), &info)

	ctx.JSON(http.StatusOK, info)
}

func task(ctx *gin.Context) {
	data := `{"took":15,"_scroll_id":"","hits":{"total":6,"max_score":0,"hits":[{"_score":0.1825316,"_index":"dataman-test-2016-12-13","_type":"dataman-test-web","_id":"AVj3kWyMIIGpJqE63T3m","_uid":"","_timestamp":0,"_ttl":600680748,"_routing":"","_parent":"","_version":null,"sort":null,"highlight":{"message":["192.168.1.98 - - [13/Dec/2016:17:33:31 +0000] \"@dataman-highlighted-field@GET@/dataman-highlighted-field@ / HTTP/1.1\" 304 0 \"-\" \"Mozilla/5.0 (Macintosh"]},"_source":{"message":"192.168.1.98 - - [13/Dec/2016:17:33:31 +0000] \"GET / HTTP/1.1\" 304 0 \"-\" \"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/47.0.2526.80 Safari/537.36\" \"-\"\n","@version":"1","@timestamp":"2016-12-13T09:44:07.282Z","host":"192.168.1.63","port":33762,"containerid":"55f3c919563f276b0566a8a2bb01167d24a7498a18a72d556fa8f630c5956958","logtime":"2016-12-14T01:33:35.421815898+08:00","path":"stdout","offset":1481650415421815800,"appid":"test-web","userid":"4","taskid":"test-web.ac4616e4-c02b-11e6-9030-024245dc84c8","groupid":"1","clusterid":"test"},"fields":null,"_explanation":null,"matched_queries":null,"inner_hits":null}]},"suggest":null,"aggregations":{"tasks":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"test-web.ac4616e4-c02b-11e6-9030-024245dc84c8","doc_count":6}]},"paths":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"test-web.ac4616e4-c02b-11e6-9030-024245dc84c8","doc_count":6}]}},"timed_out":false,"terminated_early":false,"_shards":{"total":5,"successful":5,"failed":0}}`
	var info elastic.SearchResult
	json.Unmarshal([]byte(data), &info)
	var m map[string]interface{}
	json.NewDecoder(ctx.Request.Body).Decode(&m)
	defer ctx.Request.Body.Close()
	ss, _ := json.Marshal(m)
	query := `{"from":0,"highlight":{"fields":{"message":{}},"post_tags":["@/dataman-highlighted-field@"],"pre_tags":["@dataman-highlighted-field@"]},"query":{"bool":{"filter":{"range":{"logtime.timestamp":{"format":"epoch_millis","from":null,"include_lower":true,"include_upper":true,"to":null}}}}},"size":0}`
	if string(ss) == query {
		ctx.JSON(http.StatusServiceUnavailable, info)
		return
	}
	ctx.JSON(http.StatusOK, info)

}

func nodes(ctx *gin.Context) {
	u, _ := url.Parse(baseUrl)
	data := `{"cluster_name":"elasticsearch","nodes":{"Ijb_-48ZQYmEnQ0a5BnXAw":{"name":"Choice","transport_address":"172.17.0.5:9300","host":"172.17.0.5","ip":"172.17.0.5","version":"2.4.1","build":"c67dc32","http_address":"` + u.Host + `","http":{"bound_address":["[::]:9200"],"publish_address":"172.17.0.5:9200","max_content_length_in_bytes":104857600}}}}`
	var nodes elastic.NodesInfoResponse
	json.Unmarshal([]byte(data), &nodes)

	ctx.JSON(http.StatusOK, nodes)
}

func queryResult(ctx *gin.Context) {
	data := `{"code":0,"data":{"cpu":{"usage":null},"memory":{"usage":[{"metric":{"container_label_APP_ID":"work-web","container_label_VCLUSTER":"work","id":"/docker/4f84929cb252ed0c0f2d987f2f29f133395c1b26166f99562d76e64e0af6c80c","image":"192.168.1.75/library/nginx-stress:1.10","instance":"192.168.1.102:5014","job":"cadvisor","name":"mesos-05da0395-c3c4-4a76-bd6c-bfe8454f7244-S2.104a009e-c354-4bd3-8497-8e4d3212e3cf"},"values":[[1.481853425e+09,"0.02099609375"]]}],"usage_bytes":null,"total_bytes":null},"network":{"receive":null,"transmit":null},"filesystem":{"read":null,"write":null}}}`
	var result map[string]interface{}
	json.Unmarshal([]byte(data), &result)
	fmt.Printf("%v\n", result)
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
	fmt.Printf("%v\n", result)
	ctx.JSON(http.StatusOK, result)
}

func querySlience(ctx *gin.Context) {
	data := `{"status":"success","data":{"silences":[],"totalSilences":0}}`
	var result map[string]interface{}
	json.Unmarshal([]byte(data), &result)
	fmt.Printf("%v\n", result)
	ctx.JSON(http.StatusOK, result)
}

func TestMain(m *testing.M) {
	apiserver := startApiServer()
	apiUrl = apiserver.URL
	defer apiserver.Close()
	config.InitConfig("../../env_file.template")
	server = startHttpServer()
	baseUrl = server.URL
	config.GetConfig().ES_URL = baseUrl
	config.GetConfig().PROMETHEUS_URL = baseUrl
	config.GetConfig().ALERTMANAGER_URL = baseUrl
	mo = NewMonitor()
	ret := m.Run()
	server.Close()
	os.Exit(ret)
}

func TestPing(t *testing.T) {
	if s == nil {
		s = GetSearch()
	}
	_, err := http.Get(apiUrl + "/api/v1/ping")
	if err == nil {
		t.Log("success")
	} else {
		t.Error("faild")
	}
}

func TestApplications(t *testing.T) {
	if s == nil {
		s = GetSearch()
	}
	_, err := http.Get(apiUrl + "/api/v1/applications")
	if err == nil {
		t.Log("success")
	} else {
		t.Error("faild")
	}

}

func TestTasks(t *testing.T) {
	if s == nil {
		s = GetSearch()
	}
	resp, err := http.Get(apiUrl + "/api/v1/tasks/test")
	if err == nil && resp.StatusCode == 200 {
		t.Log("success")
	} else {
		t.Error("faild")
	}
}

func TestPaths(t *testing.T) {
	if s == nil {
		s = GetSearch()
	}
	resp, err := http.Get(apiUrl + "/api/v1/paths/test")
	if err == nil && resp.StatusCode == 200 {
		t.Log("success")
	} else {
		t.Error("faild")
	}
}

func TestIndex(t *testing.T) {
	if s == nil {
		s = GetSearch()
	}
	resp, err := http.Get(apiUrl + "/api/v1/index")
	if err == nil && resp.StatusCode == 400 {
		t.Log("success")
	} else {
		t.Error("faild")
	}

	resp, err = http.Get(apiUrl + "/api/v1/index?appid=test-web")
	if err == nil && resp.StatusCode == 503 {
		t.Log("success")
	} else {
		t.Error("faild")
	}

	resp, err = http.Get(apiUrl + "/api/v1/index?appid=test-web&keyword=test")
	if err == nil && resp.StatusCode == 200 {
		t.Log("success")
	} else {
		t.Error("faild")
	}
}

func TestContext(t *testing.T) {
	if s == nil {
		s = GetSearch()
	}
	resp, err := http.Get(apiUrl + "/api/v1/context")
	if err == nil && resp.StatusCode == 400 {
		t.Log("success")
	} else {
		t.Error("faild")
	}

	resp, err = http.Get(apiUrl + "/api/v1/context?appid=test")
	if err == nil && resp.StatusCode == 400 {
		t.Log("success")
	} else {
		t.Error("faild")
	}

	resp, err = http.Get(apiUrl + "/api/v1/context?appid=test&taskid=test")
	if err == nil && resp.StatusCode == 400 {
		t.Log("success")
	} else {
		t.Error("faild")
	}

	resp, err = http.Get(apiUrl + "/api/v1/context?appid=test&taskid=test&path=appid")
	if err == nil && resp.StatusCode == 400 {
		t.Log("success")
	} else {
		t.Error("faild")
	}

	resp, err = http.Get(apiUrl + "/api/v1/context?appid=test&taskid=test&path=appid&offset=100")
	if err == nil && resp.StatusCode == 200 {
		t.Log("success")
	} else {
		t.Error("faild")
	}
}
