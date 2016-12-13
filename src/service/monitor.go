package service

import (
	"errors"
	"fmt"
	"strings"

	"github.com/Dataman-Cloud/log-proxy/src/backends"
	"github.com/Dataman-Cloud/log-proxy/src/models"
)

type Metric struct {
	CPU        *models.MetricCPU        `json:"cpu"`
	Memory     *models.MetricMemory     `json:"memory"`
	Network    *models.MetricNewtork    `json:"network"`
	Filesystem *models.MetricFilesystem `json:"filesystem"`
}

func NewMetric() *Metric {
	return &Metric{
		CPU:        &models.MetricCPU{},
		Memory:     &models.MetricMemory{},
		Network:    &models.MetricNewtork{},
		Filesystem: &models.MetricFilesystem{},
	}
}

func (m *Metric) GetQueryMetric(query *backends.Query) error {
	data, err := query.QueryMetric()
	if err != nil {
		return err
	}
	switch query.Metric {
	case "cpu":
		m.CPU.Usage = data.Data.Result
	case "memory":
		m.Memory.Percentage = data.Data.Result
	case "memory_usage":
		m.Memory.Usage = data.Data.Result
	case "memory_total":
		m.Memory.Total = data.Data.Result
	case "network_rx":
		m.Network.Receive = data.Data.Result
	case "network_tx":
		m.Network.Transmit = data.Data.Result
	case "fs_read":
		m.Filesystem.Read = data.Data.Result
	case "fs_write":
		m.Filesystem.Write = data.Data.Result
	default:
		return errors.New("No this kind of metric.")
	}
	return nil
}

type Info struct {
	Clusters     map[string]*ClusterInfo `json:"clusters"`
	Applications []string                `json:"applications"`
	Tasks        []string                `json:"tasks"`
	Nodes        []string                `json:"nodes"`
}

func NewInfo() *Info {
	return &Info{
		Clusters:     make(map[string]*ClusterInfo),
		Applications: []string{},
		Tasks:        []string{},
		Nodes:        []string{},
	}
}

type ClusterInfo struct {
	Applications map[string]*AppInfo `json:"applications"`
	Nodes        []string            `json:"nodes"`
	Tasks        []string            `json:"tasks"`
}

func NewClusterinfo() *ClusterInfo {
	return &ClusterInfo{
		Applications: make(map[string]*AppInfo),
		Tasks:        []string{},
		Nodes:        []string{},
	}
}

type AppInfo struct {
	CPU        *models.InfoCPU        `json:"cpu,omitempty"`
	Memory     *models.InfoMemory     `json:"memory,omitempty"`
	Network    *models.InfoNetwork    `json:"network,omitempty"`
	Filesystem *models.InfoFilesystem `json:"filesystem,omitempty"`
	Tasks      map[string]*TaskInfo   `json:"tasks"`
}

func NewAppInfo() *AppInfo {
	return &AppInfo{
		CPU:        &models.InfoCPU{},
		Memory:     &models.InfoMemory{},
		Network:    &models.InfoNetwork{},
		Filesystem: &models.InfoFilesystem{},
		Tasks:      make(map[string]*TaskInfo),
	}
}

type TaskInfo struct {
	Node       string                         `json:"node"`
	CPU        *models.InfoCPU                `json:"cpu"`
	Memory     *models.InfoMemory             `json:"memory"`
	Network    map[string]*models.InfoNetwork `json:"network"`
	Filesystem *models.InfoFilesystem         `json:"filesystem"`
}

func NewTaskInfo() *TaskInfo {
	return &TaskInfo{
		CPU:        &models.InfoCPU{},
		Memory:     &models.InfoMemory{},
		Network:    make(map[string]*models.InfoNetwork),
		Filesystem: &models.InfoFilesystem{},
	}
}

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
		metrics := []string{"cpu", "memory", "network_rx", "network_tx", "fs_read", "fs_write"}
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
				fmt.Printf("appid in cluster: %s\n", app)
				for _, cluster := range info.Clusters {
					for name, v := range cluster.Applications {
						if app == name {
							switch query.Metric {
							case "cpu":
								v.CPU.Usage = value
							case "memory":
								v.Memory.Percetange = value
							case "network_rx":
								v.Network.Receive = value
							case "network_tx":
								v.Network.Transmit = value
							case "fs_read":
								v.Filesystem.Read = value
							case "fs_write":
								v.Filesystem.Write = value
							}
						}
					}
				}
			}
		}

		//Fill the metric of tasks
		metrics = []string{"cpu", "memory_usage", "memory_total", "network_rx", "network_tx", "fs_read", "fs_write"}
		for _, metric := range metrics {
			query.Metric = metric
			data, err := query.QueryMetric()
			if err != nil {
				return err
			}

			if metric == "network_rx" || metric == "network_tx" {
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
								case "cpu":
									v.CPU.Usage = value
								case "memory_usage":
									v.Memory.Usage = value
								case "memory_total":
									v.Memory.Total = value
								case "network_rx":
									for nicK, nicV := range v.Network {
										if nic == nicK {
											nicV.Receive = value
										}
									}
								case "network_tx":
									for nicK, nicV := range v.Network {
										if nic == nicK {
											nicV.Transmit = value
										}
									}
								case "fs_read":
									v.Filesystem.Read = value
								case "fs_write":
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

type NodesInfo struct {
	Nodes map[string]*NodeInfo `json:"nodes"`
}

func NewNodesInfo() *NodesInfo {
	return &NodesInfo{
		Nodes: make(map[string]*NodeInfo),
	}
}

type NodeInfo struct {
	CPU        *models.InfoCPU                   `json:"cpu"`
	Memory     *models.InfoMemory                `json:"memory"`
	Network    map[string]*models.InfoNetwork    `json:"network"`
	Filesystem map[string]*models.InfoFilesystem `json:"filesystem"`
}

func NewNodeInfo() *NodeInfo {
	return &NodeInfo{
		CPU:        &models.InfoCPU{},
		Memory:     &models.InfoMemory{},
		Network:    make(map[string]*models.InfoNetwork),
		Filesystem: make(map[string]*models.InfoFilesystem),
	}
}

func (ni *NodesInfo) GetQueryNodesInfo(query *backends.Query) error {
	metrics := []string{"cpu", "memory_usage", "memory_total", "network_rx", "network_tx", "fs_usage", "fs_limit"}
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
		if metric == "network_rx" || metric == "network_tx" {
			for _, node := range ni.Nodes {
				if len(node.Network) == 0 {
					for _, originData := range data.Data.Result {
						nic := originData.Metric.Interface
						node.Network[nic] = models.NewInfoNetwork()
					}
				}
			}
		}

		if metric == "fs_usage" || metric == "fs_limit" {
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
					case "cpu":
						v.CPU.Usage = value
					case "memory_usage":
						v.Memory.Usage = value
					case "memory_total":
						v.Memory.Total = value
					case "network_rx":
						for nicK, nicV := range v.Network {
							if nic == nicK {
								nicV.Receive = value
							}
						}
					case "network_tx":
						for nicK, nicV := range v.Network {
							if nic == nicK {
								nicV.Transmit = value
							}
						}
					case "fs_usage":
						for fsK, fsV := range v.Filesystem {
							if device == fsK {
								fsV.Usage = value
							}
						}
					case "fs_limit":
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
