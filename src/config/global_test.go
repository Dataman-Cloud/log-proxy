package config

import (
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestLoadLogOptionalLabels(t *testing.T) {
	LoadLogOptionalLabels()
	assert.Equal(t, logOptionaLabels, DefaultLogOptionalLabels)
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

func TestConvertRequestQueryParams(t *testing.T) {
	testValues := map[string][]string{"app": []string{"test"}, "source": []string{"stdout"}}
	result := ConvertRequestQueryParams(testValues)
	assert.Equal(t, "test", result["DM_APP_ID"].(string))
}

func TestLogOffsetLabel(t *testing.T) {
	offset := LogOffsetLabel()
	assert.Equal(t, offset, "offset")

	tmp := logOptionaLabels
	defer func() {
		logOptionaLabels = tmp
	}()

	delete(logOptionaLabels, "offset")
	offset = LogOffsetLabel()
	assert.Equal(t, offset, "offset")
}

func TestLogAppLabel(t *testing.T) {
	app := LogAppLabel()
	assert.Equal(t, app, "DM_APP_ID")

	tmp := logOptionaLabels
	defer func() {
		logOptionaLabels = tmp
	}()

	delete(logOptionaLabels, "app")
	app = LogAppLabel()
	assert.Equal(t, app, "DM_APP_ID")
}
