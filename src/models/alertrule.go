package models

import "time"

type Rule struct {
	ID        uint64    `json:"ID" gorm:"primary_key"`
	CreatedAt time.Time `json:"CreatedAt"`
	UpdatedAt time.Time `json:"UpdatedAt"`
	Name      string    `json:"name" gorm:"not null;unique_index:idx_rule"`
	Alert     string    `json:"alert" gorm:"not null;unique_index:idx_rule"`
	Expr      string    `json:"if" gorm:"size:1020"`
	Duration  string    `json:"for"`
	Labels    string    `json:"labels"`
	Annotations
}

type Annotations struct {
	Description string `json:"description"`
	Summary     string `json:"summary"`
}

type RuleOperation struct {
	Rule *Rule
	File string
	MD5  []byte
}

func NewRuleOperation() *RuleOperation {
	return &RuleOperation{}
}
