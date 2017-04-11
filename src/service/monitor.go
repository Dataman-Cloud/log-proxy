package service

import (
	"fmt"
	"net/http"
	"sort"

	"github.com/Dataman-Cloud/log-proxy/src/models"
	"github.com/Dataman-Cloud/log-proxy/src/utils/prometheusexpr"
)

const (
	cpuUsage      = "cpu"
	memPercentage = "memory"
	memUsage      = "memory_usage"
	memTotal      = "memory_total"
	networkRX     = "network_rx"
	networkTX     = "network_tx"
	fsRead        = "fs_read"
	fsWrite       = "fs_write"
	fsUsage       = "fs_usage"
	fsLimit       = "fs_limit"
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
	for n, v := range expr.By {
		if n != len(expr.By)-1 {
			byItems = byItems + v + ", "
		} else {
			byItems = byItems + v
		}
	}

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
