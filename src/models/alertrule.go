package models

import "github.com/jinzhu/gorm"

type Rule struct {
	gorm.Model
	Name     string `json:"name";gorm:"primary_key"`
	Alert    string `json:"alert";gorm:"primary_key"`
	Expr     string `json:"if"`
	Duration string `json:"for"`
	Labels   string `json:"labels"`
	Annotations
}

type Annotations struct {
	Description string `json:"description"`
	Summary     string `json:"summary"`
}
