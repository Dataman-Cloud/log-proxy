package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/Dataman-Cloud/log-proxy/src/models"
	"github.com/Dataman-Cloud/log-proxy/src/utils/prometheusexpr"
	"github.com/prometheus/common/log"
)

const (
	// QUERYRANGEPATH define the query_range path of prometheus
	QUERYRANGEPATH = "/api/v1/query_range"
	// QUERYPATH define the query path of prometheus
	QUERYPATH = "/api/v1/query"
	//OKRESULT define the value of status field in response
	OKRESULT = "success"
)

func isInArray(array []string, value string) bool {
	for _, valueInList := range array {
		if value == valueInList {
			return true
		}
	}
	return false
}

// Query define the struct by query from prometheus
type Query struct {
	ExprTmpl   map[string]string
	HTTPClient *http.Client
	PromServer string
	Path       string
	*models.QueryParameter
}

// SetQueryExprsList get the expr strings
func SetQueryExprsList() map[string]string {
	var list = make(map[string]string, 0)
	for name, expr := range prometheusexpr.GetExprs() {
		list[name] = makeExprString(expr)
	}
	return list
}

func makeExprString(expr *prometheusexpr.Expr) string {
	var filter, byItems, queryExpr string
	filter = fmt.Sprintf("%s, %s, %s, %s, %s", expr.Filter.User, expr.Filter.Cluster, expr.Filter.App, expr.Filter.Slot, expr.Filter.Fixed)
	byItems = fmt.Sprintf("%s, %s, %s, %s", expr.By.Cluster, expr.By.User, expr.By.App, expr.By.Slot)

	queryExpr = fmt.Sprintf("%s{%s}", expr.Metric, filter)
	if expr.Function != "" {
		queryExpr = fmt.Sprintf("%s(%s[5m])", expr.Function, queryExpr)
	}
	if expr.Aggregation != "" {
		queryExpr = fmt.Sprintf("%s(%s) by (%s) keep_common", expr.Aggregation, queryExpr, byItems)
	} else {
		queryExpr = fmt.Sprintf("%s by (%s) keep_common", queryExpr, byItems)
	}
	return queryExpr
}

// GetQueryItemList set the Query exprs by utils.Expr
func (query *Query) GetQueryItemList() []string {
	list := make([]string, 0)
	for k := range prometheusexpr.GetExprs() {
		list = append(list, k)
	}
	sort.Strings(list)
	return list
}

//GetQueryUsers return the users list in the time interval
func (query *Query) GetQueryUsers() ([]string, error) {
	filter := prometheusexpr.GetFilter()
	by := prometheusexpr.GetBy()
	start, end := timeRange(query.Start, query.End)
	query.Start = start
	query.End = end
	query.Expr = fmt.Sprintf("count(container_tasks_state{%s, state='running'}) by (%s)", filter.Fixed, by.User)
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
	var list = make([]string, 0)

	for _, originData := range result.Data.Result {
		user := originData.Metric.ContainerLabelUser
		if !isInArray(list, user) {
			list = append(list, user)
		}
	}

	return list, nil
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

	log.Infof("prometheus qurey expr: %s\n", query.Expr)
	log.Infof("prometheus qurey result by url: %s\n", u.String())

	//log.Infof("prometheus qurey expr: %s", query.Expr)
	return body, u.String(), nil
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
