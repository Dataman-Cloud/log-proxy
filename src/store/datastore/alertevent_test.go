package datastore

import (
	"database/sql/driver"
	"testing"

	"github.com/Dataman-Cloud/log-proxy/src/models"

	"github.com/erikstmartin/go-testdb"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
)

func TestCreateOrIncreaseEvent(t *testing.T) {
	db, _ := gorm.Open("testdb", "")
	store := &datastore{db}

	testdb.SetQueryWithArgsFunc(func(query string, args []driver.Value) (driver.Rows, error) {
		columns := []string{"ack", "severity", "v_cluster", "app", "slot", "container_id", "alert_name"}
		rows := ""
		if args[0] != nil {
			rows = "0, Warning, wtzhou-VCluster, app-6, slot-5, /docker/aaaxefgh32e2e23rfsda, alert"
		}
		return testdb.RowsFromCSVString(columns, rows), nil
	})

	testdb.SetExecWithArgsFunc(func(query string, args []driver.Value) (result driver.Result, err error) {
		return testResult{1, 1}, nil
	})

	event := &models.Event{
		Ack:         false,
		Severity:    "Warning",
		VCluster:    "wtzhou-VCluster",
		App:         "app-6",
		Slot:        "slot-5",
		ContainerID: "/docker/aaaxefgh32e2e23rfsda",
		AlertName:   "alert",
		Description: "description",
	}

	err := store.CreateOrIncreaseEvent(event)
	assert.Nil(t, err)
}

func TestAckEvent(t *testing.T) {
	ID := 1
	userName := "user1"
	groupName := "group1"

	db, _ := gorm.Open("testdb", "")
	store := &datastore{db}
	testdb.SetQueryWithArgsFunc(func(query string, args []driver.Value) (driver.Rows, error) {
		columns := []string{"id", "user_name", "group_name", "ack"}
		rows := ""
		if args[1] == "user1" {
			rows = "1, user1, group1, 0"
		}
		return testdb.RowsFromCSVString(columns, rows), nil
	})

	testdb.SetExecWithArgsFunc(func(query string, args []driver.Value) (result driver.Result, err error) {
		return testResult{1, 1}, nil
	})

	err := store.AckEvent(ID, userName, groupName)
	assert.Nil(t, err)

	err = store.AckEvent(ID, "user2", groupName)
	assert.NotNil(t, err)
}

func TestListAckedEvent(t *testing.T) {
	page := models.Page{}
	userName := "user1"
	groupName := ""

	db, _ := gorm.Open("testdb", "")
	store := &datastore{db}

	testdb.SetQueryWithArgsFunc(func(query string, args []driver.Value) (driver.Rows, error) {
		columns := []string{"ack", "user_name", "group_name"}
		rows := ""
		if args[1] == "user1" {
			rows = "1, user1, group1"
		}
		return testdb.RowsFromCSVString(columns, rows), nil
	})
	data := store.ListAckedEvent(page, userName, groupName)
	assert.NotNil(t, data)
}

func TestListUnackedEvent(t *testing.T) {
	page := models.Page{}
	userName := "user1"
	groupName := ""

	db, _ := gorm.Open("testdb", "")
	store := &datastore{db}

	testdb.SetQueryWithArgsFunc(func(query string, args []driver.Value) (driver.Rows, error) {
		columns := []string{"ack", "user_name", "group_name"}
		rows := ""
		if args[1] == "user1" {
			rows = "0, user1, group1"
		}
		return testdb.RowsFromCSVString(columns, rows), nil
	})
	data := store.ListUnackedEvent(page, userName, groupName)
	assert.NotNil(t, data)
}
