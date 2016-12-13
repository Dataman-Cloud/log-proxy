package service

import (
	"encoding/json"
	//"io/ioutil"
	//"fmt"
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
	baseUrl string
	server  *httptest.Server
)

func startHttpServer() *httptest.Server {
	router := gin.New()
	router.HEAD("/", func(ctx *gin.Context) { ctx.String(200, "") })
	router.GET("/_nodes/http", nodes)
	router.POST("/:index/dataman-test-web/_search", task)
	router.POST("/:index/_search", app)

	return httptest.NewServer(router)
}

func app(ctx *gin.Context) {
	data := `{"took":137,"_scroll_id":"","hits":{"total":6,"max_score":0,"hits":[]},"suggest":null,"aggregations":{"apps":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"test-web","doc_count":6}]}},"timed_out":false,"terminated_early":false,"_shards":{"total":5,"successful":5,"failed":0}}`
	var info elastic.SearchResult
	json.Unmarshal([]byte(data), &info)

	ctx.JSON(http.StatusOK, info)
}

func task(ctx *gin.Context) {
	data := `{"took":15,"_scroll_id":"","hits":{"total":6,"max_score":0,"hits":[{"_score":0.1825316,"_index":"dataman-test-2016-12-13","_type":"dataman-test-web","_id":"AVj3kWyMIIGpJqE63T3m","_uid":"","_timestamp":0,"_ttl":600680748,"_routing":"","_parent":"","_version":null,"sort":null,"highlight":{"message":["192.168.1.98 - - [13/Dec/2016:17:33:31 +0000] \"@dataman-highlighted-field@GET@/dataman-highlighted-field@ / HTTP/1.1\" 304 0 \"-\" \"Mozilla/5.0 (Macintosh"]},"_source":{"message":"192.168.1.98 - - [13/Dec/2016:17:33:31 +0000] \"GET / HTTP/1.1\" 304 0 \"-\" \"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/47.0.2526.80 Safari/537.36\" \"-\"\n","@version":"1","@timestamp":"2016-12-13T09:44:07.282Z","host":"192.168.1.63","port":33762,"containerid":"55f3c919563f276b0566a8a2bb01167d24a7498a18a72d556fa8f630c5956958","logtime":"2016-12-14T01:33:35.421815898+08:00","path":"stdout","offset":1481650415421815800,"appid":"test-web","userid":"4","taskid":"test-web.ac4616e4-c02b-11e6-9030-024245dc84c8","groupid":"1","clusterid":"test"},"fields":null,"_explanation":null,"matched_queries":null,"inner_hits":null}]},"suggest":null,"aggregations":{"tasks":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"test-web.ac4616e4-c02b-11e6-9030-024245dc84c8","doc_count":6}]},"paths":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"test-web.ac4616e4-c02b-11e6-9030-024245dc84c8","doc_count":6}]}},"timed_out":false,"terminated_early":false,"_shards":{"total":5,"successful":5,"failed":0}}`
	var info elastic.SearchResult
	json.Unmarshal([]byte(data), &info)

	ctx.JSON(http.StatusOK, info)
}

func nodes(ctx *gin.Context) {
	u, _ := url.Parse(baseUrl)
	data := `{"cluster_name":"elasticsearch","nodes":{"Ijb_-48ZQYmEnQ0a5BnXAw":{"name":"Choice","transport_address":"172.17.0.5:9300","host":"172.17.0.5","ip":"172.17.0.5","version":"2.4.1","build":"c67dc32","http_address":"` + u.Host + `","http":{"bound_address":["[::]:9200"],"publish_address":"172.17.0.5:9200","max_content_length_in_bytes":104857600}}}}`
	var nodes elastic.NodesInfoResponse
	json.Unmarshal([]byte(data), &nodes)

	ctx.JSON(http.StatusOK, nodes)
}

func TestMain(m *testing.M) {
	server = startHttpServer()
	baseUrl = server.URL
	ret := m.Run()
	server.Close()
	os.Exit(ret)
}

func TestNewEsService(t *testing.T) {
	config.InitConfig("../../env_file.template")
	if NewEsService([]string{"localhost"}) == nil {
		t.Log("success")
	} else {
		t.Error("faild")
	}
	if NewEsService([]string{baseUrl}) != nil {
		t.Log("success")
	} else {
		t.Error("faild")
	}
}

func TestApplications(t *testing.T) {
	config.InitConfig("../../env_file.template")
	service := NewEsService([]string{baseUrl})
	service.Applications(models.Page{})
}

func TestTasks(t *testing.T) {
	config.InitConfig("../../env_file.template")
	service := NewEsService([]string{baseUrl})
	service.Tasks("test-web", models.Page{})
}

func TestPath(t *testing.T) {
	config.InitConfig("../../env_file.template")
	service := NewEsService([]string{baseUrl})
	service.Paths("test-web", "test-web.ac4616e4-c02b-11e6-9030-024245dc84c8", models.Page{})
}

func TestSearch(t *testing.T) {
	config.InitConfig("../../env_file.template")
	service := NewEsService([]string{baseUrl})
	service.Search("test-web", "test-web.ac4616e4-c02b-11e6-9030-024245dc84c8", "stdout", "GET", models.Page{})
}

func TestContext(t *testing.T) {
	config.InitConfig("../../env_file.template")
	service := NewEsService([]string{baseUrl})
	service.Context("test-web", "test-web.ac4616e4-c02b-11e6-9030-024245dc84c8", "stdout", "1481650415421815800", models.Page{})
}
