package models

import (
	"testing"
)

func TestNewInfoNetwork(t *testing.T) {
	if NewInfoNetwork() == nil {
		t.Error("faild")
	} else {
		t.Log("success")
	}
}

func TestNewInfoFilesystem(t *testing.T) {
	if NewInfoFilesystem() == nil {
		t.Error("faild")
	} else {
		t.Log("success")
	}
}
