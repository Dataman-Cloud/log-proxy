package models

type QueryRangeResult struct {
	Status string `json:"status"`
	Data   *Data  `json:"data"`
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
	ContainerLabelAPPID string `json:"container_label_APP_ID"`
	Group               string `json:"group"`
	ID                  string `json:"id"`
	Image               string `json:"image"`
	Instance            string `json:"instance"`
	Job                 string `json:"job"`
	Name                string `json:"name"`
}

type MetricList struct {
	CPU    []*Result `json:"cpu"`
	Memory []*Result `json:"memory"`
}
