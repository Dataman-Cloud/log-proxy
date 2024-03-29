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

	tmp := offset
	defer func() {
		logOptionaLabels["offset"] = tmp
	}()

	delete(logOptionaLabels, "offset")
	offset = LogOffsetLabel()
	assert.Equal(t, offset, "offset")
}

func TestLogAppLabel(t *testing.T) {
	app := LogAppLabel()
	assert.Equal(t, app, "DM_APP_ID")

	tmp := app
	defer func() {
		logOptionaLabels["app"] = tmp
	}()

	delete(logOptionaLabels, "app")
	app = LogAppLabel()
	assert.Equal(t, app, "DM_APP_ID")
}

func TestLogKeywordLabel(t *testing.T) {
	assert.Equal(t, LogKeywordLabel(), "keyword")
}

func TestLogConjLabel(t *testing.T) {
	assert.Equal(t, LogConjLabel(), "conj")
}

func TestLogMessageLabel(t *testing.T) {
	message := LogMessageLabel()
	assert.Equal(t, message, "message")

	tmp := message
	defer func() {
		logOptionaLabels["message"] = tmp
	}()

	delete(logOptionaLabels, "message")
	message = LogMessageLabel()
	assert.Equal(t, message, "message")
}

func TestLogSourceLabel(t *testing.T) {
	source := LogSourceLabel()
	assert.Equal(t, source, "path")

	tmp := source
	defer func() {
		logOptionaLabels["source"] = tmp
	}()

	delete(logOptionaLabels, "source")
	source = LogSourceLabel()
	assert.Equal(t, source, "path")
}

func TestLogSlotLabel(t *testing.T) {
	slot := LogSlotLabel()
	assert.Equal(t, slot, "DM_SLOT_INDEX")

	tmp := slot
	defer func() {
		logOptionaLabels["slot"] = tmp
	}()

	delete(logOptionaLabels, "slot")
	slot = LogSlotLabel()
	assert.Equal(t, slot, "DM_SLOT_INDEX")
}

func TestLogTaskLabel(t *testing.T) {
	task := LogTaskLabel()
	assert.Equal(t, task, "DM_TASK_ID")

	tmp := task
	defer func() {
		logOptionaLabels["task"] = tmp
	}()

	delete(logOptionaLabels, "task")
	task = LogTaskLabel()
	assert.Equal(t, task, "DM_TASK_ID")
}

func TestGetLogLabel(t *testing.T) {
	_, err := GetLogLabel("app")
	assert.NoError(t, err)

	_, err = GetLogLabel("xxxxxxxxxxxxx")
	assert.Error(t, err)
}

func TestMonitorAppLabel(t *testing.T) {
	app := MonitorAppLabel()
	assert.Equal(t, app, "container_label_DM_APP_ID")

	tmp := app
	defer func() {
		logOptionaLabels["app"] = tmp
	}()

	delete(logOptionaLabels, "app")
	app = MonitorAppLabel()
	assert.Equal(t, app, "container_label_DM_APP_ID")
}

func TestMonitorSlotLabel(t *testing.T) {
	slot := MonitorSlotLabel()
	assert.Equal(t, slot, "container_label_DM_SLOT_INDEX")

	tmp := slot
	defer func() {
		logOptionaLabels["slot"] = tmp
	}()

	delete(logOptionaLabels, "slot")
	slot = MonitorSlotLabel()
	assert.Equal(t, slot, "container_label_DM_SLOT_INDEX")
}

func TestLogTagLabel(t *testing.T) {
	tag := LogTagLabel()
	assert.Equal(t, tag, "container_label_DM_LOG_TAG")

	tmp := tag
	defer func() {
		logOptionaLabels["tag"] = tmp
	}()

	delete(logOptionaLabels, "tag")
	tag = LogTagLabel()
	assert.Equal(t, tag, "container_label_DM_LOG_TAG")
}
