package models

type Prometheus struct {
	CLs CommonLabels `json:"commonLabels"`
}

type CommonLabels struct {
	AlertName   string `json:"alertname"`
	AppId       string `json:"container_label_APP_ID"`
	Vcluster    string `json:"container_label_VCLUSTER"`
	ContainerId string `json:"id"`
	Instance    string `json:"instance"`
	Job         string `json:"job"`
	Name        string `json:"name"`
}
