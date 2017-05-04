package models

import "testing"

func TestNewRuleOperation(t *testing.T) {
	if NewRuleOperation() == nil {
		t.Error("faild")
	} else {
		t.Log("success")
	}
}

func TestNewRule(t *testing.T) {
	if NewRule() == nil {
		t.Error("faild")
	} else {
		t.Log("success")
	}
}

func TestNewRulesList(t *testing.T) {
	if NewRulesList() == nil {
		t.Error("faild")
	} else {
		t.Log("success")
	}
}
