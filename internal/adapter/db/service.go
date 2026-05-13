package db

import "gorm.io/gorm"

type DBServiceInterface interface {
	Create(value interface{}) DBServiceInterface
	Where(query interface{}, args ...interface{}) DBServiceInterface
	Delete(value interface{}, conds ...interface{}) DBServiceInterface
	First(dest interface{}, conds ...interface{}) DBServiceInterface
	Updates(value interface{}) DBServiceInterface
	Exec(sql string, values ...interface{}) DBServiceInterface
	GetDB() *gorm.DB
	Model(dest interface{}) DBServiceInterface
	Association(field string) *gorm.Association
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

func (s *DBService) First(dest interface{}, conds ...interface{}) DBServiceInterface {
	s.db.First(dest, conds...)
	return s
}

func (s *DBService) Updates(value interface{}) DBServiceInterface {
	s.db.Updates(value)
	return s
}

func (s *DBService) Exec(sql string, values ...interface{}) DBServiceInterface {
	s.db.Exec(sql, values...)
	return s
}

func (s *DBService) GetDB() *gorm.DB {
	return s.db
}

func (s *DBService) Model(dest interface{}) DBServiceInterface {
	s.db = s.db.Model(dest)
	return s
}

func (s *DBService) Association(field string) *gorm.Association {
	return s.db.Association(field)
}
