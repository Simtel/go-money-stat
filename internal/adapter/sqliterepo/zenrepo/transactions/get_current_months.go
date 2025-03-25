package transactions

import (
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
