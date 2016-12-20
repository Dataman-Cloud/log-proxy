package models

// Prometheus prometheus event
type Prometheus struct {
	CLs CommonLabels     `json:"commonLabels"`
	CAs CommonAnnotation `json:"commonAnnotations"`
}

// CommonLabels prometheus common labels
type CommonLabels struct {
	AlertName   string `json:"alertname"`
	AppID       string `json:"container_label_APP_ID"`
	Vcluster    string `json:"container_label_VCLUSTER"`
	ContainerID string `json:"id"`
	Instance    string `json:"instance"`
	Job         string `json:"job"`
	Name        string `json:"name"`
	Severity    string `json:"critical"`
	CreateTime  string `json:"createtime"`
	Condition   string `json:"condition"`
	Usage       string `json:"usage"`
}

// CommonAnnotation prometheus commn annotation
type CommonAnnotation struct {
	Description string `json:"description"`
	Summary     string `json:"summary"`
}
