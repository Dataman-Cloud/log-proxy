package backends

import (
	"net/http"
	"testing"
)

func initQueryMetric() *Query {
	param := &QueryParameter{
		Metric:    "cpu",
		ClusterID: "work",
		AppID:     "work-web",
		TaskID:    "abcdefg",
		Start:     "1481612202",
		End:       "1481612212",
		Step:      "100s",
	}

	query := &Query{
		HTTPClient:     http.DefaultClient,
		PromServer:     "http://127.0.0.1:9090",
		Path:           QUERYRANGEPATH,
		QueryParameter: param,
	}

	return query
}

func testGetTaskExpr() *MetricExpr {
	expr := NewMetricExpr()
	expr.CPU.Usage = "avg(irate(container_cpu_usage_seconds_total{container_label_APP_ID='work-web',id=~'/docker/abcdefg.*', name=~'mesos.*'}[5m])) by (container_label_VCLUSTER, container_label_APP_ID, group, id, image, instance, job, name, interface, device)"
	expr.Memory.Percentage = "sum(container_memory_usage_bytes{container_label_APP_ID='work-web',id=~'/docker/abcdefg.*', name=~'mesos.*'}) by (container_label_VCLUSTER, container_label_APP_ID, group, id, image, instance, job, name, interface, device) / sum(container_spec_memory_limit_bytes{container_label_APP_ID='work-web',id=~'/docker/abcdefg.*', name=~'mesos.*'}) by (container_label_VCLUSTER, container_label_APP_ID, group, id, image, instance, job, name, interface, device)"
	expr.Memory.Usage = "sum(container_memory_usage_bytes{container_label_APP_ID='work-web',id=~'/docker/abcdefg.*', name=~'mesos.*'}) by (container_label_VCLUSTER, container_label_APP_ID, group, id, image, instance, job, name, interface, device)"
	expr.Memory.Total = "sum(container_spec_memory_limit_bytes{container_label_APP_ID='work-web',id=~'/docker/abcdefg.*', name=~'mesos.*'}) by (container_label_VCLUSTER, container_label_APP_ID, group, id, image, instance, job, name, interface, device)"
	expr.Network.Receive = "sum(irate(container_network_receive_bytes_total{container_label_APP_ID='work-web',id=~'/docker/abcdefg.*', name=~'mesos.*'}[5m])) by (container_label_VCLUSTER, container_label_APP_ID, group, id, image, instance, job, name, interface, device)"
	expr.Network.Transmit = "sum(irate(container_network_transmit_bytes_total{container_label_APP_ID='work-web',id=~'/docker/abcdefg.*', name=~'mesos.*'}[5m])) by (container_label_VCLUSTER, container_label_APP_ID, group, id, image, instance, job, name, interface, device)"
	expr.Filesystem.Read = "sum(irate(container_fs_reads_total{container_label_APP_ID='work-web',id=~'/docker/abcdefg.*', name=~'mesos.*'}[5m])) by (container_label_VCLUSTER, container_label_APP_ID, group, id, image, instance, job, name, interface, device)"
	expr.Filesystem.Write = "sum(irate(container_fs_writes_total{container_label_APP_ID='work-web',id=~'/docker/abcdefg.*', name=~'mesos.*'}[5m])) by (container_label_VCLUSTER, container_label_APP_ID, group, id, image, instance, job, name, interface, device)"
	return expr
}

