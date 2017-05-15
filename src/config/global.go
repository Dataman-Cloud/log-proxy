package config

import (
	"errors"
	"net/url"

	"github.com/Sirupsen/logrus"
	"github.com/spf13/viper"
)

var logOptionaLabels map[string]string

var DefaultLogOptionalLabels = map[string]string{
	"app":     "DM_APP_ID",
	"cluster": "DM_VCLUSTER",
	"slot":    "DM_SLOT_INDEX",
	"task":    "DM_TASK_ID",
	"source":  "path",
	"offset":  "offset",
	"user":    "DM_USER",
	"group":   "DM_GROUP_NAME",
	"message": "message",
	"tag":     "DM_LOG_TAG",
}

func LoadLogOptionalLabels() {
	logOptionaLabels = DefaultLogOptionalLabels

	if c == nil || c.LabelsConfigPath == "" {
		return
	}

	logrus.Info("Init log optional labels by config file....")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.SetConfigName(c.LabelsConfigPath)

	if err := viper.ReadInConfig(); err != nil {
		logrus.Errorf("read log optional from config file failed. Error: %s \n", err.Error())
		return
	}

	logOptionaLabels = viper.GetStringMapString("log-query-options")

	logrus.Infof("load log optional labels fron %s success \n", viper.ConfigFileUsed())

	logrus.Infof("the logOptionaLabels used is %+v \n", logOptionaLabels)
	return
}

func ConvertRequestQueryParams(values url.Values) map[string]interface{} {
	result := make(map[string]interface{})
	for k, v := range values {
		if label, ok := logOptionaLabels[k]; ok && len(v) > 0 {
			result[label] = v[0]
		}
	}

	return result
}

func LogOffsetLabel() string {
	offset, ok := logOptionaLabels["offset"]
	if ok {
		return offset
	}

	logrus.Debug("log offset label not found. use default value: offset")
	return "offset"
}

func LogAppLabel() string {
	app, ok := logOptionaLabels["app"]
	if ok {
		return app
	}

	logrus.Debug("log app label not found. use default value: DM_APP_ID")
	return "DM_APP_ID"
}

func LogKeywordLabel() string {
	return "keyword"
}

func LogConjLabel() string {
	return "conj"
}

func LogMessageLabel() string {
	message, ok := logOptionaLabels["message"]
	if ok {
		return message
	}

	logrus.Debug("log message label not found. use default value: message")
	return "message"
}

func LogSourceLabel() string {
	source, ok := logOptionaLabels["source"]
	if ok {
		return source
	}

	logrus.Debug("log source label not found. use default value: path")
	return "path"
}

func LogTaskLabel() string {
	task, ok := logOptionaLabels["task"]
	if ok {
		return task
	}

	logrus.Debug("log task label not found. use default valeu: DM_TASK_ID")
	return "DM_TASK_ID"
}

func LogSlotLabel() string {
	slot, ok := logOptionaLabels["slot"]
	if ok {
		return slot
	}

	logrus.Debug("log slot label not found. use default value: DM_SLOT_INDEX")
	return "DM_SLOT_INDEX"
}

func GetLogLabel(key string) (string, error) {
	v, ok := logOptionaLabels[key]
	if ok {
		return v, nil
	}

	return "", errors.New("invalid log label: " + key)
}

const lablePrefix = "container_label_"

func MonitorAppLabel() string {
	app, ok := logOptionaLabels["app"]
	if ok {
		return lablePrefix + app
	}

	logrus.Debug("App app label not found. use default value: DM_APP_ID")
	return lablePrefix + "DM_APP_ID"
}

func MonitorSlotLabel() string {
	slot, ok := logOptionaLabels["slot"]
	if ok {
		return lablePrefix + slot
	}

	logrus.Debug("log slot label not found. use default value: DM_SLOT_INDEX")
	return lablePrefix + "DM_SLOT_INDEX"
}

func LogTagLabel() string {
	tag, ok := logOptionaLabels["tag"]
	if ok {
		return lablePrefix + tag
	}

	logrus.Debug("log slot label not found. use default value: DM_LOG_TAG")
	return lablePrefix + "DM_LOG_TAG"
}
