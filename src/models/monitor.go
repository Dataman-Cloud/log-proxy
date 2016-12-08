package models

type QueryRangeResult struct {
	Status    string `json:"status"`
	Data      *Data  `json:"data"`
	ErrorType string `json:"errorType"`
	Error     string `json:"error"`
}

type Data struct {
	ResultType string    `json:"resultType"`
	Result     []*Result `json:"result"`
}

type Result struct {
	Metric *Metric         `json:"metric"`
	Values [][]interface{} `json:"values"`
}

type Metric struct {
	ContainerLabelAPPID    string `json:"container_label_APP_ID"`
	ContainerLabelVCLUSTER string `json:"container_label_VCLUSTER"`
	ID                     string `json:"id"`
	Image                  string `json:"image"`
	Instance               string `json:"instance"`
	Job                    string `json:"job"`
	Name                   string `json:"name"`
	Group                  string `json:"group,omitempty"`
	Interface              string `json:"interface,omitempty"`
	Device                 string `json:"device,omitempty"`
}

type MetricCPU struct {
	Usage []*Result `json:"usage"`
}

type MetricMemory struct {
	Percentage []*Result `json:"usage"`
	Usage      []*Result `json:"usage_bytes"`
	Total      []*Result `json:"total_bytes"`
}

type MetricNewtork struct {
	Receive  []*Result `json:"receive"`
	Transmit []*Result `json:"transmit"`
}

type MetricFilesystem struct {
	Read  []*Result `json:"read"`
	Write []*Result `json:"write"`
}

type InfoCPU struct {
	Usage []interface{} `json:"usage"`
}

type InfoMemory struct {
	Percetange []interface{} `json:"usage"`
	Usage      []interface{} `json:"usage_bytes"`
	Total      []interface{} `json:"total_bytes"`
}

type InfoNetwork struct {
	Receive  []interface{} `json:"receive"`
	Transmit []interface{} `json:"transmit"`
}

func NewInfoNetwork() *InfoNetwork {
	return &InfoNetwork{}
}

type InfoFilesystem struct {
	Read  []interface{} `json:"read"`
	Write []interface{} `json:"write"`
	Usage []interface{} `json:"usage_bytes"`
	Limit []interface{} `json:"total_bytes"`
}

func NewInfoFilesystem() *InfoFilesystem {
	return &InfoFilesystem{}
}
