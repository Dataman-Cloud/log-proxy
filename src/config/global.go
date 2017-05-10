package config

import (
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

	logrus.Debug("log offset not found use default value: offset")
	return "offset"
}

func LogAppLabel() string {
	app, ok := logOptionaLabels["app"]
	if ok {
		return app
	}

	logrus.Debug("log app not found use default value: DM_APP_ID")
	return "DM_APP_ID"
}
