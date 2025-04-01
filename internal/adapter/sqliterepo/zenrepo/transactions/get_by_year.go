package transactions

import (
	"log"
	"money-stat/internal/model"
	"time"
)

func (r *Repository) GetByYear(year int) []model.Transaction {
	firstDay := time.Date(year, time.January, 1, 0, 0, 0, 0, time.UTC)
	lastDay := time.Date(year, time.December, 31, 23, 59, 59, 59, time.UTC)

	db := r.db

	var transactions []model.Transaction
	err := db.Model(&model.Transaction{}).
		Where(
			"date BETWEEN ? and ? AND deleted = ?",
			firstDay.Format("2006-01-02"),
			lastDay.Format("2006-01-02"),
			0).
		Order("date ASC").
		Find(&transactions).
		Error
	if err != nil {
		log.Println(err)
	}
	return transactions
}
