package config

import (
	"github.com/Sirupsen/logrus"
	"github.com/spf13/viper"
)

var LogOptionaLabels map[string]string

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
	LogOptionaLabels = DefaultLogOptionalLabels

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

	LogOptionaLabels = viper.GetStringMapString("log-query-options")

	logrus.Infof("load log optional labels fron %s success \n", viper.ConfigFileUsed())
	return
}
