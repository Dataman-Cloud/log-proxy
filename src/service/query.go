package service

import (
	"encoding/json"
	"errors"
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

type QueryRange struct {
	HttpClient *http.Client
	PromServer string
	Path       string
	Type       string
	AppID      string
	TaskID     string
	Node       string
	Metric     string
	From       string
	To         string
	Step       string
}

const (
	TIMEFORMAT = "2006-01-02 15:04:05"
)

func (query *QueryRange) QueryRangeFromProm() (*models.QueryRangeResult, error) {
	const (
		unixTime = "unix"
	)
	var expr string
	if query.Type == "app" {
		expr = query.setQueryAppExpr(query.Metric, query.AppID)
	} else if query.Type == "node" {
		expr = query.setQueryNodesExpr(query.Metric, query.Node)
	} else {
		expr = query.setQueryExpr(query.Metric, query.AppID, query.TaskID)
	}

	u, err := url.Parse(query.PromServer)
	if err != nil {
		return nil, err
	}

	u.Path = strings.TrimRight(u.Path, "/") + query.Path
	q := u.Query()
	q.Set("query", expr)

	var start, end string
	start, end, err = timeRange(query.From, query.To, unixTime)
	if err != nil {
		start = query.From
		end = query.To
	}
	q.Set("start", start)
	q.Set("end", end)

	if query.Step == "" {
		q.Set("step", "30s")
	} else {
		q.Set("step", query.Step)
	}

	u.RawQuery = q.Encode()

	resp, err := query.HttpClient.Get(u.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result *models.QueryRangeResult
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	if result.Status != "success" {
		err := fmt.Errorf("%s", result.Error)
		return nil, err
	}

	log.Printf("Get the prometheus qurey result by url: %s", u.String())

	return result, nil
}

// setQueryExpr will return the expr of prometheus query
func (query *QueryRange) setQueryExpr(metrics, appID, taskID string) (expr string) {
	switch metrics {
	case "cpu":
		expr = "avg(irate(container_cpu_usage_seconds_total{container_label_APP_ID='" + appID +
			"',id=~'/docker/" + taskID + ".*', name=~'mesos.*'}[5m])) by (container_label_APP_ID, group, " +
			"id, image, instance, job, name)"
	case "memory":
		expr = "container_memory_usage_bytes{container_label_APP_ID='" + appID +
			"',id=~'/docker/" + taskID + ".*', name=~'mesos.*'} / container_spec_memory_limit_bytes" +
			"{container_label_APP_ID='" + appID + "', id=~'/docker/" + taskID + ".*', name=~'mesos.*'}"
	case "memory_usage":
		expr = "container_memory_usage_bytes{container_label_APP_ID='" + appID +
			"',id=~'/docker/" + taskID + ".*', name=~'mesos.*'}"
	case "memory_total":
		expr = "container_spec_memory_limit_bytes{container_label_APP_ID='" + appID +
			"',id=~'/docker/" + taskID + ".*', name=~'mesos.*'}"
	case "network_rx":
		expr = "container_network_receive_bytes_total{container_label_APP_ID='" + appID +
			"',id=~'/docker/" + taskID + ".*', name=~'mesos.*'}"
	case "network_tx":
		expr = "container_network_transmit_bytes_total{container_label_APP_ID='" + appID +
			"',id=~'/docker/" + taskID + ".*', name=~'mesos.*'}"
	case "fs_read":
		expr = "container_fs_reads_total{container_label_APP_ID='" + appID +
			"',id=~'/docker/" + taskID + ".*', name=~'mesos.*'}"
	case "fs_write":
		expr = "container_fs_writes_total{container_label_APP_ID='" + appID +
			"',id=~'/docker/" + taskID + ".*', name=~'mesos.*'}"
	default:
		expr = ""
	}
	return expr
}

// setQueryExpr will return the expr of prometheus query
func (query *QueryRange) setQueryAppExpr(metrics, appID string) (expr string) {
	switch metrics {
	case "cpu":
		expr = "avg(irate(container_cpu_usage_seconds_total{container_label_APP_ID='" + appID +
			"',id=~'/docker/.*', name=~'mesos.*'}[5m])) by (container_label_APP_ID)"
	case "memory":
		expr = "avg(container_memory_usage_bytes{container_label_APP_ID='" + appID +
			"',id=~'/docker/.*'} / container_spec_memory_limit_bytes{container_label_APP_ID='" + appID +
			"',id=~'/docker/.*'}) by (container_label_APP_ID)"
	case "network_rx":
		expr = "sum(container_network_receive_bytes_total{container_label_APP_ID='" + appID +
			"',id=~'/docker/.*'}) by (container_label_APP_ID)"
	case "network_tx":
		expr = "sum(container_network_transmit_bytes_total{container_label_APP_ID='" + appID +
			"',id=~'/docker/.*'}) by (container_label_APP_ID)"
	case "fs_read":
		expr = "sum(container_fs_reads_total{container_label_APP_ID='" + appID +
			"',id=~'/docker/.*'}) by (container_label_APP_ID)"
	case "fs_write":
		expr = "sum(container_fs_writes_total{container_label_APP_ID='" + appID +
			"',id=~'/docker/.*'}) by (container_label_APP_ID)"
	default:
		expr = ""
	}
	return expr
}

func timeRange(f, t, unixTime string) (string, string, error) {
	const (
		OFFSET = "0m"
	)

	if f == "" && t == "" {
		to := time.Now()
		from := timeOffset(to, OFFSET)
		return timeConvertString(from, unixTime), timeConvertString(to, unixTime), nil
	}

	if f == "" && t != "" {
		to, err := time.Parse(TIMEFORMAT, t)
		if err != nil {
			return "", "", err
		}
		from := timeOffset(to, OFFSET)
		return timeConvertString(from, unixTime), timeConvertString(to, unixTime), nil
	}

	if f != "" && t == "" {
		from, err := time.Parse(TIMEFORMAT, f)
		if err != nil {
			return "", "", err
		}
		to := timeOffset(from, OFFSET)
		return timeConvertString(from, unixTime), timeConvertString(to, unixTime), nil
	}

	from, err := time.Parse(TIMEFORMAT, f)
	if err != nil {
		return "", "", err
	}
	to, err := time.Parse(TIMEFORMAT, t)
	if err != nil {
		return "", "", err
	}

	if to.Before(from) {
		return "", "", errors.New("The from time is after the to time.")
	}

	return timeConvertString(from, unixTime), timeConvertString(to, unixTime), nil
}

func timeConvertString(t time.Time, unixTime string) string {
	var toUnix int64
	switch unixTime {
	case "unix":
		toUnix = t.Unix()
	case "unixnano":
		toUnix = t.UnixNano() / 1000000
	default:
		toUnix = t.Unix()
	}

	return strconv.FormatInt(toUnix, 10)
}

func timeOffset(t time.Time, offset string) time.Time {
	duration, _ := time.ParseDuration(offset)
	return t.Add(duration)
}

func (query *QueryRange) QueryAppsFromProm() (*models.QueryRangeResult, error) {
	var expr string
	if query.AppID == "" {
		expr = "container_tasks_state{container_label_APP_ID=~'.*', id=~'/docker/.*', name=~'mesos.*', state='running'}"
	} else {
		expr = "container_tasks_state{container_label_APP_ID='" + query.AppID + "', id=~'/docker/.*', name=~'mesos.*', state='running'}"
	}

	u, err := url.Parse(query.PromServer)
	if err != nil {
		return nil, err
	}

	u.Path = strings.TrimRight(u.Path, "/") + query.Path
	q := u.Query()
	q.Set("query", expr)

	u.RawQuery = q.Encode()

	resp, err := query.HttpClient.Get(u.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result *models.QueryRangeResult
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	if result.Status != "success" {
		err := fmt.Errorf("%s", result.Error)
		return nil, err
	}

	log.Printf("Get the prometheus qurey result by url: %s", u.String())

	return result, nil
}

// setQueryExpr will return the expr of prometheus query
func (query *QueryRange) setQueryNodesExpr(metric, node string) (expr string) {
	if node != "" {
		switch metric {
		case "cpu":
			expr = "avg(irate(container_cpu_usage_seconds_total{id='/',instance='" + node + ":5014'}[5m])) by (instance)"
		case "memory":
			expr = "sum(container_memory_usage_bytes{id='/',instance='" + node +
				":5014'} / container_spec_memory_limit_bytes{id='/',instance='" + node + ":5014'}) by (instance)"
		case "memory_usage":
			expr = "sum(container_memory_usage_bytes{id='/',instance='" + node + ":5014'}) by (instance)"
		case "memory_total":
			expr = "sum(container_spec_memory_limit_bytes{id='/',instance='" + node + ":5014'}) by (instance)"
		case "network_rx":
			expr = "sum(container_network_receive_bytes_total{id=~'/',instance='" + node + ":5014'}) by (instance)"
		case "network_tx":
			expr = "sum(container_network_transmit_bytes_total{id=~'/',instance='" + node + ":5014'}) by (instance)"
		default:
			expr = ""
		}
		return expr
	}

	switch metric {
	case "cpu":
		expr = "avg(irate(container_cpu_usage_seconds_total{id='/'}[5m])) by (instance)"
	case "memory":
		expr = "sum(container_memory_usage_bytes{id='/'} / container_spec_memory_limit_bytes{id='/'}) by (instance)"
	case "memory_usage":
		expr = "sum(container_memory_usage_bytes{id='/'}) by (instance)"
	case "memory_total":
		expr = "sum(container_spec_memory_limit_bytes{id='/'}) by (instance)"
	case "network_rx":
		expr = "sum(container_network_receive_bytes_total{id=~'/'}) by (instance)"
	case "network_tx":
		expr = "sum(container_network_transmit_bytes_total{id=~'/'}) by (instance)"
	default:
		expr = ""
	}
	return expr
}

func (query *QueryRange) QueryNodesFromProm() (*models.QueryRangeResult, error) {
	const (
		unixTime = "unix"
	)
	expr := query.setQueryNodesExpr(query.Metric, query.Node)

	u, err := url.Parse(query.PromServer)
	if err != nil {
		return nil, err
	}

	u.Path = strings.TrimRight(u.Path, "/") + query.Path
	q := u.Query()
	q.Set("query", expr)

	start, end, err := timeRange(query.From, query.To, unixTime)
	if err != nil {
		return nil, err
	}
	q.Set("start", start)
	q.Set("end", end)

	if query.Step == "" {
		q.Set("step", "30s")
	} else {
		q.Set("step", query.Step)
	}

	u.RawQuery = q.Encode()

	resp, err := query.HttpClient.Get(u.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result *models.QueryRangeResult
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	if result.Status != "success" {
		err := fmt.Errorf("%s", result.Error)
		return nil, err
	}
	log.Printf("Get the prometheus qurey result by url: %s", u.String())

	return result, nil
}
