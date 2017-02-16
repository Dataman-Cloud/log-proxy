package api

import (
	"bytes"
	"container/list"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"gopkg.in/olivere/elastic.v3"

	"github.com/Dataman-Cloud/log-proxy/src/config"
	"github.com/Dataman-Cloud/log-proxy/src/models"
)

var (
	exporterURL string
	ep          *Exporter
)

func startESErrorClient() *httptest.Server {
	router := gin.New()
	router.HEAD("/", func(ctx *gin.Context) { ctx.String(200, "") })
	router.GET("/_nodes/http", nodes)
	router.POST("/:index/_search", func(ctx *gin.Context) { ctx.String(503, "error") })
	return httptest.NewServer(router)
}

func startExporterServer(exporter *Exporter) *httptest.Server {
	router := gin.New()

	exporterv1 := router.Group("/api/v1/exporter", func(ctx *gin.Context) { ctx.Set("page", models.Page{}) })
	{
		exporterv1.POST("/keyword", func(ctx *gin.Context) { exporter.CreateFilter(ctx) })
		exporterv1.GET("/keyword", func(ctx *gin.Context) { exporter.GetFilters(ctx) })
		exporterv1.GET("/keyword/:id", func(ctx *gin.Context) { exporter.GetFilter(ctx) })
		exporterv1.PUT("/keyword", func(ctx *gin.Context) { exporter.UpdateFilter(ctx) })
		exporterv1.DELETE("/keyword/:id", func(ctx *gin.Context) { exporter.DeleteFilter(ctx) })
		exporterv1.POST("/input", receiverLog)
	}
	return httptest.NewServer(router)
}

func startESHTTPServer() *httptest.Server {
	router := gin.New()
	router.HEAD("/", func(ctx *gin.Context) { ctx.String(200, "") })
	router.GET("/_nodes/http", nodes)
	//router.POST("/:index/dataman-test-web", webs)
	//router.POST("/:index/dataman-test-web/_search", task)
	//router.POST("/:index/dataman-prometheus/*action", pro)
	router.POST("/:index/dataman-keyword-filter/*action", filters)
	router.POST("/:index/_search", app)
	//router.GET("/.dataman-prometheus/dataman-prometheus/AVj3kWyMIIGpJqE63T3m", getp)
	router.GET("/.dataman-keyword-filter/dataman-keyword-filter/test", getf)
	router.DELETE("/.dataman-keyword-filter/dataman-keyword-filter/test", getf)
	//router.GET("/api/v1/query_rang", queryResult)
	//router.GET("/api/v1/alerts", queryAlerts)
	//router.GET("/api/v1/alerts/groups", queryAlertsGroups)
	//router.GET("/api/v1/alerts/status", queryAlertsStatus)
	//router.GET("/api/v1/silences", querySliences)

	return httptest.NewServer(router)
}

func getf(ctx *gin.Context) {
	data := `{"_index":"1","_type":"1","_id":"1","_source":{"test":"value"}}`
	var info elastic.GetResult
	json.Unmarshal([]byte(data), &info)

	ctx.JSON(http.StatusOK, info)
}

func nodes(ctx *gin.Context) {
	u, _ := url.Parse(config.GetConfig().EsURL)
	data := `{"cluster_name":"elasticsearch","nodes":{"Ijb_-48ZQYmEnQ0a5BnXAw":{"name":"Choice","transport_address":"172.17.0.5:9300","host":"172.17.0.5","ip":"172.17.0.5","version":"2.4.1","build":"c67dc32","http_address":"` + u.Host + `","http":{"bound_address":["[::]:9200"],"publish_address":"172.17.0.5:9200","max_content_length_in_bytes":104857600}}}}`
	var nodes elastic.NodesInfoResponse
	json.Unmarshal([]byte(data), &nodes)

	ctx.JSON(http.StatusOK, nodes)
}

func app(ctx *gin.Context) {
	data := `{"took":137,"_scroll_id":"","hits":{"total":6,"max_score":0,"hits":[]},"suggest":null,"aggregations":{"apps":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"test-web","doc_count":6}]}},"timed_out":false,"terminated_early":false,"_shards":{"total":5,"successful":5,"failed":0}}`
	var info elastic.SearchResult
	json.Unmarshal([]byte(data), &info)

	ctx.JSON(http.StatusOK, info)
}

