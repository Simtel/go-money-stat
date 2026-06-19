package db

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type DBServiceInterface interface {
	Create(value interface{}) DBServiceInterface
	Where(query interface{}, args ...interface{}) DBServiceInterface
	Delete(value interface{}, conds ...interface{}) DBServiceInterface
	First(dest interface{}, conds ...interface{}) DBServiceInterface
	Updates(value interface{}) DBServiceInterface
	Save(value interface{}) DBServiceInterface
	Exec(sql string, values ...interface{}) DBServiceInterface
	GetDB() *gorm.DB
	Model(dest interface{}) DBServiceInterface
	Association(field string) *gorm.Association
	Clauses(onConflict *clause.OnConflict) DBServiceInterface
	Select(query interface{}, args ...interface{}) DBServiceInterface
}

type DBService struct {
	db *gorm.DB
}

func NewDBService(db *gorm.DB) DBServiceInterface {
	return &DBService{db: db}
}

func (s *DBService) Create(value interface{}) (tx DBServiceInterface) {
	return NewDBService(s.db.Create(value))
}

func (s *DBService) Where(query interface{}, args ...interface{}) (tx DBServiceInterface) {
	return NewDBService(s.db.Where(query, args))
}

func (s *DBService) Delete(value interface{}, conds ...interface{}) (tx DBServiceInterface) {
	return NewDBService(s.db.Delete(value, conds))
}

func (s *DBService) First(dest interface{}, conds ...interface{}) DBServiceInterface {
	return NewDBService(s.db.First(dest, conds...))
}

func (s *DBService) Updates(value interface{}) DBServiceInterface {
	return NewDBService(s.db.Updates(value))
}

func (s *DBService) Save(value interface{}) DBServiceInterface {
	return NewDBService(s.db.Save(value))
}

func (s *DBService) Exec(sql string, values ...interface{}) DBServiceInterface {
	s.db.Exec(sql, values...)
	return s
}

func (s *DBService) GetDB() *gorm.DB {
	return s.db
}

func (s *DBService) Model(dest interface{}) DBServiceInterface {
	return NewDBService(s.db.Model(dest))
}

func (s *DBService) Association(field string) *gorm.Association {
	return s.db.Association(field)
}

func (s *DBService) Clauses(onConflict *clause.OnConflict) (tx DBServiceInterface) {
	return NewDBService(s.db.Clauses(*onConflict))
}

func (s *DBService) Select(query interface{}, args ...interface{}) (tx DBServiceInterface) {
	return NewDBService(s.db.Select(query, args...))
}