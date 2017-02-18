package datastore

import (
	"github.com/Dataman-Cloud/log-proxy/src/models"
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
func InitDB(driver, dsn string) {
	db := database.DB(driver, dsn)

	db.AutoMigrate(&models.Rule{})
	db.AutoMigrate(&models.AlertEvent{})
	// gorm don't support the union unique key, use raw SQL here.
	db.Exec("ALTER TABLE rules ADD UNIQUE KEY(name, alert);")
	db.LogMode(false)
}