func filters(ctx *gin.Context) {
	data := `{"took":15,"_scroll_id":"","hits":{"total":6,"max_score":0,"hits":[{"_score":0.1825316,"_index":"dataman-test-2016-12-13","_type":"dataman-test-web","_id":"AVj3kWyMIIGpJqE63T3m","_uid":"","_timestamp":0,"_ttl":600680748,"_routing":"","_parent":"","_version":null,"sort":null,"highlight":{"message":["192.168.1.98 - - [13/Dec/2016:17:33:31 +0000] \"@dataman-highlighted-field@GET@/dataman-highlighted-field@ / HTTP/1.1\" 304 0 \"-\" \"Mozilla/5.0 (Macintosh"]},"_source":{"message":"192.168.1.98 - - [13/Dec/2016:17:33:31 +0000] \"GET / HTTP/1.1\" 304 0 \"-\" \"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/47.0.2526.80 Safari/537.36\" \"-\"\n","@version":"1","@timestamp":"2016-12-13T09:44:07.282Z","host":"192.168.1.63","port":33762,"containerid":"55f3c919563f276b0566a8a2bb01167d24a7498a18a72d556fa8f630c5956958","logtime":"2016-12-14T01:33:35.421815898+08:00","path":"stdout","offset":1481650415421815800,"app":"test-web","user":"4","task":"test-web.ac4616e4-c02b-11e6-9030-024245dc84c8","groupid":"1","clusterid":"test"},"fields":null,"_explanation":null,"matched_queries":null,"inner_hits":null}]},"suggest":null,"aggregations":{"tasks":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"test-web.ac4616e4-c02b-11e6-9030-024245dc84c8","doc_count":6}]},"paths":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"test-web.ac4616e4-c02b-11e6-9030-024245dc84c8","doc_count":6}]}},"timed_out":false,"terminated_early":false,"_shards":{"total":5,"successful":5,"failed":0}}`
	var info elastic.SearchResult
	json.Unmarshal([]byte(data), &info)

	ctx.JSON(http.StatusOK, info)
}

func TestMain(m *testing.M) {
	config.InitConfig("../../../env_file.template")
	server := startESHTTPServer()
	ret := m.Run()
	server.Close()
	os.Exit(ret)
}

func TestCreateFilter(t *testing.T) {
	sr := startESErrorClient()
	config.GetConfig().EsURL = sr.URL
	ep = NewExporter()
	se := startExporterServer(ep)
	ep.KeywordFilter = make(map[string]*list.List)
	filter := models.KWFilter{
		AppID:   "test",
		Keyword: "test1",
	}
	data, _ := json.Marshal(filter)
	resp, _ := http.NewRequest("POST", se.URL+"/api/v1/exporter/keyword", bytes.NewBuffer(data))
	req, err := http.DefaultClient.Do(resp)
	if err == nil && req.StatusCode == 503 {
		t.Log("success")
	} else {
		t.Error("failed")
	}

	sr = startESHTTPServer()
	config.GetConfig().EsURL = sr.URL
	resp, _ = http.NewRequest("POST", se.URL+"/api/v1/exporter/keyword", nil)
	req, err = http.DefaultClient.Do(resp)
	if err == nil && req.StatusCode == 400 {
		t.Log("success")
	} else {
		t.Error("failed")
	}

	filter = models.KWFilter{}
	data, _ = json.Marshal(filter)
	resp, _ = http.NewRequest("POST", se.URL+"/api/v1/exporter/keyword", bytes.NewBuffer(data))
	req, err = http.DefaultClient.Do(resp)
	if err == nil && req.StatusCode == 400 {
		t.Log("success")
	} else {
		t.Error("failed")
	}

	filter = models.KWFilter{
		AppID: "test",
	}
	data, _ = json.Marshal(filter)
	resp, _ = http.NewRequest("POST", se.URL+"/api/v1/exporter/keyword", bytes.NewBuffer(data))
	req, err = http.DefaultClient.Do(resp)
	if err == nil && req.StatusCode == 400 {
		t.Log("success")
	} else {
		t.Error("failed")
	}

	l := list.New()
	l.PushFront("test")
	ep.KeywordFilter = map[string]*list.List{
		"test": l,
	}
	filter = models.KWFilter{
		AppID:   "test",
		Keyword: "test",
	}
	data, _ = json.Marshal(filter)
	resp, _ = http.NewRequest("POST", se.URL+"/api/v1/exporter/keyword", bytes.NewBuffer(data))
	req, err = http.DefaultClient.Do(resp)
	if err == nil && req.StatusCode == 503 {
		t.Log("success")
	} else {
		t.Error("failed")
	}

	ep.KeywordFilter = make(map[string]*list.List)
	ep.KeywordFilter["test"] = list.New()
	filter = models.KWFilter{
		AppID:   "test",
		Keyword: "test",
	}
	data, _ = json.Marshal(filter)
	resp, _ = http.NewRequest("POST", se.URL+"/api/v1/exporter/keyword", bytes.NewBuffer(data))
	req, err = http.DefaultClient.Do(resp)
	if err == nil && req.StatusCode == 200 {
		t.Log("success")
	} else {
		t.Error("failed")
	}
}

