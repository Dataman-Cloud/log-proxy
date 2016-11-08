package service

import (
	"encoding/json"
	"errors"
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
	AppID      string
	Metric     string
	From       string
	To         string
	Step       string
}

const (
	TIMEFORMAT = "2006-01-02 15:04:05"
)

// GetMetricList will return the list of metric of query result.
// metric=cpu, will only query CPU usage
// metric=memory, will only query memory usage
// metric=all, will query each metic.
func (query *QueryRange) GetMetricList() (*models.MetricList, error) {
	list := new(models.MetricList)
	if query.Metric != "all" {
		data, err := query.QueryRangeFromProm()
		if err != nil {
			return nil, err
		}
		switch query.Metric {
		case "cpu":
			list.CPU = data.Data.Result
		case "memory":
			list.Memory = data.Data.Result
		default:
			return nil, errors.New("No this kind of metric.")
		}
		return list, nil
	}

	metrics := [...]string{"cpu", "memory"}
	for _, metric := range metrics {
		query.Metric = metric
		data, err := query.QueryRangeFromProm()
		if err != nil {
			return nil, err
		}
		switch metric {
		case "cpu":
			list.CPU = data.Data.Result
		case "memory":
			list.Memory = data.Data.Result
		default:
			return nil, errors.New("No this kind of metric.")
		}
	}

	return list, nil
}

func (query *QueryRange) QueryRangeFromProm() (*models.QueryRangeResult, error) {
	const (
		unixTime = "unix"
	)
	expr := query.setQueryExpr(query.Metric, query.AppID)

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
		q.Set("step", "15s")
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
	log.Printf("Get the prometheus qurey result by url: %s", u.String())

	return result, nil
}

// setQueryExpr will return the expr of prometheus query
func (query *QueryRange) setQueryExpr(metrics, appID string) (expr string) {
	switch metrics {
	case "cpu":
		expr = "avg(irate(container_cpu_usage_seconds_total{container_label_APP_ID='" + appID + "',id=~'/docker/.*', name=~'mesos.*'}[5m])) by (container_label_APP_ID, group, id, image, instance, job, name)"
	case "memory":
		expr = "container_memory_usage_bytes{container_label_APP_ID='" + appID + "',id=~'/docker/.*', name=~'mesos.*'} / container_spec_memory_limit_bytes{container_label_APP_ID='nginx-stress', id=~'/docker/.*', name=~'mesos.*'}"
	default:
		expr = ""
	}
	return expr
}

func timeRange(f, t, unixTime string) (string, string, error) {
	if f == "" && t == "" {
		to := time.Now()
		from := timeOffset(to, "-30m")
		return timeConvertString(from, unixTime), timeConvertString(to, unixTime), nil
	}

	if f == "" && t != "" {
		to, err := time.Parse(TIMEFORMAT, t)
		if err != nil {
			return "", "", err
		}
		from := timeOffset(to, "-30m")
		return timeConvertString(from, unixTime), timeConvertString(to, unixTime), nil
	}

	if f != "" && t == "" {
		from, err := time.Parse(TIMEFORMAT, f)
		if err != nil {
			return "", "", err
		}
		to := timeOffset(from, "30m")
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
