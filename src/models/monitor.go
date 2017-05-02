package models

// QueryParameter define the fields of query paramter,
// Frontend pass the parameter in 2 way:
// 1. Expr, Start, End, Step
// 2. Metric, Cluster, User, App, Task, Start, End, Step
type QueryParameter struct {
	Expr    string
	Metric  string
	Cluster string
	User    string
	App     string
	Task    string
	Start   string
	End     string
	Step    string
}

// QueryExprResult define the JSON format
type QueryExprResult struct {
	Expr      string                 `json:"expr"`
	Status    string                 `json:"status"`
	Data      map[string]interface{} `json:"data"`
	ErrorType string                 `json:"errorType"`
	Error     string                 `json:"error"`
}

// QueryRangeResult define the JSON format of query result from prometheus
type QueryRangeResult struct {
	Expr      string `json:"expr"`
	Status    string `json:"status"`
	Data      *Data  `json:"data"`
	ErrorType string `json:"errorType"`
	Error     string `json:"error"`
}

// Data is the sub in QueryRangeResult
type Data struct {
	ResultType string    `json:"resultType"`
	Result     []*Result `json:"result"`
}

// Result is the sub in Data
type Result struct {
	Metric *Metric         `json:"metric"`
	Values [][]interface{} `json:"values"`
}

// Metric is the sub in Result
type Metric struct {
	ContainerLabelAppID    string `json:"container_label_DM_APP_ID,omitempty"`
	ContainerLabelApp      string `json:"container_label_DM_APP_NAME,omitempty"`
	ContainerLabelCluster  string `json:"container_label_DM_CLUSTER,omitempty"`
	ContainerLabelSlot     string `json:"container_label_DM_SLOT_INDEX,omitempty"`
	ContainerLabelSlotID   string `json:"container_label_DM_SLOT_ID,omitempty"`
	ContainerLabelTaskID   string `json:"container_label_DM_TASK_ID,omitempty"`
	ContainerLabelUser     string `json:"container_label_DM_USER_NAME,omitempty"`
	ContainerLabelGroup    string `json:"container_label_DM_GROUP_NAME,omitempty"`
	ContainerLabelVcluster string `json:"container_label_DM_VCLUSTER,omitempty"`
	ID                     string `json:"id"`
	Image                  string `json:"image"`
	Instance               string `json:"instance"`
	Job                    string `json:"job"`
	Name                   string `json:"name"`
	Group                  string `json:"group,omitempty"`
	Interface              string `json:"interface,omitempty"`
	Device                 string `json:"device,omitempty"`
}
