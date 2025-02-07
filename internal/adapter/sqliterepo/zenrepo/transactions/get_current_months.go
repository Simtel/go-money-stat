package transactions

import (
	"fmt"
	"log"
	"money-stat/internal/model"
	"time"
)

func (r *Repository) GetCurrentMonth() []model.Transaction {

	now := time.Now()

	firstDayOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	firstOfNextMonth := time.Date(now.Year(), now.Month()+1, 1, 0, 0, 0, 0, now.Location())
	lastOfCurrentMonth := firstOfNextMonth.AddDate(0, 0, -1)

	return r.GetBetweenDate(firstDayOfMonth, lastOfCurrentMonth)
}

func (r *Repository) GetPreviousMonth() []model.Transaction {
	now := time.Now()
	firstDayOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	previousMonth := firstDayOfMonth.AddDate(0, -1, 0)

	firstOfNextMonth := time.Date(previousMonth.Year(), previousMonth.Month()+1, 1, 23, 59, 59, 0, now.Location())
	lastDayMonth := firstOfNextMonth.AddDate(0, 0, -1)

	return r.GetBetweenDate(previousMonth, lastDayMonth)
}

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
