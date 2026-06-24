package app

import (
	"money-stat/internal/dbinit"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DB struct {
	db *gorm.DB
}

type DbInterface interface {
	GetGorm() *gorm.DB
}

func NewDB() DbInterface {
	db, err := gorm.Open(sqlite.Open("zenmoney.db?cache=shared&mode=rwc"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Error),
	})
	if err != nil {
		panic("failed to connect database")
	}

	// Инициализация базы данных через отдельный пакет
	if err := dbinit.InitializeDB(db); err != nil {
		panic(err)
	}
	return &DB{db: db}
}

func (db *DB) GetGorm() *gorm.DB {
	return db.db
}
