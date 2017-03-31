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

	"github.com/Dataman-Cloud/log-proxy/src/backends"
	"github.com/Dataman-Cloud/log-proxy/src/models"
	"github.com/Dataman-Cloud/log-proxy/src/utils/prometheusexpr"

	log "github.com/Sirupsen/logrus"
)

const (
	// QUERYRANGEPATH define the query_range path of prometheus
	QUERYRANGEPATH = "/api/v1/query_range"
	// QUERYPATH define the query path of prometheus
	QUERYPATH = "/api/v1/query"
	//OKRESULT define the value of status field in response
	OKRESULT = "success"

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

// Metric defines the JSON fomat from the Query Result
type Metric struct {
	CPU        *models.MetricCPU        `json:"cpu"`
	Memory     *models.MetricMemory     `json:"memory"`
	Network    *models.MetricNewtork    `json:"network"`
	Filesystem *models.MetricFilesystem `json:"filesystem"`
}

// NewMetric init Metric
func NewMetric() *Metric {
	return &Metric{
		CPU:        &models.MetricCPU{},
		Memory:     &models.MetricMemory{},
		Network:    &models.MetricNewtork{},
		Filesystem: &models.MetricFilesystem{},
	}
}

// GetQueryMetric gets the result by calling func QueryMetric
// and then set the vaule of the fields in struct Metric
func (m *Metric) GetQueryMetric(query *backends.Query) error {
	data, err := query.QueryMetric()
	if err != nil {
		return err
	}
	switch query.Metric {
	case cpuUsage:
		m.CPU.Usage = data.Data.Result
	case memPercentage:
		m.Memory.Percentage = data.Data.Result
	case memUsage:
		m.Memory.Usage = data.Data.Result
	case memTotal:
		m.Memory.Total = data.Data.Result
	case networkRX:
		m.Network.Receive = data.Data.Result
	case networkTX:
		m.Network.Transmit = data.Data.Result
	case fsRead:
		m.Filesystem.Read = data.Data.Result
	case fsWrite:
		m.Filesystem.Write = data.Data.Result
	default:
		return errors.New("No this kind of metric")
	}
	return nil
}

// Info defines the JSON format of information in clusters/cluster/apps/nodes
type Info struct {
	Clusters     map[string]*ClusterInfo `json:"clusters"`
	Users        []string                `json:"users"`
	Applications []string                `json:"applications"`
	Nodes        []string                `json:"nodes"`
}

// NewInfo init struct Info
func NewInfo() *Info {
	return &Info{
		Clusters:     make(map[string]*ClusterInfo),
		Users:        []string{},
		Applications: []string{},
		Nodes:        []string{},
	}
}

// ClusterInfo defines the JSON format of information in Cluster
type ClusterInfo struct {
	Users map[string]*UserInfo `json:"users"`
}

// NewClusterInfo init struct ClusterInfo
func NewClusterInfo() *ClusterInfo {
	return &ClusterInfo{
		Users: make(map[string]*UserInfo),
	}
}

// UserInfo defines the JSON format of information in User
type UserInfo struct {
	Applications map[string]*AppInfo `json:"applications"`
	Tasks        []string            `json:"tasks"`
	Nodes        []string            `json:"nodes"`
}

// NewUserInfo init the struct UserInfo
func NewUserInfo() *UserInfo {
	return &UserInfo{
		Applications: make(map[string]*AppInfo),
		Tasks:        []string{},
		Nodes:        []string{},
	}
}

// AppInfo defines the JSON format of information in application
type AppInfo struct {
	CPU        *models.InfoCPU        `json:"cpu,omitempty"`
	Memory     *models.InfoMemory     `json:"memory,omitempty"`
	Network    *models.InfoNetwork    `json:"network,omitempty"`
	Filesystem *models.InfoFilesystem `json:"filesystem,omitempty"`
	Tasks      map[string]*TaskInfo   `json:"tasks"`
}

// NewAppInfo init struct AppInfo
func NewAppInfo() *AppInfo {
	return &AppInfo{
		CPU:        &models.InfoCPU{},
		Memory:     &models.InfoMemory{},
		Network:    &models.InfoNetwork{},
		Filesystem: &models.InfoFilesystem{},
		Tasks:      make(map[string]*TaskInfo),
	}
}

// TaskInfo defines the JSON format of information in task(container)
type TaskInfo struct {
	Slot       string                         `json:"slot"`
	Node       string                         `json:"node"`
	CPU        *models.InfoCPU                `json:"cpu"`
	Memory     *models.InfoMemory             `json:"memory"`
	Network    map[string]*models.InfoNetwork `json:"network"`
	Filesystem *models.InfoFilesystem         `json:"filesystem"`
}

// NewTaskInfo init the struct TaskInfo
func NewTaskInfo() *TaskInfo {
	return &TaskInfo{
		CPU:        &models.InfoCPU{},
		Memory:     &models.InfoMemory{},
		Network:    make(map[string]*models.InfoNetwork),
		Filesystem: &models.InfoFilesystem{},
	}
}

// GetQueryInfo get the info of Cluster, User, App, Task from Query Result
func (info *Info) GetQueryInfo(query *backends.Query) error {
	data, err := query.QueryInfo()
	if err != nil {
		return err
	}
	info.GetClustersInfo(query, data)
	if query.Cluster != "" && query.User != "" {
		info.GetAppInfo(query, data)
	}

	if query.Cluster != "" && query.User != "" && query.App != "" {
		err = info.GetAppInfoMetric(query)
		if err != nil {
			return err
		}
		err = info.GetTaskInfoMetric(query)
		if err != nil {
			return err
		}
	}

	return nil
}

// GetClustersInfo get the info of Clusters from Query Result
func (info *Info) GetClustersInfo(query *backends.Query, data *models.QueryRangeResult) {
	// Set the cluster into list
	for _, originData := range data.Data.Result {
		cluster := originData.Metric.ContainerLabelCluster
		info.Clusters[cluster] = NewClusterInfo()
	}

	for _, originData := range data.Data.Result {
		cluster := originData.Metric.ContainerLabelCluster
		user := originData.Metric.ContainerLabelUser

		for name, value := range info.Clusters {
			if cluster == name {
				value.Users[user] = NewUserInfo()
				if !isInArray(info.Users, user) {
					info.Users = append(info.Users, user)
				}
			}
		}
	}
	for _, originData := range data.Data.Result {
		cluster := originData.Metric.ContainerLabelCluster
		user := originData.Metric.ContainerLabelUser
		app := originData.Metric.ContainerLabelApp
		task := originData.Metric.ID
		node := strings.Split(originData.Metric.Instance, ":")[0]
		for clusterName, ClusterValue := range info.Clusters {
			if cluster == clusterName {
				for name, value := range ClusterValue.Users {
					if user == name {
						value.Applications[app] = NewAppInfo()
						if !isInArray(value.Tasks, task) {
							value.Tasks = append(value.Tasks, task)
						}
						if !isInArray(value.Nodes, node) {
							value.Nodes = append(value.Nodes, node)
						}
						appUID := fmt.Sprintf("%s.%s.%s", cluster, user, app)
						if !isInArray(info.Applications, appUID) {
							info.Applications = append(info.Applications, appUID)
						}
						if !isInArray(info.Nodes, node) {
							info.Nodes = append(info.Nodes, node)
						}
					}
				}
			}
		}
	}
}

// GetAppInfo get the info of application from Query Result
func (info *Info) GetAppInfo(query *backends.Query, data *models.QueryRangeResult) {
	// Fill the info of container in application
	for _, originData := range data.Data.Result {
		cluster := originData.Metric.ContainerLabelCluster
		user := originData.Metric.ContainerLabelUser
		app := originData.Metric.ContainerLabelApp
		task := originData.Metric.ID
		slot := originData.Metric.ContainerLabelSlot
		node := strings.Split(originData.Metric.Instance, ":")[0]
		for clusterName, ClusterValue := range info.Clusters {
			if cluster == clusterName {
				for userName, userValue := range ClusterValue.Users {
					if user == userName {
						for name, value := range userValue.Applications {
							if app == name {
								value.Tasks[task] = NewTaskInfo()
								value.Tasks[task].Slot = slot
								value.Tasks[task].Node = node
							}
						}
					}
				}
			}
		}
	}
}

// GetAppInfoMetric get the metric data from Query Result
func (info *Info) GetAppInfoMetric(query *backends.Query) error {
	metrics := []string{cpuUsage, memPercentage, networkRX, networkTX, fsRead, fsWrite}
	for _, metric := range metrics {
		query.Metric = metric
		query.Path = backends.QUERYRANGEPATH
		data, err := query.QueryAppMetric()
		if err != nil {
			return err
		}
		for _, originData := range data.Data.Result {
			cluster := originData.Metric.ContainerLabelCluster
			user := originData.Metric.ContainerLabelUser
			app := originData.Metric.ContainerLabelApp
			metricValue := originData.Values[0]
			for clusterName, ClusterValue := range info.Clusters {
				if cluster == clusterName {
					for userName, userValue := range ClusterValue.Users {
						if user == userName {
							for name, value := range userValue.Applications {
								if app == name {
									switch query.Metric {
									case cpuUsage:
										value.CPU.Usage = metricValue
									case memPercentage:
										value.Memory.Percetange = metricValue
									case networkRX:
										value.Network.Receive = metricValue
									case networkTX:
										value.Network.Transmit = metricValue
									case fsRead:
										value.Filesystem.Read = metricValue
									case fsWrite:
										value.Filesystem.Write = metricValue
									}
								}
							}
						}
					}
				}
			}
		}
	}
	return nil
}

// GetTaskInfoMetric get the metric data from Query Result
func (info *Info) GetTaskInfoMetric(query *backends.Query) error {
	metrics := []string{cpuUsage, memUsage, memTotal, networkRX, networkTX, fsRead, fsWrite}
	for _, metric := range metrics {
		query.Metric = metric
		data, err := query.QueryMetric()
		if err != nil {
			return nil
		}

		if metric == networkRX || metric == networkTX {
			for _, clusterValue := range info.Clusters {
				for _, userValue := range clusterValue.Users {
					for _, appValue := range userValue.Applications {
						for _, value := range appValue.Tasks {
							if len(value.Network) == 0 {
								for _, originData := range data.Data.Result {
									nic := originData.Metric.Interface
									value.Network[nic] = models.NewInfoNetwork()
								}
							}
						}
					}
				}
			}
		}

		for _, originData := range data.Data.Result {
			task := originData.Metric.ID
			metricValue := originData.Values[0]
			nic := originData.Metric.Interface

			for _, clusterValue := range info.Clusters {
				for _, userValue := range clusterValue.Users {
					for _, appValue := range userValue.Applications {
						for name, value := range appValue.Tasks {
							if task == name {
								switch query.Metric {
								case cpuUsage:
									value.CPU.Usage = metricValue
								case memUsage:
									value.Memory.Usage = metricValue
								case memTotal:
									value.Memory.Total = metricValue
								case networkRX:
									for nicK, nicV := range value.Network {
										if nic == nicK {
											nicV.Receive = metricValue
										}
									}
								case networkTX:
									for nicK, nicV := range value.Network {
										if nic == nicK {
											nicV.Transmit = metricValue
										}
									}
								case fsRead:
									value.Filesystem.Read = metricValue
								case fsWrite:
									value.Filesystem.Write = metricValue
								}
							}
						}
					}
				}
			}
		}
	}
	return nil
}

func isInArray(array []string, value string) bool {
	for _, valueInList := range array {
		if value == valueInList {
			return true
		}
	}
	return false
}

// NodesInfo defines the JSON format of Nodes list
type NodesInfo struct {
	Nodes map[string]*NodeInfo `json:"nodes"`
}

// NewNodesInfo init NodesInfo
func NewNodesInfo() *NodesInfo {
	return &NodesInfo{
		Nodes: make(map[string]*NodeInfo),
	}
}

// NodeInfo defines the JSON format of information in Node
type NodeInfo struct {
	CPU        *models.InfoCPU                   `json:"cpu"`
	Memory     *models.InfoMemory                `json:"memory"`
	Network    map[string]*models.InfoNetwork    `json:"network"`
	Filesystem map[string]*models.InfoFilesystem `json:"filesystem"`
}

// NewNodeInfo init the struct NodeInfo
func NewNodeInfo() *NodeInfo {
	return &NodeInfo{
		CPU:        &models.InfoCPU{},
		Memory:     &models.InfoMemory{},
		Network:    make(map[string]*models.InfoNetwork),
		Filesystem: make(map[string]*models.InfoFilesystem),
	}
}

// GetQueryNodesInfo gets the result by calling QueryNodeMetric
// then set the value of fields in NodesInfo
func (ni *NodesInfo) GetQueryNodesInfo(query *backends.Query) error {
	metrics := []string{cpuUsage, memUsage, memTotal, networkRX, networkTX, fsUsage, fsLimit}
	for _, metric := range metrics {
		query.Metric = metric
		data, err := query.QueryNodeMetric()
		if err != nil {
			return err
		}

		if len(ni.Nodes) == 0 {
			for _, originData := range data.Data.Result {
				name := strings.Split(originData.Metric.Instance, ":")[0]
				ni.Nodes[name] = NewNodeInfo()
			}
		}
		if metric == networkRX || metric == networkTX {
			for _, node := range ni.Nodes {
				if len(node.Network) == 0 {
					for _, originData := range data.Data.Result {
						nic := originData.Metric.Interface
						node.Network[nic] = models.NewInfoNetwork()
					}
				}
			}
		}

		if metric == fsUsage || metric == fsLimit {
			for _, node := range ni.Nodes {
				if len(node.Filesystem) == 0 {
					for _, originData := range data.Data.Result {
						device := originData.Metric.Device
						node.Filesystem[device] = models.NewInfoFilesystem()
					}
				}
			}
		}

		for _, originData := range data.Data.Result {
			name := strings.Split(originData.Metric.Instance, ":")[0]
			nic := originData.Metric.Interface
			device := originData.Metric.Device
			value := originData.Values[0]
			for k, v := range ni.Nodes {
				if name == k {
					switch query.Metric {
					case cpuUsage:
						v.CPU.Usage = value
					case memUsage:
						v.Memory.Usage = value
					case memTotal:
						v.Memory.Total = value
					case networkRX:
						for nicK, nicV := range v.Network {
							if nic == nicK {
								nicV.Receive = value
							}
						}
					case networkTX:
						for nicK, nicV := range v.Network {
							if nic == nicK {
								nicV.Transmit = value
							}
						}
					case fsUsage:
						for fsK, fsV := range v.Filesystem {
							if device == fsK {
								fsV.Usage = value
							}
						}
					case fsLimit:
						for fsK, fsV := range v.Filesystem {
							if device == fsK {
								fsV.Limit = value
							}
						}
					}
				}
			}
		}
	}
	return nil
}

type Query struct {
	ExprTmpl   map[string]string
	HTTPClient *http.Client
	PromServer string
	Path       string
	*models.QueryParameter
}

// GetQueryItemList set the Query exprs by utils.Expr
func (query *Query) GetQueryItemList() []string {
	list := make([]string, 0)
	for k := range prometheusexpr.GetExprs() {
		list = append(list, k)
	}
	return list
}

// GetQueryClusters set the Query exprs by utils.Expr
func (query *Query) GetQueryClusters() ([]string, error) {
	start, end := timeRange(query.Start, query.End)
	query.Start = start
	query.End = end
	query.Expr = "count(container_tasks_state{id=~'/docker/.*', name=~'mesos.*', state='running', container_label_DM_LOG_TAG!='ignore'}) by (container_label_VCLUSTER)"

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
		cluster := originData.Metric.ContainerLabelVcluster
		if !isInArray(list, cluster) {
			list = append(list, cluster)
		}
	}

	return list, nil
}

