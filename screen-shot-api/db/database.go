package db

import (
	"sync"

	"github.com/jinzhu/gorm"
)

var (
	db   *gorm.DB
	once sync.Once
)

// InitDB : initialize database object.
func InitDB(driver, connectionstring string) error {
	db, err := gorm.Open(driver, connectionstring)
	if err != nil {
		return err
	}
	setDB(db)
	return nil
}

// setDB : set database object.
// singleton pattern.
func setDB(newdb *gorm.DB) {
	once.Do(func() {
		db = newdb
	})
}

// DB : get database object.
func DB() *gorm.DB {
	return db
}

// Close : close database.
func Close() error {
	return db.Close()
}
