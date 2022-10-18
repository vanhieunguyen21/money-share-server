package database

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"money_share/pkg/model"
)

const (
	host = "localhost"
	user = "admin"
	password = "123456"
	dbname = "money_share"
)

type PostgresDB struct {
	DB *gorm.DB
}

var Postgres = &PostgresDB{}

func Connect() *PostgresDB {
	dsn := fmt.Sprintf("host=%s user=%s dbname=%s password=%s sslmode=disable", host, user, dbname, password)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	err = db.AutoMigrate(&model.User{}, &model.Group{}, &model.Member{}, &model.Expense{})
	if err != nil {
		panic(err)
	}

	Postgres.DB = db
	return Postgres
}
