package dbinit

import (
	"gorm.io/gorm"
	"money-stat/internal/model"
)

// InitializeDB инициализирует базу данных, выполняя миграции
func InitializeDB(db *gorm.DB) error {
	return db.AutoMigrate(&model.Transaction{}, &model.Tag{}, &model.Instrument{}, &model.Account{})
}
