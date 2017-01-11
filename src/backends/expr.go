package backends

import "fmt"

// MetricExpr define the expr strings of metric type
type MetricExpr struct {
	CPU        *MetricExprCPU
	Memory     *MetricExprMemory
	Network    *MetricExprNewtork
	Filesystem *MetricExprFilesystem
}

// NewMetricExpr init the struct MetricExpr
func NewMetricExpr() *MetricExpr {
	return &MetricExpr{
		CPU:        &MetricExprCPU{},
		Memory:     &MetricExprMemory{},
		Network:    &MetricExprNewtork{},
		Filesystem: &MetricExprFilesystem{},
	}
}

// MetricExprCPU define the expr string of CPU usage
type MetricExprCPU struct {
	Usage string
}

// MetricExprMemory define the expr strings of CPU Percentage/Usage/Total
type MetricExprMemory struct {
	Percentage string
	Usage      string
	Total      string
}

// MetricExprNewtork define the expr strings of network receive/transmit
type MetricExprNewtork struct {
	Receive  string
	Transmit string
}

// MetricExprFilesystem define the expr strings of Filesystem Read/Write
// Usage/Limit
type MetricExprFilesystem struct {
	Read  string
	Write string
	Usage string
	Limit string
}

// GetExpr return the struct MetricExpr by Query and Level string
func GetExpr(query *Query, level string) *MetricExpr {
	me := NewMetricExpr()

	switch level {
	case "task":
		me.GetMetricExpr(query)
	case "app":
		me.GetAppExpr(query)
	case "cluster":
		me.GetNodeExpr(query)
	}
	return me
}

// GetMetricExpr return expr by the Tasks labels
func (expr *MetricExpr) GetMetricExpr(query *Query) {
	byItems := "container_label_CLUSTER_ID,container_label_APP_ID," +
		"container_label_SLOT_ID,container_label_TASK_ID,container_label_USER_ID," +
		"group,id,image,instance,job,name,interface,device"
	expr.setMetricExpr(query, byItems)
}

func (expr *MetricExpr) setMetricExpr(query *Query, byItems string) {
	clusterid := query.ClusterID
	appid := query.AppID
	userid := query.UserID
	slotid := query.SlotID

	expr.CPU.Usage = fmt.Sprintf("avg(irate(container_cpu_usage_seconds_total{container_label_CLUSTER_ID='%s',"+
		"container_label_APP_ID='%s',container_label_SLOT_ID=~'%s',container_label_USER_ID='%s',"+
		"id=~'/docker/.*',name=~'mesos.*'}[5m])) by (%s)", clusterid, appid, slotid, userid, byItems)

	expr.Memory.Usage = fmt.Sprintf("sum(container_memory_usage_bytes{container_label_CLUSTER_ID='%s',"+
		"container_label_APP_ID='%s',container_label_SLOT_ID=~'%s',container_label_USER_ID='%s',"+
		"id=~'/docker/.*',name=~'mesos.*'}) by (%s)", clusterid, appid, slotid, userid, byItems)

	expr.Memory.Total = fmt.Sprintf("sum(container_spec_memory_limit_bytes{container_label_CLUSTER_ID='%s',"+
		"container_label_APP_ID='%s',container_label_SLOT_ID=~'%s',container_label_USER_ID='%s',"+
		"id=~'/docker/.*',name=~'mesos.*'}) by (%s)", clusterid, appid, slotid, userid, byItems)

	expr.Memory.Percentage = fmt.Sprintf("%s / %s", expr.Memory.Usage, expr.Memory.Total)

	expr.Network.Receive = fmt.Sprintf("sum(irate(container_network_receive_bytes_total{container_label_CLUSTER_ID='%s',"+
		"container_label_APP_ID='%s',container_label_SLOT_ID=~'%s',container_label_USER_ID='%s',"+
		"id=~'/docker/.*',name=~'mesos.*'}[5m])) by (%s)", clusterid, appid, slotid, userid, byItems)

	expr.Network.Transmit = fmt.Sprintf("sum(irate(container_network_transmit_bytes_total{container_label_CLUSTER_ID='%s',"+
		"container_label_APP_ID='%s', container_label_SLOT_ID=~'%s',container_label_USER_ID='%s',"+
		"id=~'/docker/.*',name=~'mesos.*'}[5m])) by (%s)", clusterid, appid, slotid, userid, byItems)

	expr.Filesystem.Read = fmt.Sprintf("sum(irate(container_fs_reads_total{container_label_CLUSTER_ID='%s',"+
		"container_label_APP_ID='%s',container_label_SLOT_ID=~'%s',container_label_USER_ID='%s',"+
		"id=~'/docker/.*',name=~'mesos.*'}[5m])) by (%s)", clusterid, appid, slotid, userid, byItems)

	expr.Filesystem.Write = fmt.Sprintf("sum(irate(container_fs_writes_total{container_label_CLUSTER_ID='%s',"+
		"container_label_APP_ID='%s',container_label_SLOT_ID=~'%s',container_label_USER_ID='%s',"+
		"id=~'/docker/.*',name=~'mesos.*'}[5m])) by (%s)", clusterid, appid, slotid, userid, byItems)
}

