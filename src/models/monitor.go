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

// QueryExpreResult define the JSON format
type QueryExprResult struct {
	Expr      string                 `json:"expr"`
	Status    string                 `json:"status"`
	Data      map[string]interface{} `json:"data"`
	ErrorType string                 `json:"errorType"`
	Error     string                 `json:"error"`
}

// Metric is the sub in Result
type Metric struct {
	ContainerLabelAppID    string `json:"container_label_APP_ID,omitempty"`
	ContainerLabelApp      string `json:"container_label_APP,omitempty"`
	ContainerLabelCluster  string `json:"container_label_CLUSTER,omitempty"`
	ContainerLabelSlot     string `json:"container_label_SLOT,omitempty"`
	ContainerLabelTask     string `json:"container_env_mesos_task_id,omitempty"`
	ContainerLabelUser     string `json:"container_label_DM_USER,omitempty"`
	ContainerLabelVcluster string `json:"container_label_VCLUSTER,omitempty"`
	ID                     string `json:"id"`
	Image                  string `json:"image"`
	Instance               string `json:"instance"`
	Job                    string `json:"job"`
	Name                   string `json:"name"`
	Group                  string `json:"group,omitempty"`
	Interface              string `json:"interface,omitempty"`
	Device                 string `json:"device,omitempty"`
}

// MetricCPU defines the CPU metric to contain the Results
type MetricCPU struct {
	Usage []*Result `json:"usage"`
}

// MetricMemory defines the Memory metrics to contain the Results
type MetricMemory struct {
	Percentage []*Result `json:"usage"`
	Usage      []*Result `json:"usage_bytes"`
	Total      []*Result `json:"total_bytes"`
}

// MetricNewtork defines the Network metrics to contain the Results
type MetricNewtork struct {
	Receive  []*Result `json:"receive"`
	Transmit []*Result `json:"transmit"`
}

// MetricFilesystem defines the Filesystem metrics to contain the Results
type MetricFilesystem struct {
	Read  []*Result `json:"read"`
	Write []*Result `json:"write"`
}
