package models

type Prometheus struct {
	CLs CommonLabels     `json:"commonLabels"`
	CAs CommonAnnotation `json:"commonAnnotations"`
}

type CommonLabels struct {
	AlertName   string `json:"alertname"`
	AppId       string `json:"container_label_APP_ID"`
	Vcluster    string `json:"container_label_VCLUSTER"`
	ContainerId string `json:"id"`
	Instance    string `json:"instance"`
	Job         string `json:"job"`
	Name        string `json:"name"`
	Severity    string `json:"critical"`
	CreateTime  string `json:"createtime"`
	Condition   string `json:"condition"`
	Usage       string `json:"usage"`
}

type CommonAnnotation struct {
	Description string `json:"description"`
	Summary     string `json:"summary"`
}
