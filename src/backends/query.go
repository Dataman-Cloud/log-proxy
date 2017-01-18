package backends

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/Dataman-Cloud/log-proxy/src/models"

	log "github.com/Sirupsen/logrus"
)

const (
	// QUERYRANGEPATH define the query_range path of prometheus
	QUERYRANGEPATH = "/api/v1/query_range"
	// QUERYPATH define the query path of prometheus
	QUERYPATH = "/api/v1/query"
	//OKRESULT define the value of status field in response
	OKRESULT = "success"
)

// Query define the fields required by query from prometheus
type Query struct {
	HTTPClient *http.Client
	PromServer string
	Path       string
	*QueryParameter
}

// QueryParameter define the fields of query paramter
type QueryParameter struct {
	Metric  string
	Cluster string
	App     string
	Task    string
	Slot    string
	User    string
	Node    string
	Start   string
	End     string
	Step    string
	Period  string
	Expr    string
}

// QueryExpr return the result by calling function getExprResponse
func (query *Query) QueryExpr() (map[string]interface{}, error) {
	response, request, err := query.getExprResponse()
	if err != nil {
		return nil, err
	}
	var result map[string]interface{}
	err = json.Unmarshal(response, &result)
	if err != nil {
		err = fmt.Errorf("Failed to parse the response from %s", request)
		return nil, err
	}
	return result, nil
}

// getExprResponse set the request url/query parameter and then
// return the result by query from Prometheus
func (query *Query) getExprResponse() ([]byte, string, error) {
	u, err := url.Parse(query.PromServer)
	if err != nil {
		return nil, "", err
	}
	u.Path = strings.TrimRight(u.Path, "/") + query.Path
	q := u.Query()
	q.Set("query", query.Expr)
	if query.Path == QUERYRANGEPATH {
		q.Set("start", query.Start)
		q.Set("end", query.End)
		if query.Step == "" {
			q.Set("step", "30s")
		} else {
			q.Set("step", query.Step)
		}
	} else if query.Path == QUERYPATH {
	}
	u.RawQuery = q.Encode()
	resp, err := query.HTTPClient.Get(u.String())
	if err != nil {
		err = fmt.Errorf("Failed to get response from %s", u.String())
		return nil, u.String(), err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, u.String(), err
	}

	log.Printf("Get the prometheus qurey result by url: %s", u.String())
	return body, u.String(), nil
}

// QueryMetric return the result by calling function getExprResponse
func (query *Query) QueryMetric() (*models.QueryRangeResult, error) {
	// if start/end as "", set them as now
	start, end := timeRange(query.Start, query.End)
	query.Start = start
	query.End = end
	query.Expr = query.getQueryMetricExpr("task")

	response, request, err := query.getExprResponse()
	if err != nil {
		return nil, err
	}

	var result *models.QueryRangeResult
	err = json.Unmarshal(response, &result)
	if err != nil {
		err = fmt.Errorf("Failed to parse the response from %s", request)
		return nil, err
	}

	if result.Status != OKRESULT {
		err := fmt.Errorf("%s", result.Error)
		return nil, err
	}

	return result, nil
}

// getQueryMetricExpr return the expr for query.Metric
func (query *Query) getQueryMetricExpr(level string) (expr string) {
	// the level string declare the query precision, app, id and so on
	me := GetExpr(query, level)

	switch query.Metric {
	case "cpu":
		expr = me.CPU.Usage
	case "memory":
		expr = me.Memory.Percentage
	case "memory_usage":
		expr = me.Memory.Usage
	case "memory_total":
		expr = me.Memory.Total
	case "network_rx":
		expr = me.Network.Receive
	case "network_tx":
		expr = me.Network.Transmit
	case "fs_read":
		expr = me.Filesystem.Read
	case "fs_write":
		expr = me.Filesystem.Write
	case "fs_usage":
		expr = me.Filesystem.Usage
	case "fs_limit":
		expr = me.Filesystem.Limit
	}

	//fmt.Printf("expr: %s", expr)
	return expr
}

func timeRange(start, end string) (string, string) {
	const (
		OFFSET = "0m"
	)
	if start == "" && end == "" {
		endtime := time.Now()
		starttime := timeOffset(endtime, OFFSET)
		return timetoString(endtime), timetoString(starttime)
	}
	return start, end
}

func timetoString(t time.Time) string {
	return strconv.FormatInt(t.Unix(), 10)
}

func timeOffset(t time.Time, offset string) time.Time {
	duration, _ := time.ParseDuration(offset)
	return t.Add(duration)
}

// QueryInfo get expr string by calling func GetInfoExpr,
// then get response by calling func getExprResponse,
// return the success result
func (query *Query) QueryInfo() (*models.QueryRangeResult, error) {
	infoExpr := NewInfoExpr()
	infoExpr.GetInfoExpr(query)
	if query.Cluster == "" && query.App == "" && query.User == "" {
		query.Expr = infoExpr.Clusters
	}

	if query.Cluster != "" && query.User == "" && query.App == "" {
		query.Expr = infoExpr.Cluster
	}

	if query.Cluster != "" && query.User != "" && query.App == "" {
		query.Expr = infoExpr.User
	}

	if query.Cluster != "" && query.User != "" && query.App != "" {
		query.Expr = infoExpr.Application
	}

	response, request, err := query.getExprResponse()
	if err != nil {
		return nil, err
	}

	var result *models.QueryRangeResult
	err = json.Unmarshal(response, &result)
	if err != nil {
		err = fmt.Errorf("Failed to parse the response from %s", request)
		return nil, err
	}

	if result.Status != OKRESULT {
		err := fmt.Errorf("%s", result.Error)
		return nil, err
	}

	return result, nil
}

// QueryAppMetric get the expr string by calling getQueryMetricExpr
// then get response by calling getExprResponse
// then return the success result
func (query *Query) QueryAppMetric() (*models.QueryRangeResult, error) {
	start, end := timeRange(query.Start, query.End)
	query.Start = start
	query.End = end
	query.Expr = query.getQueryMetricExpr("app")

	response, request, err := query.getExprResponse()
	if err != nil {
		return nil, err
	}

	var result *models.QueryRangeResult
	err = json.Unmarshal(response, &result)
	if err != nil {
		err = fmt.Errorf("Failed to parse the response from %s", request)
		return nil, err
	}

	if result.Status != OKRESULT {
		err := fmt.Errorf("%s", result.Error)
		return nil, err
	}

	return result, nil
}

// QueryNodeMetric get the expr string by calling getQueryMetricExpr
// then get the response by calling getExprResponse
// then return the success result
func (query *Query) QueryNodeMetric() (*models.QueryRangeResult, error) {
	start, end := timeRange(query.Start, query.End)
	query.Start = start
	query.End = end
	query.Expr = query.getQueryMetricExpr("cluster")

	response, request, err := query.getExprResponse()
	if err != nil {
		return nil, err
	}

	var result *models.QueryRangeResult
	err = json.Unmarshal(response, &result)
	if err != nil {
		err = fmt.Errorf("Failed to parse the response from %s", request)
		return nil, err
	}

	if result.Status != OKRESULT {
		err := fmt.Errorf("%s", result.Error)
		return nil, err
	}

	return result, nil
}
