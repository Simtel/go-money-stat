package transactions

import (
	"fmt"
	"log"
	"money-stat/internal/model"
	"time"
)

func (r *Repository) GetBetweenDate(first time.Time, last time.Time) []model.Transaction {
	db := r.db
	fmt.Println("Ищем транзакции между " + first.Format("2006-01-02") + " и " + last.Format("2006-01-02"))
	var transactions []model.Transaction
	err := db.Model(&model.Transaction{}).
		Where("date BETWEEN ? and ? AND deleted = ?", first.Format("2006-01-02"), last.Format("2006-01-02"), 0).
		Preload("Tag").
		Joins("InAccount").
		Joins("OutAccount").
		Order("date ASC").
		Find(&transactions).
		Error
	if err != nil {
		log.Println(err)
	}
	return transactions
}
