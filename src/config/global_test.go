package config

import (
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestLoadLogOptionalLabels(t *testing.T) {
	LoadLogOptionalLabels()
	assert.Equal(t, LogOptionaLabels, DefaultLogOptionalLabels)
	labelsConfigPath := c.LabelsConfigPath
	defer func() {
		c.LabelsConfigPath = labelsConfigPath
	}()

	c.LabelsConfigPath = ""
	LoadLogOptionalLabels()

	viper.AddConfigPath("../../")
	c.LabelsConfigPath = "mola-conf.template"
	LoadLogOptionalLabels()
}