// GetQueryApps set the Query exprs by utils.Expr
func (query *Query) GetQueryApps() ([]string, error) {
	start, end := timeRange(query.Start, query.End)
	query.Start = start
	query.End = end
	query.Expr = fmt.Sprintf("count(container_tasks_state{container_label_VCLUSTER='%s', container_label_DM_LOG_TAG!='ignore'}) by (container_label_APP_ID)", query.Cluster)

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
		app := originData.Metric.ContainerLabelAppID
		if !isInArray(list, app) {
			list = append(list, app)
		}
	}

	return list, nil
}

// GetQueryTasks set the Query exprs by utils.Expr
func (query *Query) GetQueryAppTasks() ([]string, error) {
	start, end := timeRange(query.Start, query.End)
	query.Start = start
	query.End = end
	query.Expr = fmt.Sprintf("count(container_tasks_state{container_label_VCLUSTER='%s', container_label_APP_ID='%s', id=~'/docker/.*', name=~'mesos.*', state='running', container_label_DM_LOG_TAG!='ignore'}) by (container_env_mesos_task_id)", query.Cluster, query.App)

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
		task := originData.Metric.ContainerLabelTask
		if !isInArray(list, task) {
			list = append(list, task)
		}
	}

	return list, nil
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
	filter = fmt.Sprintf("%s, %s, %s, %s", expr.Filter.Cluster, expr.Filter.App, expr.Filter.Task, expr.Filter.Fixed)
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

