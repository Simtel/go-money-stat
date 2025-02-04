package app

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"money-stat/internal/model"
)

type DB struct {
	db *gorm.DB
}

func NewDB() *DB {
	db, err := gorm.Open(sqlite.Open("zenmoney.db?cache=shared&mode=rwc"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	errMigrate := db.AutoMigrate(&model.Transaction{}, &model.Tag{}, &model.Instrument{}, &model.Account{})
	if errMigrate != nil {
		panic(errMigrate)
	}
	return &DB{db: db}
}

func (db *DB) GetGorm() *gorm.DB {
	return db.db
}
