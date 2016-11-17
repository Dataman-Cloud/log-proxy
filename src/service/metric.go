package service

import (
	"errors"

	"github.com/Dataman-Cloud/log-proxy/src/models"
)

type MetricList struct {
	CPU        *MetricCPU        `json:"cpu"`
	Memory     *MetricMemory     `json:"memory"`
	Network    *MetricNewtork    `json:"network"`
	Filesystem *MetricFilesystem `json:"filesystem"`
}

func NewMetricList() *MetricList {
	return &MetricList{
		CPU:        &MetricCPU{},
		Memory:     &MetricMemory{},
		Network:    &MetricNewtork{},
		Filesystem: &MetricFilesystem{},
	}
}

type MetricCPU struct {
	Usage []*models.Result `json:"usage"`
}

type MetricMemory struct {
	Usage []*models.Result `json:"usage"`
}

type MetricNewtork struct {
	Receive  []*models.Result `json:"receive"`
	Transmit []*models.Result `json:"transmit"`
}

type MetricFilesystem struct {
	Read  []*models.Result `json:"read"`
	Write []*models.Result `json:"write"`
}

// GetMetricList will return the list of metric of query result.
// metric=cpu, will only query CPU usage
// metric=memory, will only query memory usage
// metric=all, will query each metic.
func (ml *MetricList) GetMetricList(query *QueryRange) error {
	if query.Metric == "all" {
		metrics := [...]string{"cpu", "memory", "network_rx", "network_tx", "fs_read", "fs_write"}
		for _, metric := range metrics {
			query.Metric = metric
			err := ml.SetMetricList(query)
			if err != nil {
				return err
			}
		}
	}

	err := ml.SetMetricList(query)
	if err != nil {
		return err
	}
	return nil
}

func (ml *MetricList) SetMetricList(query *QueryRange) error {
	data, err := query.QueryRangeFromProm()
	if err != nil {
		return err
	}

	switch query.Metric {
	case "cpu":
		ml.CPU.Usage = data.Data.Result
	case "memory":
		ml.Memory.Usage = data.Data.Result
	case "network_rx":
		ml.Network.Receive = data.Data.Result
	case "network_tx":
		ml.Network.Transmit = data.Data.Result
	case "fs_read":
		ml.Filesystem.Read = data.Data.Result
	case "fs_write":
		ml.Filesystem.Write = data.Data.Result
	default:
		return errors.New("No this kind of metric.")
	}
	return nil
}

type AppsList struct {
	Apps map[string][]string `json:"apps"`
}

func NewAppsList() *AppsList {
	return &AppsList{
		Apps: make(map[string][]string),
	}
}

func (a *AppsList) GetAppsList(query *QueryRange) error {
	err := a.SetAppsList(query)
	if err != nil {
		return err
	}
	return nil
}

func (a *AppsList) SetAppsList(query *QueryRange) error {
	data, err := query.QueryAppsFromProm()
	if err != nil {
		return err
	}

	// appIDs save the APPID in list;
	// Get the appID and taskID, make the appID as Key, the array of taskID as
	// value, and then append them into the list of apps
	appIDs := []string{}
	for _, originData := range data.Data.Result {
		appID := originData.Metric.ContainerLabelAPPID
		taskID := originData.Metric.ID

		if len(a.Apps) == 0 {
			appIDs = append(appIDs, appID)
		}

		isExist := false
		for _, value := range appIDs {
			if appID == value {
				a.Apps[appID] = append(a.Apps[appID], taskID)
				isExist = true
				break
			}
		}
		if !isExist {
			tasks := []string{taskID}
			a.Apps[appID] = tasks
			appIDs = append(appIDs, appID)
		}
	}
	return nil
}
