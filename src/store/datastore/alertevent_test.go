package datastore

import (
	"database/sql/driver"
	"testing"

	"github.com/Dataman-Cloud/log-proxy/src/models"

	"github.com/erikstmartin/go-testdb"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
)

func TestCreateEvent(t *testing.T) {
	var rows string
	db, _ := gorm.Open("testdb", "")
	store := &datastore{db}

	testdb.SetQueryWithArgsFunc(func(query string, args []driver.Value) (driver.Rows, error) {
		columns := []string{"alertname", "groupname", "app", "task", "ack", "severity"}
		if args[0] != nil {
			rows = "web_zdou_datamanmesos_cpu_usage_warning, dev, web-zdou-datamanmesos, 0, 0, warning"
		}
		return testdb.RowsFromCSVString(columns, rows), nil
	})

	testdb.SetExecWithArgsFunc(func(query string, args []driver.Value) (result driver.Result, err error) {
		return testResult{1, 1}, nil
	})

	event := &models.Event{
		AlertName: "web_zdou_datamanmesos_cpu_usage_warning",
		Ack:       false,
		Group:     "dev",
		App:       "web-zdou-datamanmesos",
		Task:      "0",
		Severity:  "warning",
	}

	err := store.CreateOrIncreaseEvent(event)
	assert.Nil(t, err)
}

func TestIncreaseEvent(t *testing.T) {
	var rows string
	db, _ := gorm.Open("testdb", "")
	store := &datastore{db}

	testdb.SetQueryWithArgsFunc(func(query string, args []driver.Value) (driver.Rows, error) {
		columns := []string{"alertname", "groupname", "app", "task", "ack", "severity"}
		if args[0] != nil {
			rows = ""
		}
		return testdb.RowsFromCSVString(columns, rows), nil
	})

	testdb.SetExecWithArgsFunc(func(query string, args []driver.Value) (result driver.Result, err error) {
		return testResult{1, 1}, nil
	})

	event := &models.Event{
		AlertName: "web_zdou_datamanmesos_cpu_usage_warning",
		Group:     "dev",
		App:       "web-zdou-datamanmesos",
		Task:      "0",
		Severity:  "warning",
	}

	err := store.CreateOrIncreaseEvent(event)
	assert.Nil(t, err)
}

func TestAckAlertEvent(t *testing.T) {
	ID := 1
	groupName := "dev"
	app := "web-zdou-datamanmesos"

	db, _ := gorm.Open("testdb", "")
	store := &datastore{db}
	testdb.SetQueryWithArgsFunc(func(query string, args []driver.Value) (driver.Rows, error) {
		columns := []string{"id", "groupname", "app", "ack"}
		rows := ""
		if args[1] == "dev" {
			rows = "1, dev, web-zdou-datamanmesos, 0"
		}
		return testdb.RowsFromCSVString(columns, rows), nil
	})

	testdb.SetExecWithArgsFunc(func(query string, args []driver.Value) (result driver.Result, err error) {
		return testResult{1, 1}, nil
	})

	err := store.AckEvent(ID, groupName, app)
	assert.Nil(t, err)

	err = store.AckEvent(ID, "test", app)
	assert.NotNil(t, err)
}

func TestListEvents(t *testing.T) {
	var rows string
	page := models.Page{}
	page.PageFrom = 0
	page.PageSize = 100

	options := make(map[string]interface{})
	options["group"] = "dev"
	options["app"] = "web-zdou-datamanmesos"
	options["ack"] = false
	options["start"] = "1494139284"
	options["end"] = "1494139315"

	db, _ := gorm.Open("testdb", "")
	store := &datastore{db}

	testdb.SetQueryWithArgsFunc(func(query string, args []driver.Value) (driver.Rows, error) {
		columns := []string{"ack", "groupname", "app"}
		if args[0] != nil {
			rows = "0, dev, web-zdou-datamanmesos"
		}

		return testdb.RowsFromCSVString(columns, rows), nil
	})
	data, _ := store.ListEvents(page, options)
	assert.Nil(t, data)
}
