package database

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"money_share/pkg/model"
)

const (
	host     = "localhost"
	user     = "admin"
	password = "123456"
	dbname   = "money_share_test"
)

func Connect() (db *gorm.DB, err error) {
	dsn := fmt.Sprintf("host=%s user=%s dbname=%s password=%s sslmode=disable", host, user, dbname, password)

	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return
	}

	// Drop old database schema
	err = db.Migrator().DropTable(&model.User{}, &model.Group{}, &model.Member{}, &model.Expense{})
	if err != nil {
		return
	}

	// Auto migrate new database schema
	err = db.AutoMigrate(&model.User{}, &model.Group{}, &model.Member{}, &model.Expense{})
	if err != nil {
		return
	}

	return
}