// getQueryMetricExpr return the expr string
func (query *Query) getQueryMetricExpr() string {
	tmpl := query.ExprTmpl[query.Metric]
	cluster := query.Cluster
	app := query.App
	task := query.Task
	exprTempl := fmt.Sprintf(tmpl, cluster, app, task)
	return exprTempl
}

// QueryExpr return the result by calling function getExprResponse
func (query *Query) QueryExpr() (*models.QueryExprResult, error) {
	start, end := timeRange(query.Start, query.End)
	query.Start = start
	query.End = end

	response, request, err := query.getExprResponse()
	if err != nil {
		return nil, err
	}

	var result *models.QueryExprResult
	err = json.Unmarshal(response, &result)
	if err != nil {
		err = fmt.Errorf("Failed to parse the response from %s", request)
		return nil, err
	}
	result.Expr = query.Expr
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

	log.Infof("prometheus qurey result by url: %s", u.String())
	//log.Infof("prometheus qurey expr: %s", query.Expr)
	return body, u.String(), nil
}

func (query *Query) QueryMetric() (*models.QueryRangeResult, error) {
	// if start/end as "", set them as now
	start, end := timeRange(query.Start, query.End)
	query.Start = start
	query.End = end

	query.Expr = query.getQueryMetricExpr()
	response, request, err := query.getExprResponse()
	if err != nil {
		log.Infoln("QueryMetric: err", err)
		return nil, err
	}
	log.Infoln("QueryMetric: request", request)

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

	result.Expr = query.Expr

	log.Infoln("QueryMetric: result", result)
	return result, nil
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