func testGetAppExpr(appid string) *MetricExpr {
	expr := NewMetricExpr()
	if appid == "" {
		expr.CPU.Usage = "avg(irate(container_cpu_usage_seconds_total{container_label_APP_ID=~'.*',id=~'/docker/abcdefg.*', name=~'mesos.*'}[5m])) by (container_label_VCLUSTER, container_label_APP_ID)"
		expr.Memory.Usage = "sum(container_memory_usage_bytes{container_label_APP_ID=~'.*',id=~'/docker/abcdefg.*', name=~'mesos.*'}) by (container_label_VCLUSTER, container_label_APP_ID)"
		expr.Network.Receive = "sum(irate(container_network_receive_bytes_total{container_label_APP_ID=~'.*',id=~'/docker/abcdefg.*', name=~'mesos.*'}[5m])) by (container_label_VCLUSTER, container_label_APP_ID)"
		expr.Network.Transmit = "sum(irate(container_network_transmit_bytes_total{container_label_APP_ID=~'.*',id=~'/docker/abcdefg.*', name=~'mesos.*'}[5m])) by (container_label_VCLUSTER, container_label_APP_ID)"
		expr.Filesystem.Read = "sum(irate(container_fs_reads_total{container_label_APP_ID=~'.*',id=~'/docker/abcdefg.*', name=~'mesos.*'}[5m])) by (container_label_VCLUSTER, container_label_APP_ID)"
		expr.Filesystem.Write = "sum(irate(container_fs_writes_total{container_label_APP_ID=~'.*',id=~'/docker/abcdefg.*', name=~'mesos.*'}[5m])) by (container_label_VCLUSTER, container_label_APP_ID)"
		return expr
	}
	expr.CPU.Usage = "avg(irate(container_cpu_usage_seconds_total{container_label_APP_ID='work-web',id=~'/docker/abcdefg.*', name=~'mesos.*'}[5m])) by (container_label_VCLUSTER, container_label_APP_ID)"
	expr.Memory.Usage = "sum(container_memory_usage_bytes{container_label_APP_ID='work-web',id=~'/docker/abcdefg.*', name=~'mesos.*'}) by (container_label_VCLUSTER, container_label_APP_ID)"
	expr.Network.Receive = "sum(irate(container_network_receive_bytes_total{container_label_APP_ID='work-web',id=~'/docker/abcdefg.*', name=~'mesos.*'}[5m])) by (container_label_VCLUSTER, container_label_APP_ID)"
	expr.Network.Transmit = "sum(irate(container_network_transmit_bytes_total{container_label_APP_ID='work-web',id=~'/docker/abcdefg.*', name=~'mesos.*'}[5m])) by (container_label_VCLUSTER, container_label_APP_ID)"
	expr.Filesystem.Read = "sum(irate(container_fs_reads_total{container_label_APP_ID='work-web',id=~'/docker/abcdefg.*', name=~'mesos.*'}[5m])) by (container_label_VCLUSTER, container_label_APP_ID)"
	expr.Filesystem.Write = "sum(irate(container_fs_writes_total{container_label_APP_ID='work-web',id=~'/docker/abcdefg.*', name=~'mesos.*'}[5m])) by (container_label_VCLUSTER, container_label_APP_ID)"
	return expr
}

func testGetClusterExpr(nodeid string) *MetricExpr {
	expr := NewMetricExpr()
	if nodeid == "" {
		expr.CPU.Usage = "avg(irate(container_cpu_usage_seconds_total{id='/'}[5m])) by (instance)"
		expr.Memory.Usage = "sum(container_memory_usage_bytes{id='/'}) by (instance)"
		expr.Memory.Total = "sum(container_spec_memory_limit_bytes{id='/'}) by (instance)"
		expr.Network.Receive = "irate(container_network_receive_bytes_total{id=~'/'}[5m])"
		expr.Network.Transmit = "irate(container_network_transmit_bytes_total{id=~'/'}[5m])"
		expr.Filesystem.Usage = "container_fs_usage_bytes{id=~'/'}"
		expr.Filesystem.Limit = "container_fs_limit_bytes{id=~'/'}"
		return expr
	}
	expr.CPU.Usage = "avg(irate(container_cpu_usage_seconds_total{id='/',instance='192.168.1.1:5014'}[5m])) by (instance)"
	expr.Memory.Usage = "sum(container_memory_usage_bytes{id='/',instance='192.168.1.1:5014'}) by (instance)"
	expr.Memory.Total = "sum(container_spec_memory_limit_bytes{id='/',instance='192.168.1.1:5014'}) by (instance)"
	expr.Network.Receive = "irate(container_network_receive_bytes_total{id=~'/',instance='192.168.1.1:5014'}[5m])"
	expr.Network.Transmit = "irate(container_network_transmit_bytes_total{id=~'/',instance='192.168.1.1:5014'}[5m])"
	expr.Filesystem.Usage = "container_fs_usage_bytes{id=~'/',instance='192.168.1.1:5014'}"
	expr.Filesystem.Limit = "container_fs_limit_bytes{id=~'/',instance='192.168.1.1:5014'}"
	return expr
}

