package models

import "github.com/jinzhu/gorm"

type Rule struct {
	gorm.Model
	Name    string `json:"name";gorm:"primary_key"`
	Alert   string `json:"alert";gorm:"primary_key"`
	ForTime string `json:"for_time"`
	Labels  string `json:"labels"`
	Annotations
}

type Annotations struct {
	Description string `json:"description"`
	Summary     string `json:"summary"`
}

type Rules struct {
	Rules []*Rule `json:"rules"`
}
