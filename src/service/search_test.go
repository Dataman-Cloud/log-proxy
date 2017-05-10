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
	"github.com/stretchr/testify/assert"

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
	router.POST("/:index/_search", app)

	return httptest.NewServer(router)
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

func nodes(ctx *gin.Context) {
	u, _ := url.Parse(baseURL)
	data := `
	{
    "cluster_name":"elasticsearch",
    "nodes":{
        "Ijb_-48ZQYmEnQ0a5BnXAw":{
            "name":"Choice",
            "transport_address":"172.17.0.5:9300",
            "host":"172.17.0.5",
            "ip":"172.17.0.5",
            "version":"2.4.1",
            "build":"c67dc32",
            "http_address":"` + u.Host + `",
            "http":{
                "bound_address":[
                    "[::]:9200"
                ],
                "publish_address":"172.17.0.5:9200",
                "max_content_length_in_bytes":104857600
            }
        }
    }
}
	`
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

func TestApplications(t *testing.T) {
	service := NewEsService([]string{baseURL})
	service.Applications(models.Page{})
}

func TestSlots(t *testing.T) {
	service := NewEsService([]string{baseURL})
	service.Slots("test", models.Page{})
}

func TestTasks(t *testing.T) {
	service := NewEsService([]string{baseURL})
	service.Tasks("test", "test", models.Page{})
}

func TestSources(t *testing.T) {
	service := NewEsService([]string{baseURL})
	opts := make(map[string]interface{})
	opts["slot"] = "0"
	opts["task"] = "test"
	opts["page"] = models.Page{}
	service.Sources("app", opts)
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
	service.Search("user", opts)
}

func TestContext(t *testing.T) {
	service := NewEsService([]string{baseURL})
	opts := make(map[string]interface{})
	opts["app"] = "test"
	opts["slot"] = "0"
	opts["task"] = "test"
	opts["offset"] = "1123123130000"
	opts["source"] = "stdout"
	_, err := service.Context(opts, models.Page{})
	assert.NoError(t, err)
}

func TestContextOffsetError(t *testing.T) {
	service := NewEsService([]string{baseURL})
	opts := make(map[string]interface{})
	opts["app"] = "test"
	opts["slot"] = "0"
	opts["task"] = "test"
	opts["source"] = "stdout"
	_, err := service.Context(opts, models.Page{})
	assert.Error(t, err)

	opts["offset"] = "a11111"
	_, err = service.Context(opts, models.Page{})
	assert.Error(t, err)
}