// GetAppExpr return the expr by the App labels
func (expr *MetricExpr) GetAppExpr(query *Query) {
	byItems := "container_label_CLUSTER_ID,container_label_APP_ID,container_label_USER_ID"
	expr.setMetricExpr(query, byItems)
}

// GetNodeExpr return the struct MetricExpr with the expr strings of nodes
func (expr *MetricExpr) GetNodeExpr(query *Query) {
	node := query.NodeID
	if node != "" {
		expr.CPU.Usage = fmt.Sprintf("avg(irate(container_cpu_usage_seconds_total{id='/',instance='%s:5014'}[5m])) by (instance)", node)
		expr.Memory.Usage = fmt.Sprintf("sum(container_memory_usage_bytes{id='/',instance='%s:5014'}) by (instance)", node)
		expr.Memory.Total = fmt.Sprintf("sum(container_spec_memory_limit_bytes{id='/',instance='%s:5014'}) by (instance)", node)
		expr.Network.Receive = fmt.Sprintf("irate(container_network_receive_bytes_total{id=~'/',instance='%s:5014'}[5m])", node)
		expr.Network.Transmit = fmt.Sprintf("irate(container_network_transmit_bytes_total{id=~'/',instance='%s:5014'}[5m])", node)
		expr.Filesystem.Usage = fmt.Sprintf("container_fs_usage_bytes{id=~'/',instance='%s:5014'}", node)
		expr.Filesystem.Limit = fmt.Sprintf("container_fs_limit_bytes{id=~'/',instance='%s:5014'}", node)
		return
	}
	expr.CPU.Usage = fmt.Sprintf("avg(irate(container_cpu_usage_seconds_total{id='/'}[5m])) by (instance)")
	expr.Memory.Usage = fmt.Sprintf("sum(container_memory_usage_bytes{id='/'}) by (instance)")
	expr.Memory.Total = fmt.Sprintf("sum(container_spec_memory_limit_bytes{id='/'}) by (instance)")
	expr.Network.Receive = fmt.Sprintf("irate(container_network_receive_bytes_total{id=~'/'}[5m])")
	expr.Network.Transmit = fmt.Sprintf("irate(container_network_transmit_bytes_total{id=~'/'}[5m])")
	expr.Filesystem.Usage = fmt.Sprintf("container_fs_usage_bytes{id=~'/'}")
	expr.Filesystem.Limit = fmt.Sprintf("container_fs_limit_bytes{id=~'/'}")
	return
}

// InfoExpr define the expr strings of Clusters/Cluster/Application
type InfoExpr struct {
	Clusters    string
	Cluster     string
	User        string
	Application string
}

// NewInfoExpr init the struct InfoExpr
func NewInfoExpr() *InfoExpr {
	return &InfoExpr{}
}

// GetInfoExpr set the value for InfoExpr fileds
func (expr *InfoExpr) GetInfoExpr(query *Query) {
	expr.Clusters = fmt.Sprintf("container_tasks_state{container_label_TASK_ID=~'[0-9]-.*',id=~'/docker/.*',name=~'mesos.*',state='running'}")
	expr.Cluster = fmt.Sprintf("container_tasks_state{container_label_TASK_ID=~'[0-9]-.*',id=~'/docker/.*',name=~'mesos.*',state='running',container_label_CLUSTER_ID='%s'}", query.ClusterID)
	expr.User = fmt.Sprintf("container_tasks_state{container_label_TASK_ID=~'[0-9]-.*',id=~'/docker/.*',name=~'mesos.*',state='running',container_label_CLUSTER_ID='%s', container_label_USER_ID='%s'}", query.ClusterID, query.UserID)
	expr.Application = fmt.Sprintf("container_tasks_state{container_label_TASK_ID=~'[0-9]-.*',id=~'/docker/.*',name=~'mesos.*', state='running',container_label_CLUSTER_ID='%s', container_label_USER_ID='%s', container_label_APP_ID='%s'}", query.ClusterID, query.UserID, query.AppID)
}