func TestGetCPUExpr(t *testing.T) {
	query := initQueryMetric()
	level := "id"
	expr := testGetTaskExpr()

	newExpr := GetExpr(query, level)
	if expr.CPU.Usage != newExpr.CPU.Usage {
		t.Errorf("Expect %v, got wrong expr %v", expr.CPU.Usage, newExpr.CPU.Usage)
	}
	if expr.Memory.Usage != newExpr.Memory.Usage {
		t.Errorf("Expect %v, got wrong expr %v", expr.Memory.Usage, newExpr.Memory.Usage)
	}
	if expr.Network.Receive != newExpr.Network.Receive {
		t.Errorf("Expect %v, got wrong expr %v", expr.Network.Receive, newExpr.Network.Receive)
	}
	if expr.Network.Transmit != newExpr.Network.Transmit {
		t.Errorf("Expect %v, got wrong expr %v", expr.Network.Transmit, newExpr.Network.Transmit)
	}
	if expr.Filesystem.Read != newExpr.Filesystem.Read {
		t.Errorf("Expect %v, got wrong expr %v", expr.Filesystem.Read, newExpr.Filesystem.Read)
	}
	if expr.Filesystem.Write != newExpr.Filesystem.Write {
		t.Errorf("Expect %v, got wrong expr %v", expr.Filesystem.Write, newExpr.Filesystem.Write)
	}

	level = "app"
	expr = testGetAppExpr(query.AppID)
	newExpr = GetExpr(query, level)
	if expr.CPU.Usage != newExpr.CPU.Usage {
		t.Errorf("Expect %v, got wrong expr %v", expr.CPU.Usage, newExpr.CPU.Usage)
	}
	if expr.Memory.Usage != newExpr.Memory.Usage {
		t.Errorf("Expect %v, got wrong expr %v", expr.Memory.Usage, newExpr.Memory.Usage)
	}
	if expr.Network.Receive != newExpr.Network.Receive {
		t.Errorf("Expect %v, got wrong expr %v", expr.Network.Receive, newExpr.Network.Receive)
	}
	if expr.Network.Transmit != newExpr.Network.Transmit {
		t.Errorf("Expect %v, got wrong expr %v", expr.Network.Transmit, newExpr.Network.Transmit)
	}
	if expr.Filesystem.Read != newExpr.Filesystem.Read {
		t.Errorf("Expect %v, got wrong expr %v", expr.Filesystem.Read, newExpr.Filesystem.Read)
	}
	if expr.Filesystem.Write != newExpr.Filesystem.Write {
		t.Errorf("Expect %v, got wrong expr %v", expr.Filesystem.Write, newExpr.Filesystem.Write)
	}
	query.AppID = ""
	expr = testGetAppExpr(query.AppID)
	newExpr = GetExpr(query, level)
	if expr.CPU.Usage != newExpr.CPU.Usage {
		t.Errorf("Expect %v, got wrong expr %v", expr.CPU.Usage, newExpr.CPU.Usage)
	}
	if expr.Memory.Usage != newExpr.Memory.Usage {
		t.Errorf("Expect %v, got wrong expr %v", expr.Memory.Usage, newExpr.Memory.Usage)
	}
	if expr.Network.Receive != newExpr.Network.Receive {
		t.Errorf("Expect %v, got wrong expr %v", expr.Network.Receive, newExpr.Network.Receive)
	}
	if expr.Network.Transmit != newExpr.Network.Transmit {
		t.Errorf("Expect %v, got wrong expr %v", expr.Network.Transmit, newExpr.Network.Transmit)
	}
	if expr.Filesystem.Read != newExpr.Filesystem.Read {
		t.Errorf("Expect %v, got wrong expr %v", expr.Filesystem.Read, newExpr.Filesystem.Read)
	}
	if expr.Filesystem.Write != newExpr.Filesystem.Write {
		t.Errorf("Expect %v, got wrong expr %v", expr.Filesystem.Write, newExpr.Filesystem.Write)
	}

	level = "cluster"
	expr = testGetClusterExpr(query.NodeID)
	newExpr = GetExpr(query, level)
	if expr.CPU.Usage != newExpr.CPU.Usage {
		t.Errorf("Expect %v, got wrong expr %v", expr.CPU.Usage, newExpr.CPU.Usage)
	}
	if expr.Memory.Usage != newExpr.Memory.Usage {
		t.Errorf("Expect %v, got wrong expr %v", expr.Memory.Usage, newExpr.Memory.Usage)
	}
	if expr.Network.Receive != newExpr.Network.Receive {
		t.Errorf("Expect %v, got wrong expr %v", expr.Network.Receive, newExpr.Network.Receive)
	}
	if expr.Network.Transmit != newExpr.Network.Transmit {
		t.Errorf("Expect %v, got wrong expr %v", expr.Network.Transmit, newExpr.Network.Transmit)
	}
	if expr.Filesystem.Read != newExpr.Filesystem.Read {
		t.Errorf("Expect %v, got wrong expr %v", expr.Filesystem.Read, newExpr.Filesystem.Read)
	}
	if expr.Filesystem.Write != newExpr.Filesystem.Write {
		t.Errorf("Expect %v, got wrong expr %v", expr.Filesystem.Write, newExpr.Filesystem.Write)
	}
	query.NodeID = "192.168.1.1"
	expr = testGetClusterExpr(query.NodeID)
	newExpr = GetExpr(query, level)
	if expr.CPU.Usage != newExpr.CPU.Usage {
		t.Errorf("Expect %v, got wrong expr %v", expr.CPU.Usage, newExpr.CPU.Usage)
	}
	if expr.Memory.Usage != newExpr.Memory.Usage {
		t.Errorf("Expect %v, got wrong expr %v", expr.Memory.Usage, newExpr.Memory.Usage)
	}
	if expr.Network.Receive != newExpr.Network.Receive {
		t.Errorf("Expect %v, got wrong expr %v", expr.Network.Receive, newExpr.Network.Receive)
	}
	if expr.Network.Transmit != newExpr.Network.Transmit {
		t.Errorf("Expect %v, got wrong expr %v", expr.Network.Transmit, newExpr.Network.Transmit)
	}
	if expr.Filesystem.Read != newExpr.Filesystem.Read {
		t.Errorf("Expect %v, got wrong expr %v", expr.Filesystem.Read, newExpr.Filesystem.Read)
	}
	if expr.Filesystem.Write != newExpr.Filesystem.Write {
		t.Errorf("Expect %v, got wrong expr %v", expr.Filesystem.Write, newExpr.Filesystem.Write)
	}
}

func TestInfoExpr(t *testing.T) {
	query := initQueryMetric()
	infoExpr := NewInfoExpr()
	infoExpr.GetInfoExpr(query)

	appexpr := "container_tasks_state{id=~'/docker/.*', name=~'mesos.*', state='running',container_label_APP_ID='work-web'}"
	if appexpr != infoExpr.Application {
		t.Errorf("newExpr is Not right. Expect %s, got wrong expr %s", appexpr, infoExpr.Application)
	}

	clusterexpr := "container_tasks_state{id=~'/docker/.*', name=~'mesos.*', state='running',container_label_VCLUSTER='work'}"
	if clusterexpr != infoExpr.Cluster {
		t.Errorf("newExpr is Not right. Expect %s, got wrong expr %s", clusterexpr, infoExpr.Cluster)
	}

	clustersexpr := "container_tasks_state{id=~'/docker/.*', name=~'mesos.*', state='running'}"
	if clustersexpr != infoExpr.Clusters {
		t.Errorf("newExpr is Not right. Expect %s, got wrong expr %s", clustersexpr, infoExpr.Clusters)
	}
}