func TestDeleteFilter(t *testing.T) {
	sr := startESHTTPServer()
	config.GetConfig().EsURL = sr.URL
	ep = NewExporter()
	se := startExporterServer(ep)
	l := list.New()
	l.PushFront("test")
	l.PushFront("")
	ep.KeywordFilter[""] = l
	req, _ := http.NewRequest("DELETE", se.URL+"/api/v1/exporter/keyword/test", nil)
	resp, err := http.DefaultClient.Do(req)

	if err == nil && resp.StatusCode == 200 {
		t.Log("success")
	} else {
		t.Error("failed")
	}

}

func TestGetFilters(t *testing.T) {
	sr := startESHTTPServer()
	config.GetConfig().EsURL = sr.URL
	ep = NewExporter()
	se := startExporterServer(ep)

	resp, err := http.Get(se.URL + "/api/v1/exporter/keyword")
	if err == nil && resp.StatusCode == 200 {
		t.Log("success")
	} else {
		t.Error("failed")
	}
}

func TestGetFilter(t *testing.T) {
	sr := startESHTTPServer()
	config.GetConfig().EsURL = sr.URL
	ep = NewExporter()
	se := startExporterServer(ep)

	resp, err := http.Get(se.URL + "/api/v1/exporter/keyword/test")
	if err == nil && resp.StatusCode == 200 {
		t.Log("success")
	} else {
		t.Error("failed")
	}
}

func TestUpdateFilter(t *testing.T) {
	sr := startESHTTPServer()
	config.GetConfig().EsURL = sr.URL
	ep = NewExporter()
	se := startExporterServer(ep)

	req, _ := http.NewRequest("PUT", se.URL+"/api/v1/exporter/keyword", nil)
	resp, err := http.DefaultClient.Do(req)
	if err == nil && resp.StatusCode == 400 {
		t.Log("success")
	} else {
		t.Error("failed")
	}

	filter := models.KWFilter{}
	data, _ := json.Marshal(filter)
	req, _ = http.NewRequest("PUT", se.URL+"/api/v1/exporter/keyword", bytes.NewBuffer(data))
	resp, err = http.DefaultClient.Do(req)
	if err == nil && resp.StatusCode == 400 {
		t.Log("success")
	} else {
		t.Error("failed")
	}

	filter = models.KWFilter{
		ID: "test",
	}
	data, _ = json.Marshal(filter)
	req, _ = http.NewRequest("PUT", se.URL+"/api/v1/exporter/keyword", bytes.NewBuffer(data))
	resp, err = http.DefaultClient.Do(req)
	if err == nil && resp.StatusCode == 400 {
		t.Log("success")
	} else {
		t.Error("failed")
	}

	filter = models.KWFilter{
		ID:    "test",
		AppID: "appid",
	}
	data, _ = json.Marshal(filter)
	req, _ = http.NewRequest("PUT", se.URL+"/api/v1/exporter/keyword", bytes.NewBuffer(data))
	resp, err = http.DefaultClient.Do(req)
	if err == nil && resp.StatusCode == 400 {
		t.Log("success")
	} else {
		t.Error("failed")
	}

	ep.KeywordFilter = make(map[string]*list.List)
	ep.KeywordFilter[""] = list.New()
	filter = models.KWFilter{
		ID:      "test",
		AppID:   "appid",
		Keyword: "test",
	}
	data, _ = json.Marshal(filter)
	req, _ = http.NewRequest("PUT", se.URL+"/api/v1/exporter/keyword", bytes.NewBuffer(data))
	resp, err = http.DefaultClient.Do(req)
	if err == nil && resp.StatusCode == 200 {
		t.Log("success")
	} else {
		t.Error("faild")
	}

	l := list.New()
	l.PushFront("test")
	ep.KeywordFilter[""] = l
	req, _ = http.NewRequest("PUT", se.URL+"/api/v1/exporter/keyword", bytes.NewBuffer(data))
	resp, err = http.DefaultClient.Do(req)
	if err == nil && resp.StatusCode == 200 {
		t.Log("success")
	} else {
		t.Error("failed")
	}
}

func receiverLog(ctx *gin.Context) {
	ep.ReceiveLog(ctx)
}

