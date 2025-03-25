package transactions

import (
	"money-stat/internal/model"
	"time"
)

func (r *Repository) GetByYear(year int) []model.Transaction {
	firstDay := time.Date(year, time.January, 1, 0, 0, 0, 0, time.UTC)
	lastDay := time.Date(year, time.December, 31, 23, 59, 59, 59, time.UTC)

	return r.GetBetweenDate(firstDay, lastDay)
}
