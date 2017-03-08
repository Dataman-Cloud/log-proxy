package datastore

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"testing"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/stretchr/testify/assert"
)

func TestInitDBError(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in f", r)
		}
	}()
	_ = InitDB("sqlite3", "")
}

func TestInitDB(t *testing.T) {
	sqliteDBPath := filepath.Join(os.TempDir(), strconv.FormatInt(time.Now().UnixNano(), 10)+"test.db")
	store := InitDB("sqlite3", sqliteDBPath)
	defer os.Remove(sqliteDBPath)
	assert.NotNil(t, store)
}

func TestFrom(t *testing.T) {
	store := From(&gorm.DB{})
	assert.NotNil(t, store)
}