func TestReceiveLog(t *testing.T) {
	if ep == nil {
		ep = NewExporter()
	}
	apiserver := startExporterServer(ep)
	exporterURL = apiserver.URL
	defer apiserver.Close()

	m := make(map[string]interface{})
	data, _ := json.Marshal(m)

	req, _ := http.NewRequest("POST", exporterURL+"/api/v1/exporter/input", bytes.NewBuffer(data))
	resp, err := http.DefaultClient.Do(req)
	if err == nil && resp.StatusCode == 503 {
		t.Log("success")
	} else {
		t.Error("failed")
	}

	m = map[string]interface{}{
		"app": "test",
	}

	data, _ = json.Marshal(m)
	req, _ = http.NewRequest("POST", exporterURL+"/api/v1/exporter/input", bytes.NewBuffer(data))
	resp, err = http.DefaultClient.Do(req)
	if err == nil && resp.StatusCode == 503 {
		t.Log("success")
	} else {
		t.Error("failed")
	}

	m = map[string]interface{}{
		"app":  "test",
		"task": "test",
	}

	data, _ = json.Marshal(m)
	req, _ = http.NewRequest("POST", exporterURL+"/api/v1/exporter/input", bytes.NewBuffer(data))
	resp, err = http.DefaultClient.Do(req)
	if err == nil && resp.StatusCode == 503 {
		t.Log("success")
	} else {
		t.Error("failed")
	}

	m = map[string]interface{}{
		"app":  "test",
		"task": "test",
		"path": "test",
	}

	data, _ = json.Marshal(m)
	req, _ = http.NewRequest("POST", exporterURL+"/api/v1/exporter/input", bytes.NewBuffer(data))
	resp, err = http.DefaultClient.Do(req)
	if err == nil && resp.StatusCode == 503 {
		t.Log("success")
	} else {
		t.Error("failed")
	}

	m = map[string]interface{}{
		"app":  "test",
		"task": "test",
		"path": "test",
		"user": "test",
	}

	data, _ = json.Marshal(m)
	req, _ = http.NewRequest("POST", exporterURL+"/api/v1/exporter/input", bytes.NewBuffer(data))
	resp, err = http.DefaultClient.Do(req)
	if err == nil && resp.StatusCode == 503 {
		t.Log("success")
	} else {
		t.Error("failed")
	}

	m = map[string]interface{}{
		"app":     "test",
		"task":    "test",
		"path":    "test",
		"user":    "test",
		"cluster": "test",
	}

	data, _ = json.Marshal(m)
	req, _ = http.NewRequest("POST", exporterURL+"/api/v1/exporter/input", bytes.NewBuffer(data))
	resp, err = http.DefaultClient.Do(req)
	if err == nil && resp.StatusCode == 503 {
		t.Log("success")
	} else {
		t.Error("failed")
	}

	m = map[string]interface{}{
		"app":     "test",
		"task":    "test",
		"path":    "test",
		"user":    "test",
		"cluster": "test",
		"offset":  111,
	}

	data, _ = json.Marshal(m)
	req, _ = http.NewRequest("POST", exporterURL+"/api/v1/exporter/input", bytes.NewBuffer(data))
	resp, err = http.DefaultClient.Do(req)
	if err == nil && resp.StatusCode == 503 {
		t.Log("success")
	} else {
		t.Error("failed")
	}

	m = map[string]interface{}{
		"app":     "test",
		"task":    "test",
		"path":    "test",
		"user":    "test",
		"cluster": "test",
		"offset":  111,
		"message": "get",
	}

	data, _ = json.Marshal(m)
	req, _ = http.NewRequest("POST", exporterURL+"/api/v1/exporter/input", bytes.NewBuffer(data))
	resp, err = http.DefaultClient.Do(req)
	if err == nil && resp.StatusCode == 200 {
		t.Log("success")
	} else {
		t.Error("failed")
	}

	m = map[string]interface{}{
		"app":     "test",
		"task":    "test",
		"path":    "test",
		"user":    "test",
		"cluster": "test",
		"offset":  111,
		"message": "get",
	}

	l := list.New()
	l.PushFront("get")
	ep.KeywordFilter["testtest"] = l

	data, _ = json.Marshal(m)
	req, _ = http.NewRequest("POST", exporterURL+"/api/v1/exporter/input", bytes.NewBuffer(data))
	resp, err = http.DefaultClient.Do(req)
	if err == nil && resp.StatusCode == 200 {
		t.Log("success")
	} else {
		t.Error("failed")
	}
}
