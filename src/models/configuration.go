package models

import "github.com/jinzhu/gorm"

type Configuration struct {
	gorm.Model
	CamaCmdbDefaultAppID string `json:"cama_cmdb_default_app_id"`
	CamaNotifactionADDR  string `json:"cama_notifaction_addr"`
	CamaNotifactionTempl string `json:"cama_notifaction_templ"`
}

func NewConfiguration() *Configuration {
	return &Configuration{
		CamaCmdbDefaultAppID: "cama_cmdb_id",
		CamaNotifactionADDR:  "127.0.0.1:8030",
		CamaNotifactionTempl: "cama notification template",
	}
}
