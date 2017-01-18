package backends

import (
	"net/http"
	"testing"
)

func initQueryMetric() *Query {
	param := &QueryParameter{
		Metric:  "cpu",
		Cluster: "work",
		App:     "web",
		Slot:    "0-1",
		User:    "user1",
		Start:   "1481612202",
		End:     "1481612212",
		Step:    "100s",
	}

	query := &Query{
		HTTPClient:     http.DefaultClient,
		PromServer:     "http://127.0.0.1:9090",
		Path:           QUERYRANGEPATH,
		QueryParameter: param,
	}

	return query
}

func expectGetTaskExpr() *MetricExpr {
	expr := NewMetricExpr()
	expr.CPU.Usage = "avg(irate(container_cpu_usage_seconds_total{container_label_CLUSTER='work'," +
		"container_label_APP='web',container_label_SLOT=~'0-1',container_label_USER='user1'," +
		"id=~'/docker/.*',name=~'mesos.*'}[5m])) by (container_label_CLUSTER,container_label_APP," +
		"container_label_SLOT,container_label_USER," +
		"group,id,image,instance,job,name,interface,device)"
	expr.Memory.Usage = "sum(container_memory_usage_bytes{container_label_CLUSTER='work'," +
		"container_label_APP='web',container_label_SLOT=~'0-1',container_label_USER='user1'," +
		"id=~'/docker/.*',name=~'mesos.*'}) by (container_label_CLUSTER,container_label_APP," +
		"container_label_SLOT,container_label_USER," +
		"group,id,image,instance,job,name,interface,device)"
	expr.Network.Receive = "sum(irate(container_network_receive_bytes_total{container_label_CLUSTER='work'," +
		"container_label_APP='web',container_label_SLOT=~'0-1',container_label_USER='user1'," +
		"id=~'/docker/.*',name=~'mesos.*'}[5m])) by (container_label_CLUSTER,container_label_APP," +
		"container_label_SLOT,container_label_USER," +
		"group,id,image,instance,job,name,interface,device)"
	expr.Filesystem.Read = "sum(irate(container_fs_reads_total{container_label_CLUSTER='work'," +
		"container_label_APP='web',container_label_SLOT=~'0-1',container_label_USER='user1'," +
		"id=~'/docker/.*',name=~'mesos.*'}[5m])) by (container_label_CLUSTER,container_label_APP," +
		"container_label_SLOT,container_label_USER," +
		"group,id,image,instance,job,name,interface,device)"
	return expr
}

func expectGetAppExpr(appid string) *MetricExpr {
	expr := NewMetricExpr()
	expr.CPU.Usage = "avg(irate(container_cpu_usage_seconds_total{container_label_CLUSTER='work',container_label_APP='web',container_label_SLOT=~'0-1',container_label_USER='user1',id=~'/docker/.*',name=~'mesos.*'}[5m])) by (container_label_CLUSTER,container_label_APP,container_label_USER)"
	expr.Memory.Usage = "sum(container_memory_usage_bytes{container_label_CLUSTER='work'," +
		"container_label_APP='web',container_label_SLOT=~'0-1',container_label_USER='user1'," +
		"id=~'/docker/.*',name=~'mesos.*'}) by (container_label_CLUSTER,container_label_APP," +
		"container_label_USER)"
	expr.Network.Receive = "sum(irate(container_network_receive_bytes_total{container_label_CLUSTER='work'," +
		"container_label_APP='web',container_label_SLOT=~'0-1',container_label_USER='user1'," +
		"id=~'/docker/.*',name=~'mesos.*'}[5m])) by (container_label_CLUSTER,container_label_APP," +
		"container_label_USER)"
	expr.Filesystem.Read = "sum(irate(container_fs_reads_total{container_label_CLUSTER='work'," +
		"container_label_APP='web',container_label_SLOT=~'0-1',container_label_USER='user1'," +
		"id=~'/docker/.*',name=~'mesos.*'}[5m])) by (container_label_CLUSTER,container_label_APP," +
		"container_label_USER)"

	return expr
}

func expectGetClusterExpr(node string) *MetricExpr {
	expr := NewMetricExpr()
	if node == "" {
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

func TestGetExpr(t *testing.T) {
	query := initQueryMetric()
	level := "task"
	expr := expectGetTaskExpr()

	newExpr := GetExpr(query, level)
	if expr.CPU.Usage != newExpr.CPU.Usage {
		t.Errorf("Expect %v\n, got wrong expr %v", expr.CPU.Usage, newExpr.CPU.Usage)
	}

	if expr.Memory.Usage != newExpr.Memory.Usage {
		t.Errorf("Expect %v\n, got wrong expr %v", expr.Memory.Usage, newExpr.Memory.Usage)
	}

	if expr.Network.Receive != newExpr.Network.Receive {
		t.Errorf("Expect %v, got wrong expr %v", expr.Network.Receive, newExpr.Network.Receive)
	}

	if expr.Filesystem.Read != newExpr.Filesystem.Read {
		t.Errorf("Expect %v, got wrong expr %v", expr.Filesystem.Read, newExpr.Filesystem.Read)
	}

	level = "app"
	expr = expectGetAppExpr(query.App)
	newExpr = GetExpr(query, level)
	if expr.CPU.Usage != newExpr.CPU.Usage {
		t.Errorf("Expect %v\n, got wrong expr %v", expr.CPU.Usage, newExpr.CPU.Usage)
	}

	if expr.Memory.Usage != newExpr.Memory.Usage {
		t.Errorf("Expect %v, got wrong expr %v", expr.Memory.Usage, newExpr.Memory.Usage)
	}

	if expr.Network.Receive != newExpr.Network.Receive {
		t.Errorf("Expect %v, got wrong expr %v", expr.Network.Receive, newExpr.Network.Receive)
	}

	if expr.Filesystem.Read != newExpr.Filesystem.Read {
		t.Errorf("Expect %v, got wrong expr %v", expr.Filesystem.Read, newExpr.Filesystem.Read)
	}

	level = "cluster"
	expr = expectGetClusterExpr(query.Node)
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
	query.Node = "192.168.1.1"
	expr = expectGetClusterExpr(query.Node)
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

	appexpr := "container_tasks_state{container_label_license!='GPLv2',id=~'/docker/.*',name=~'mesos.*', state='running',container_label_CLUSTER='work', container_label_USER='user1', container_label_APP='web'}"
	if appexpr != infoExpr.Application {
		t.Errorf("newExpr is Not right. Expect %s, got wrong expr %s", appexpr, infoExpr.Application)
	}

	clusterexpr := "container_tasks_state{container_label_license!='GPLv2',id=~'/docker/.*',name=~'mesos.*',state='running',container_label_CLUSTER='work'}"
	if clusterexpr != infoExpr.Cluster {
		t.Errorf("newExpr is Not right. Expect %s, got wrong expr %s", clusterexpr, infoExpr.Cluster)
	}

	clustersexpr := "container_tasks_state{container_label_license!='GPLv2',id=~'/docker/.*',name=~'mesos.*',state='running'}"
	if clustersexpr != infoExpr.Clusters {
		t.Errorf("newExpr is Not right. Expect %s, got wrong expr %s", clustersexpr, infoExpr.Clusters)
	}

}
