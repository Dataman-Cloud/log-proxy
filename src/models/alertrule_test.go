package models

import "testing"

func TestNewRuleOperation(t *testing.T) {
	if NewRuleOperation() == nil {
		t.Error("faild")
	} else {
		t.Log("success")
	}
}
