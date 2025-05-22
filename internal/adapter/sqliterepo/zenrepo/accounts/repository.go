package accounts

import (
	"gorm.io/gorm"
	"money-stat/internal/model"
)

type RepositoryInterface interface {
	GetAll() []model.Account
}

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) RepositoryInterface {
	return &Repository{db: db}
}
