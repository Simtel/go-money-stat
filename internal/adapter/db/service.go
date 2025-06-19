package db

import "gorm.io/gorm"

type DBServiceInterface interface {
	Create(value interface{}) DBServiceInterface
	Where(query interface{}, args ...interface{}) DBServiceInterface
	Delete(value interface{}, conds ...interface{}) DBServiceInterface
}

type DBService struct {
	db *gorm.DB
}

func NewDBService(db *gorm.DB) DBServiceInterface {
	return &DBService{db: db}
}

func (s *DBService) Create(value interface{}) (tx DBServiceInterface) {
	s.db.Create(value)
	return s
}

func (s *DBService) Where(query interface{}, args ...interface{}) (tx DBServiceInterface) {
	return NewDBService(s.db.Where(query, args))
}

func (s *DBService) Delete(value interface{}, conds ...interface{}) (tx DBServiceInterface) {
	return NewDBService(s.db.Delete(value, conds))
}
