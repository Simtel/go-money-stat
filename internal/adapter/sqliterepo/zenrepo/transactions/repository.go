package transactions

import (
	"gorm.io/gorm"
	"money-stat/internal/model"
	"time"
)

type Repository struct {
	db *gorm.DB
}

type RepositoryInterface interface {
	GetCurrentMonth() []model.Transaction
	GetPreviousMonth() []model.Transaction
	GetBetweenDate(first time.Time, last time.Time) []model.Transaction
	GetAll() ([]model.Transaction, error)
	GetByYear(year int) []model.Transaction
}

func NewRepository(db *gorm.DB) RepositoryInterface {
	return &Repository{db: db}
}
