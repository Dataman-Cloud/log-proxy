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
	TIMEFORMAT     = "2006-01-02 15:04:05"
	QUERYRANGEPATH = "/api/v1/query_range"
	QUERYPATH      = "/api/v1/query"
	RULESPATH      = "/api/v1/rules"
)

type Query struct {
	HttpClient *http.Client
	PromServer string
	Path       string
	*QueryParameter
}

type QueryParameter struct {
	Metric    string
	ClusterID string
	AppID     string
	TaskID    string
	NodeID    string
	Start     string
	End       string
	Step      string
	Period    string
	Expr      string
}

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
	resp, err := query.HttpClient.Get(u.String())
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

func (query *Query) QueryMetric() (*models.QueryRangeResult, error) {
	// if start/end as "", set them as now
	start, end := timeRange(query.Start, query.End)
	query.Start = start
	query.End = end
	query.Expr = query.getQueryMetricExpr("id")

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

	if result.Status != "success" {
		err := fmt.Errorf("%s", result.Error)
		return nil, err
	}

	return result, nil
}

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

func (query *Query) QueryInfo() (*models.QueryRangeResult, error) {
	infoExpr := NewInfoExpr()
	infoExpr.GetInfoExpr(query)
	if query.ClusterID == "" && query.AppID == "" {
		query.Expr = infoExpr.Clusters
	} else if query.ClusterID != "" {
		query.Expr = infoExpr.Cluster
	} else if query.AppID != "" {
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

	if result.Status != "success" {
		err := fmt.Errorf("%s", result.Error)
		return nil, err
	}

	return result, nil
}

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

	if result.Status != "success" {
		err := fmt.Errorf("%s", result.Error)
		return nil, err
	}

	return result, nil
}

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

	if result.Status != "success" {
		err := fmt.Errorf("%s", result.Error)
		return nil, err
	}

	return result, nil
}
