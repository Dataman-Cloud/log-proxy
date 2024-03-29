package datastore

import (
	"database/sql/driver"
	"errors"
	"testing"

	"github.com/Dataman-Cloud/log-proxy/src/models"

	testdb "github.com/erikstmartin/go-testdb"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
)

func TestCreateLogAlertRule(t *testing.T) {
	db, _ := gorm.Open("testdb", "")
	assert.NotNil(t, db)
	store := &datastore{db}

	testdb.SetQueryWithArgsFunc(func(query string, args []driver.Value) (driver.Rows, error) {
		columns := []string{"id", "cluster", "app", "keyword"}
		rows := ""
		return testdb.RowsFromCSVString(columns, rows), nil
	})

	testdb.SetExecWithArgsFunc(func(query string, args []driver.Value) (result driver.Result, err error) {
		return testResult{1, 1}, nil
	})

	rule := &models.LogAlertRule{Cluster: "c1", App: "a1", Keyword: "k1"}
	err := store.CreateLogAlertRule(rule)
	assert.NoError(t, err)
}

func TestCreateLogAlertRuleExistedError(t *testing.T) {
	db, _ := gorm.Open("testdb", "")
	assert.NotNil(t, db)
	store := &datastore{db}

	testdb.SetQueryWithArgsFunc(func(query string, args []driver.Value) (driver.Rows, error) {
		columns := []string{"id", "cluster", "app", "keyword"}
		rows := "1, c, a, k"
		return testdb.RowsFromCSVString(columns, rows), nil
	})

	rule := &models.LogAlertRule{Cluster: "c1", App: "a1", Keyword: "k1"}
	err := store.CreateLogAlertRule(rule)
	assert.Error(t, err)
}

func TestUpdateLogAlertRuleNotExistedError(t *testing.T) {
	db, _ := gorm.Open("testdb", "")
	assert.NotNil(t, db)
	store := &datastore{db}

	testdb.SetQueryWithArgsFunc(func(query string, args []driver.Value) (driver.Rows, error) {
		columns := []string{"id", "cluster", "app", "keyword"}
		rows := ""
		return testdb.RowsFromCSVString(columns, rows), nil
	})

	rule := &models.LogAlertRule{Cluster: "c1", App: "a1", Keyword: "k1"}
	err := store.UpdateLogAlertRule(rule)
	assert.Error(t, err)
}

func TestUpdateLogAlertRule(t *testing.T) {
	db, _ := gorm.Open("testdb", "")
	assert.NotNil(t, db)
	store := &datastore{db}

	testdb.SetQueryWithArgsFunc(func(query string, args []driver.Value) (driver.Rows, error) {
		columns := []string{"id", "cluster", "app", "keyword"}
		rows := "1, c, a, k"
		return testdb.RowsFromCSVString(columns, rows), nil
	})

	testdb.SetExecWithArgsFunc(func(query string, args []driver.Value) (result driver.Result, err error) {
		return testResult{1, 1}, nil
	})

	rule := &models.LogAlertRule{Cluster: "c1", App: "a1", Keyword: "k1"}
	err := store.UpdateLogAlertRule(rule)
	assert.NoError(t, err)
}

func TestDeleteDeleteLogAlertRule(t *testing.T) {
	db, _ := gorm.Open("testdb", "")
	assert.NotNil(t, db)
	store := &datastore{db}

	testdb.SetExecWithArgsFunc(func(query string, args []driver.Value) (result driver.Result, err error) {
		return testResult{1, 1}, nil
	})

	err := store.DeleteLogAlertRule("test")
	assert.NoError(t, err)
}

func TestGetLogAlertRule(t *testing.T) {
	db, _ := gorm.Open("testdb", "")
	assert.NotNil(t, db)
	store := &datastore{db}

	testdb.SetQueryWithArgsFunc(func(query string, args []driver.Value) (driver.Rows, error) {
		columns := []string{"id", "cluster", "app", "keyword"}
		rows := "1, c, a, k"
		return testdb.RowsFromCSVString(columns, rows), nil
	})

	rule, err := store.GetLogAlertRule("test")
	assert.NoError(t, err)
	assert.Equal(t, rule.ID, uint64(1))
}

