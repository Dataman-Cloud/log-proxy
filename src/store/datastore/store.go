package datastore

import (
	"github.com/Dataman-Cloud/log-proxy/src/models"
	"github.com/Dataman-Cloud/log-proxy/src/store"
	"github.com/Dataman-Cloud/log-proxy/src/utils/database"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type datastore struct {
	*gorm.DB
}

// InitDB init the table struct in DB
//  - driver   mysql|sqlite|pgsql
//  - dsn      user:password@/dbname?charset=utf8&parseTime=True&loc=Local
func InitDB(driver, dsn string) store.Store {
	db := database.DB(driver, dsn)

	db.AutoMigrate(&models.Rule{})
	db.AutoMigrate(&models.Event{})
	db.AutoMigrate(&models.LogAlertRule{})
	db.AutoMigrate(&models.LogAlertEvent{})

	return From(db)
}

func From(db *gorm.DB) store.Store {
	return &datastore{db}
}
