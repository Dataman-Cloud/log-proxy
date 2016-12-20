package service

import (
	"errors"
	"strings"

	"github.com/Dataman-Cloud/log-proxy/src/backends"
	"github.com/Dataman-Cloud/log-proxy/src/models"
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
		return errors.New("No this kind of metric.")
	}
	return nil
}

// Info defines the JSON format of information in clusters/cluster/apps/nodes
type Info struct {
	Clusters     map[string]*ClusterInfo `json:"clusters"`
	Applications []string                `json:"applications"`
	Tasks        []string                `json:"tasks"`
	Nodes        []string                `json:"nodes"`
}

// NewInfo init struct Info
func NewInfo() *Info {
	return &Info{
		Clusters:     make(map[string]*ClusterInfo),
		Applications: []string{},
		Tasks:        []string{},
		Nodes:        []string{},
	}
}

// ClusterInfo defines the JSON format of information in cluster
type ClusterInfo struct {
	Applications map[string]*AppInfo `json:"applications"`
	Nodes        []string            `json:"nodes"`
	Tasks        []string            `json:"tasks"`
}

// NewClusterinfo init the struct ClusterInfo
func NewClusterinfo() *ClusterInfo {
	return &ClusterInfo{
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

// GetQueryInfo get result by calling QueryInfo
// then set the vaules of fileds in Info
func (info *Info) GetQueryInfo(query *backends.Query) error {
	data, err := query.QueryInfo()
	if err != nil {
		return err
	}

	// Fill the info of Clusters
	if len(info.Clusters) == 0 {
		for _, originData := range data.Data.Result {
			name := originData.Metric.ContainerLabelVCLUSTER
			app := originData.Metric.ContainerLabelAPPID
			task := originData.Metric.ID
			node := strings.Split(originData.Metric.Instance, ":")[0]
			info.Clusters[name] = NewClusterinfo()
			if !isInArray(info.Applications, app) {
				info.Applications = append(info.Applications, app)
			}
			if !isInArray(info.Tasks, task) {
				info.Tasks = append(info.Tasks, task)
			}
			if !isInArray(info.Nodes, node) {
				info.Nodes = append(info.Nodes, node)
			}
		}
	}
	// Fill the info of cluster
	for _, originData := range data.Data.Result {
		cluster := originData.Metric.ContainerLabelVCLUSTER
		app := originData.Metric.ContainerLabelAPPID
		task := originData.Metric.ID
		node := strings.Split(originData.Metric.Instance, ":")[0]
		for name, value := range info.Clusters {
			if cluster == name {
				value.Applications[app] = NewAppInfo()
				if !isInArray(value.Tasks, task) {
					value.Tasks = append(value.Tasks, task)
				}
				if !isInArray(value.Nodes, node) {
					value.Nodes = append(value.Nodes, node)
				}
			}
		}
	}

	// Fill the info of app
	for _, originData := range data.Data.Result {
		app := originData.Metric.ContainerLabelAPPID
		task := originData.Metric.ID
		for _, cluster := range info.Clusters {
			for name, value := range cluster.Applications {
				if app == name {
					value.Tasks[task] = NewTaskInfo()
				}
			}
		}
	}

	// Fill the metric of app
	if query.ClusterID != "" || query.AppID != "" {
		metrics := []string{cpuUsage, memPercentage, networkRX, networkTX, fsRead, fsWrite}
		for _, metric := range metrics {
			query.Metric = metric
			query.Path = backends.QUERYRANGEPATH
			data, err := query.QueryAppMetric()
			if err != nil {
				return err
			}
			for _, originData := range data.Data.Result {
				app := originData.Metric.ContainerLabelAPPID
				value := originData.Values[0]
				for _, cluster := range info.Clusters {
					for name, v := range cluster.Applications {
						if app == name {
							switch query.Metric {
							case cpuUsage:
								v.CPU.Usage = value
							case memPercentage:
								v.Memory.Percetange = value
							case networkRX:
								v.Network.Receive = value
							case networkTX:
								v.Network.Transmit = value
							case fsRead:
								v.Filesystem.Read = value
							case fsWrite:
								v.Filesystem.Write = value
							}
						}
					}
				}
			}
		}

		//Fill the metric of tasks
		metrics = []string{cpuUsage, memUsage, memTotal, networkRX, networkTX, fsRead, fsWrite}
		for _, metric := range metrics {
			query.Metric = metric
			data, err := query.QueryMetric()
			if err != nil {
				return err
			}

			if metric == networkRX || metric == networkTX {
				for _, cluster := range info.Clusters {
					for _, app := range cluster.Applications {
						for _, task := range app.Tasks {
							if len(task.Network) == 0 {
								for _, originData := range data.Data.Result {
									nic := originData.Metric.Interface
									task.Network[nic] = models.NewInfoNetwork()
								}
							}
						}
					}
				}
			}
			for _, originData := range data.Data.Result {
				task := originData.Metric.ID
				nic := originData.Metric.Interface
				node := strings.Split(originData.Metric.Instance, ":")[0]
				value := originData.Values[0]
				for _, cluster := range info.Clusters {
					for _, app := range cluster.Applications {
						for name, v := range app.Tasks {
							if task == name {
								v.Node = node
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
								case fsRead:
									v.Filesystem.Read = value
								case fsWrite:
									v.Filesystem.Write = value
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
