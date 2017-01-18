package models

// QueryRangeResult define the JSON format of query result from prometheus
type QueryRangeResult struct {
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
	ContainerLabelApp      string `json:"container_label_APP"`
	ContainerLabelCluster  string `json:"container_label_CLUSTER"`
	ContainerLabelSlot     string `json:"container_label_SLOT"`
	ContainerLabelTask     string `json:"container_label_TASK"`
	ContainerLabelUser     string `json:"container_label_USER"`
	ContainerLabelVcluster string `json:"container_label_VCLUSTER"`
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

// InfoCPU defines the CPU metric to contain the Results in Info
type InfoCPU struct {
	Usage []interface{} `json:"usage"`
}

// InfoMemory defines the Memory metric to contain the Results in Info
type InfoMemory struct {
	Percetange []interface{} `json:"usage"`
	Usage      []interface{} `json:"usage_bytes"`
	Total      []interface{} `json:"total_bytes"`
}

// InfoNetwork defines the Network metric to contain the Results in Info
type InfoNetwork struct {
	Receive  []interface{} `json:"receive"`
	Transmit []interface{} `json:"transmit"`
}

// NewInfoNetwork init the InfoNetwork
func NewInfoNetwork() *InfoNetwork {
	return &InfoNetwork{}
}

// InfoFilesystem defines the Filesystem metric to contain the Results in Info
type InfoFilesystem struct {
	Read  []interface{} `json:"read"`
	Write []interface{} `json:"write"`
	Usage []interface{} `json:"usage_bytes"`
	Limit []interface{} `json:"total_bytes"`
}

// NewInfoFilesystem init the nfoFilesystem
func NewInfoFilesystem() *InfoFilesystem {
	return &InfoFilesystem{}
}
