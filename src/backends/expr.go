package backends

import "fmt"

type MetricExpr struct {
	CPU        *MetricExprCPU
	Memory     *MetricExprMemory
	Network    *MetricExprNewtork
	Filesystem *MetricExprFilesystem
}

func NewMetricExpr() *MetricExpr {
	return &MetricExpr{
		CPU:        &MetricExprCPU{},
		Memory:     &MetricExprMemory{},
		Network:    &MetricExprNewtork{},
		Filesystem: &MetricExprFilesystem{},
	}
}

type MetricExprCPU struct {
	Usage string
}

type MetricExprMemory struct {
	Percentage string
	Usage      string
	Total      string
}

type MetricExprNewtork struct {
	Receive  string
	Transmit string
}

type MetricExprFilesystem struct {
	Read  string
	Write string
	Usage string
	Limit string
}

func GetExpr(query *Query, level string) *MetricExpr {
	me := NewMetricExpr()
	switch level {
	case "id":
		me.GetMetricExpr(query)
	case "app":
		me.GetAppExpr(query)
	case "cluster":
		me.GetNodeExpr(query)
	}
	return me
}

func (expr *MetricExpr) GetMetricExpr(query *Query) {
	byItems := "container_label_VCLUSTER, container_label_APP_ID, group, id, image, instance, job, name, interface, device"
	expr.setExpr(query, byItems)
}

func (expr *MetricExpr) GetAppExpr(query *Query) {
	byItems := "container_label_VCLUSTER, container_label_APP_ID"
	expr.setExpr(query, byItems)
}

func (expr *MetricExpr) setExpr(query *Query, byItems string) {
	expr.CPU.Usage = fmt.Sprintf("avg(irate(container_cpu_usage_seconds_total{container_label_APP_ID=~'.*',"+
		"id=~'/docker/%s.*', name=~'mesos.*'}[5m])) by (%s)",
		query.TaskID, byItems)

	expr.Memory.Usage = fmt.Sprintf("sum(container_memory_usage_bytes{container_label_APP_ID=~'.*',"+
		"id=~'/docker/%s.*', name=~'mesos.*'}) by (%s)", query.TaskID, byItems)

	expr.Memory.Total = fmt.Sprintf("sum(container_spec_memory_limit_bytes{container_label_APP_ID=~'.*',"+
		"id=~'/docker/%s.*', name=~'mesos.*'}) by (%s)", query.TaskID, byItems)

	expr.Memory.Percentage = fmt.Sprintf("%s / %s", expr.Memory.Usage, expr.Memory.Total)

	expr.Network.Receive = fmt.Sprintf("sum(irate(container_network_receive_bytes_total{container_label_APP_ID=~'.*',"+
		"id=~'/docker/%s.*', name=~'mesos.*'}[5m])) by (%s)", query.TaskID, byItems)

	expr.Network.Transmit = fmt.Sprintf("sum(irate(container_network_transmit_bytes_total{container_label_APP_ID=~'.*',"+
		"id=~'/docker/%s.*', name=~'mesos.*'}[5m])) by (%s)", query.TaskID, byItems)

	expr.Filesystem.Read = fmt.Sprintf("sum(irate(container_fs_reads_total{container_label_APP_ID=~'.*',"+
		"id=~'/docker/%s.*', name=~'mesos.*'}[5m])) by (%s)", query.TaskID, byItems)

	expr.Filesystem.Write = fmt.Sprintf("sum(irate(container_fs_writes_total{container_label_APP_ID=~'.*',"+
		"id=~'/docker/%s.*', name=~'mesos.*'}[5m])) by (%s)", query.TaskID, byItems)
}

func (expr *MetricExpr) GetNodeExpr(query *Query) {
	node := query.NodeID
	if node != "" {
		expr.CPU.Usage = fmt.Sprintf("avg(irate(container_cpu_usage_seconds_total{id='/',instance='%s:5014'}[5m])) by (instance)", node)
		expr.Memory.Usage = fmt.Sprintf("sum(container_memory_usage_bytes{id='/',instance='%s:5014'}) by (instance)", node)
		expr.Memory.Total = fmt.Sprintf("sum(container_spec_memory_limit_bytes{id='/',instance='%s:5014'}) by (instance)", node)
		expr.Network.Receive = fmt.Sprintf("irate(container_network_receive_bytes_total{id=~'/',instance='%s:5014'}[5m])", node)
		expr.Network.Transmit = fmt.Sprintf("irate(container_network_transmit_bytes_total{id=~'/',instance='%s:5014'}[5m])", node)
		expr.Filesystem.Usage = fmt.Sprintf("irate(container_fs_usage_bytes{id=~'/',instance='%s:5014'}[5m])", node)
		expr.Filesystem.Limit = fmt.Sprintf("irate(container_fs_limit_bytes{id=~'/',instance='%s:5014'}[5m])", node)
		return
	}
	expr.CPU.Usage = fmt.Sprintf("avg(irate(container_cpu_usage_seconds_total{id='/'}[5m])) by (instance)")
	expr.Memory.Usage = fmt.Sprintf("sum(container_memory_usage_bytes{id='/'}) by (instance)")
	expr.Memory.Total = fmt.Sprintf("sum(container_spec_memory_limit_bytes{id='/'}) by (instance)")
	expr.Network.Receive = fmt.Sprintf("irate(container_network_receive_bytes_total{id=~'/'}[5m])")
	expr.Network.Transmit = fmt.Sprintf("irate(container_network_transmit_bytes_total{id=~'/'}[5m])")
	expr.Filesystem.Usage = fmt.Sprintf("irate(container_fs_usage_bytes{id=~'/'}[5m])")
	expr.Filesystem.Limit = fmt.Sprintf("irate(container_fs_limit_bytes{id=~'/'}[5m])")
	return
}

type InfoExpr struct {
	Clusters    string
	Cluster     string
	Application string
}

func NewInfoExpr() *InfoExpr {
	return &InfoExpr{}
}

func (expr *InfoExpr) GetInfoExpr(query *Query) {
	expr.Clusters = fmt.Sprintf("container_tasks_state{id=~'/docker/.*', name=~'mesos.*', state='running'}")
	expr.Cluster = fmt.Sprintf("container_tasks_state{id=~'/docker/.*', name=~'mesos.*', state='running',container_label_VCLUSTER='%s'}", query.ClusterID)
	expr.Application = fmt.Sprintf(fmt.Sprintf("container_tasks_state{id=~'/docker/.*', name=~'mesos.*', state='running',container_label_APP_ID='%s'}", query.AppID))
}