func TestGetLogAlertRules(t *testing.T) {
	db, _ := gorm.Open("testdb", "")
	assert.NotNil(t, db)
	store := &datastore{db}

	testdb.SetQueryWithArgsFunc(func(query string, args []driver.Value) (driver.Rows, error) {
		columns := []string{"count"}
		rows := "1"
		return testdb.RowsFromCSVString(columns, rows), nil
	})

	_, err := store.GetLogAlertRules(map[string]interface{}{"test": "test"}, models.Page{})
	assert.NoError(t, err)

	groups := []string{"work"}
	_, err = store.GetLogAlertRules(map[string]interface{}{"test": "test", "groups": groups}, models.Page{})
	assert.NoError(t, err)
}

func TestCreateLogAlertEvent(t *testing.T) {
	db, _ := gorm.Open("testdb", "")
	assert.NotNil(t, db)
	store := &datastore{db}

	testdb.SetExecWithArgsFunc(func(query string, args []driver.Value) (result driver.Result, err error) {
		return testResult{1, 1}, nil
	})

	err := store.CreateLogAlertEvent(&models.LogAlertEvent{})
	assert.Nil(t, err)
}

func TestGetLogAlertEvents(t *testing.T) {
	db, _ := gorm.Open("testdb", "")
	assert.NotNil(t, db)
	store := &datastore{db}

	testdb.SetQueryWithArgsFunc(func(query string, args []driver.Value) (driver.Rows, error) {
		columns := []string{"count"}
		rows := "1"
		return testdb.RowsFromCSVString(columns, rows), nil
	})

	_, err := store.GetLogAlertEvents(map[string]interface{}{"test": "test"}, models.Page{RangeFrom: 0, RangeTo: 10})
	assert.NoError(t, err)

	testdb.SetQueryWithArgsFunc(func(query string, args []driver.Value) (driver.Rows, error) {
		columns := []string{"count"}
		rows := "1"
		return testdb.RowsFromCSVString(columns, rows), errors.New("test")
	})
	_, err = store.GetLogAlertEvents(map[string]interface{}{"test": "test"}, models.Page{RangeFrom: 0, RangeTo: 10})
	assert.NotNil(t, err)

	groups := []string{"work"}
	_, err = store.GetLogAlertEvents(map[string]interface{}{"test": "test", "groups": groups}, models.Page{RangeFrom: 0, RangeTo: 10})
	assert.NotNil(t, err)
}

func TestGetLogAlertEvent(t *testing.T) {
	db, _ := gorm.Open("testdb", "")
	assert.NotNil(t, db)
	store := &datastore{db}

	testdb.SetQueryWithArgsFunc(func(query string, args []driver.Value) (driver.Rows, error) {
		columns := []string{"count"}
		rows := "1"
		return testdb.RowsFromCSVString(columns, rows), nil
	})

	_, err := store.GetLogAlertEvent("test")
	assert.NoError(t, err)
}

func TestAckLogAlertEvent(t *testing.T) {
	db, _ := gorm.Open("testdb", "")
	assert.NotNil(t, db)
	store := &datastore{db}

	testdb.SetExecWithArgsFunc(func(query string, args []driver.Value) (result driver.Result, err error) {
		return testResult{1, 1}, nil
	})

	err := store.AckLogAlertEvent("test")
	assert.Nil(t, err)
}

func TestGetLogAlertApps(t *testing.T) {
	db, _ := gorm.Open("testdb", "")
	assert.NotNil(t, db)
	store := &datastore{db}

	testdb.SetQueryWithArgsFunc(func(query string, args []driver.Value) (driver.Rows, error) {
		columns := []string{"count"}
		rows := "1"
		return testdb.RowsFromCSVString(columns, rows), nil
	})

	_, err := store.GetLogAlertApps(map[string]interface{}{"test": "test"}, models.Page{RangeFrom: 0, RangeTo: 10})
	assert.NoError(t, err)

	testdb.SetQueryWithArgsFunc(func(query string, args []driver.Value) (driver.Rows, error) {
		columns := []string{"count"}
		rows := "1"
		return testdb.RowsFromCSVString(columns, rows), errors.New("test")
	})
	_, err = store.GetLogAlertApps(map[string]interface{}{"test": "test"}, models.Page{RangeFrom: 0, RangeTo: 10})
	assert.NotNil(t, err)

	groups := []string{"work"}
	_, err = store.GetLogAlertApps(map[string]interface{}{"test": "test", "groups": groups}, models.Page{RangeFrom: 0, RangeTo: 10})
	assert.NotNil(t, err)

}
