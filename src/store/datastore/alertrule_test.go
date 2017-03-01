package datastore

import (
	"database/sql/driver"
	"errors"
	"testing"

	"github.com/Dataman-Cloud/log-proxy/src/models"

	"github.com/erikstmartin/go-testdb"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
)

type testResult struct {
	lastID       int64
	affectedRows int64
}

func (r testResult) LastInsertId() (int64, error) {
	return r.lastID, nil
}

func (r testResult) RowsAffected() (int64, error) {
	return r.affectedRows, nil
}

func TestListAlertRules(t *testing.T) {
	db, _ := gorm.Open("testdb", "")
	store := &datastore{db}

	testdb.SetQueryFunc(func(query string) (driver.Rows, error) {
		columns := []string{"id", "name", "alert"}
		result := `
    1, user1, alert
    `
		return testdb.RowsFromCSVString(columns, result), nil
	})

	page := models.Page{}
	result, err := store.ListAlertRules(page, "user1")
	assert.Nil(t, err)
	assert.NotNil(t, result)

	result, err = store.ListAlertRules(page, "")
	assert.Nil(t, err)
	assert.NotNil(t, result)

	testdb.SetQueryFunc(func(query string) (driver.Rows, error) {
		return nil, errors.New("db error")
	})
	_, err = store.ListAlertRules(page, "user1")
	assert.NotNil(t, err)

	_, err = store.ListAlertRules(page, "")
	assert.NotNil(t, err)
}

func TestGetAlertRule(t *testing.T) {
	db, _ := gorm.Open("testdb", "")
	store := &datastore{db}

	testdb.SetQueryFunc(func(query string) (driver.Rows, error) {
		columns := []string{"id", "name", "alert"}
		result := `
    1, user1, alert
    `
		return testdb.RowsFromCSVString(columns, result), nil
	})
	result, err := store.GetAlertRule(1)

	assert.Equal(t, uint64(1), result.ID, "Get the id 1 rule")
	assert.Nil(t, err)

	testdb.SetQueryFunc(func(query string) (driver.Rows, error) {
		return nil, errors.New("db error")
	})
	_, err = store.GetAlertRule(1)
	assert.NotNil(t, err)
}

func TestGetAlertRuleByName(t *testing.T) {
	db, _ := gorm.Open("testdb", "")
	store := &datastore{db}

	testdb.SetQueryFunc(func(query string) (driver.Rows, error) {
		columns := []string{"id", "name", "alert"}
		result := `
    1, user1, alert
    `
		return testdb.RowsFromCSVString(columns, result), nil
	})
	result, err := store.GetAlertRuleByName("user1", "alert")
	assert.Equal(t, "user1", result.Name, "Get the rule name is user1")
	assert.Nil(t, err)

	testdb.SetQueryFunc(func(query string) (driver.Rows, error) {
		return nil, errors.New("db error")
	})
	_, err = store.GetAlertRuleByName("user1", "alert")
	assert.NotNil(t, err)
}

func TestCreateAlertRule(t *testing.T) {
	db, _ := gorm.Open("testdb", "")
	store := &datastore{db}

	testdb.SetQueryWithArgsFunc(func(query string, args []driver.Value) (driver.Rows, error) {
		columns := []string{"id", "name", "alert"}
		rows := ""
		if args[0] == "user1" && args[1] == "alert" {
			rows = "1, user1, alert"
		}
		return testdb.RowsFromCSVString(columns, rows), nil
	})

	testdb.SetExecWithArgsFunc(func(query string, args []driver.Value) (result driver.Result, err error) {
		return testResult{1, 1}, nil
	})

	rule := &models.Rule{
		Name:  "user2",
		Alert: "alert",
	}
	err := store.CreateAlertRule(rule)
	assert.Nil(t, err)

	rule.Name = "user1"
	err = store.CreateAlertRule(rule)
	assert.NotNil(t, err)
}

func TestUpdateAlertRule(t *testing.T) {
	db, _ := gorm.Open("testdb", "")
	store := &datastore{db}

	testdb.SetQueryWithArgsFunc(func(query string, args []driver.Value) (driver.Rows, error) {
		columns := []string{"id", "name", "alert"}
		rows := ""
		if args[0] == "user1" && args[1] == "alert" {
			rows = "1, user1, alert"
		}
		return testdb.RowsFromCSVString(columns, rows), nil
	})

	testdb.SetExecWithArgsFunc(func(query string, args []driver.Value) (result driver.Result, err error) {
		return testResult{1, 1}, nil
	})

	rule := &models.Rule{
		Name:  "user2",
		Alert: "alert",
	}
	err := store.UpdateAlertRule(rule)
	assert.NotNil(t, err)

	rule.Name = "user1"
	err = store.UpdateAlertRule(rule)
	assert.Nil(t, err)
}

func TestDeleteAlertRuleByIDName(t *testing.T) {
	db, _ := gorm.Open("testdb", "")
	store := &datastore{db}

	testdb.SetExecWithArgsFunc(func(query string, args []driver.Value) (result driver.Result, err error) {
		return testResult{1, 1}, nil
	})

	_, err := store.DeleteAlertRuleByIDName(uint64(1), "user1")
	assert.Nil(t, err)

	testdb.SetExecWithArgsFunc(func(query string, args []driver.Value) (result driver.Result, err error) {
		return nil, errors.New("db error")
	})

	_, err = store.DeleteAlertRuleByIDName(uint64(1), "user2")
	assert.NotNil(t, err)

}
