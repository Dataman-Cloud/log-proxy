package service

import (
	"encoding/json"
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
	baseURL string
	server  *httptest.Server
)

func startHTTPServer() *httptest.Server {
	router := gin.New()
	router.HEAD("/", func(ctx *gin.Context) { ctx.String(200, "") })
	router.GET("/_nodes/http", nodes)
	router.POST("/:index/dataman-user-test-web", task)
	router.POST("/:index/dataman-user-test-web/_search", task)
	router.POST("/:index/dataman-prometheus/*action", task)
	router.POST("/:index/dataman-alerts/*action", task)
	router.POST("/:index/_search", app)
	router.GET("/.dataman-prometheus/dataman-prometheus/AVj3kWyMIIGpJqE63T3m", getp)
	router.GET("/.dataman-alerts/dataman-alerts/test", getp)
	router.DELETE("/.dataman-alerts/dataman-alerts/test", getp)

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
        "total":20,
        "max_score":null,
        "hits":[
            {
                "_index":"dataman-%{cluster}-2017-04-19",
                "_type":"dataman-%{user}-%{app}",
                "_id":"AVuFY41XB1jJ-pLrRrCV",
                "_score":null,
                "_ttl":517057719,
				"highlight":{
					"message":["aaaaaa","xxxxxhhh"]
				},
                "_source":{
                    "message":"xxx test aaa keyword ",
                    "@version":"1",
                    "@timestamp":"2017-04-19T08:45:45.667Z",
                    "host":"192.168.1.102",
                    "port":51208,
                    "containerid":"f888cae13d80641f0e79db5cdbaa837635f079926dd8204ec33bbb1d9ba2f325",
                    "DM_VCLUSTER":"mola",
                    "DM_CLUSTER":"datamanmesos",
                    "DM_SLOT_INDEX":"0",
                    "DM_TASK_ID":"0-2xxxxxxxx-yaoyun-datamanmesos-f2969953966e40d08a9fcf2e8ef0bd5f",
                    "DM_USER":"yaoyun",
                    "DM_USER_NAME":"yaoyun",
                    "DM_SLOT_ID":"0-2xxxxxxxx-yaoyun-datamanmesos",
                    "AAAAAAAAAA":"12344",
                    "BBBBBBBBBB":"2222222",
                    "DM_APP_ID":"2xxxxxxxx-yaoyun-datamanmesos",
                    "DM_APP_NAME":"2xxxxxxxx",
                    "DM_GROUP_NAME":"yaoyun",
                    "logtime":"2017-04-19T16:45:45.586140305+08:00",
                    "path":"stdout",
                    "offset":1492591545586140400
                },
                "fields":{
                    "logtime.timestamp":[
                        1492591545586
                    ],
                    "@timestamp":[
                        1492591545667
                    ]
                },
                "sort":[
                    1492591545667
                ]
            }
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
        },
        "clusters":{
            "doc_count_error_upper_bound":0,
            "sum_other_doc_count":0,
            "buckets":[
                {
                    "key":"test-web",
                    "doc_count":6
                }
            ]
        },
        "slots":{
            "doc_count_error_upper_bound":0,
            "sum_other_doc_count":0,
            "buckets":[
                {
                    "key":"test-web",
                    "doc_count":6
                }
            ]
        },
        "tasks":{
            "doc_count_error_upper_bound":0,
            "sum_other_doc_count":0,
            "buckets":[
                {
                    "key":"test-web",
                    "doc_count":6
                }
            ]
        },
        "sources":{
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

func task(ctx *gin.Context) {
	data := `{"took":15,"_scroll_id":"","hits":{"total":6,"max_score":0,"hits":[{"_score":0.1825316,"_index":"dataman-test-2016-12-13","_type":"dataman-test-web","_id":"AVj3kWyMIIGpJqE63T3m","_uid":"","_timestamp":0,"_ttl":600680748,"_routing":"","_parent":"","_version":null,"sort":null,"highlight":{"message":["192.168.1.98 - - [13/Dec/2016:17:33:31 +0000] \"@dataman-highlighted-field@GET@/dataman-highlighted-field@ / HTTP/1.1\" 304 0 \"-\" \"Mozilla/5.0 (Macintosh"]},"_source":{"message":"192.168.1.98 - - [13/Dec/2016:17:33:31 +0000] \"GET / HTTP/1.1\" 304 0 \"-\" \"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/47.0.2526.80 Safari/537.36\" \"-\"\n","@version":"1","@timestamp":"2016-12-13T09:44:07.282Z","host":"192.168.1.63","port":33762,"containerid":"55f3c919563f276b0566a8a2bb01167d24a7498a18a72d556fa8f630c5956958","logtime":"2016-12-14T01:33:35.421815898+08:00","path":"stdout","offset":1481650415421815800,"app":"test-web","user":"4","task":"test-web.ac4616e4-c02b-11e6-9030-024245dc84c8","groupid":"1","cluster":"test"},"fields":null,"_explanation":null,"matched_queries":null,"inner_hits":null}]},"suggest":null,"aggregations":{"tasks":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"test-web.ac4616e4-c02b-11e6-9030-024245dc84c8","doc_count":6}]},"paths":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"test-web.ac4616e4-c02b-11e6-9030-024245dc84c8","doc_count":6}]}},"timed_out":false,"terminated_early":false,"_shards":{"total":5,"successful":5,"failed":0}}`
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

func TestMain(m *testing.M) {
	config.InitConfig("../../env_file.template")
	server = startHTTPServer()
	baseURL = server.URL
	config.GetConfig().EsURL = baseURL
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
	if NewEsService([]string{baseURL}) != nil {
		t.Log("success")
	} else {
		t.Error("faild")
	}
}

func TestClusters(t *testing.T) {
	service := NewEsService([]string{baseURL})
	service.Clusters(models.Page{})
}

func TestApplications(t *testing.T) {
	service := NewEsService([]string{baseURL})
	service.Applications("test", models.Page{})
}

func TestSlots(t *testing.T) {
	service := NewEsService([]string{baseURL})
	service.Slots("test", "test", models.Page{})
}

func TestTasks(t *testing.T) {
	service := NewEsService([]string{baseURL})
	service.Tasks("test-web", "test", "test", models.Page{})
}

func TestSources(t *testing.T) {
	service := NewEsService([]string{baseURL})
	opts := make(map[string]interface{})
	opts["slot"] = "0"
	opts["task"] = "test"
	opts["page"] = models.Page{}
	service.Sources("cluster", "app", opts)
}

func TestSearch(t *testing.T) {
	service := NewEsService([]string{baseURL})
	opts := make(map[string]interface{})
	opts["slot"] = "0"
	opts["task"] = "test"
	opts["page"] = models.Page{}
	opts["keyword"] = "keyword"
	opts["conj"] = "or"
	opts["source"] = "stdout"
	service.Search("cluster", "user", opts)
}

func TestContext(t *testing.T) {
	service := NewEsService([]string{baseURL})
	service.Context("cluster", "user", "test-web", "test-web.ac4616e4-c02b-11e6-9030-024245dc84c8", "stdout", "1481650415421815800", models.Page{})
}
