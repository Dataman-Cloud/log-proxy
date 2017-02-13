package database

import (
	log "github.com/Sirupsen/logrus"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/mattes/migrate/driver/mysql"
)

var db *gorm.DB

// DB returns the struct gorm.DB
func DB(driver, dsn string) *gorm.DB {
	if db == nil {
		OpenDB(driver, dsn)
	}
	return db
}

func GetDB() *gorm.DB {
	return db
}

// OpenDB initilize the db connection
func OpenDB(driver, dsn string) {
	var err error
	log.Infof("connecting mysql uri: %s", dsn)

	db, err = gorm.Open(driver, dsn)
	if err != nil {
		log.Fatalf("init mysql error: %v", err)
		panic("database connection failed")
	}

	db.DB().SetMaxIdleConns(1)
	db.DB().SetMaxOpenConns(2)
	db.SetLogger(log.StandardLogger())

	db.LogMode(false)
}
