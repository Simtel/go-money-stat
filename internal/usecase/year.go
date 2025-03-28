package usecase

import (
	"money-stat/internal/adapter/sqliterepo/zenrepo/transactions"
	"sort"
	"time"
)

type MonthStat struct {
	Month   string
	Income  float64
	OutCome float64
}

type YearInterface interface {
	GetYearStat()
}
type Year struct {
	repository transactions.RepositoryInterface
}

func NewYear(repository transactions.RepositoryInterface) *Year {
	return &Year{repository: repository}
}

func (y *Year) GetYearStat(selectYear int) []MonthStat {

	stats := make(map[string]MonthStat)

	for _, transaction := range y.repository.GetByYear(selectYear) {
		layout := "2006-01-02"
		tTime, _ := time.Parse(layout, transaction.Date)
		key := tTime.Format("2006-01")
		stat, exists := stats[key]
		if !exists {
			stat = MonthStat{Month: tTime.Format("2006-01")}
		}
		if transaction.Outcome > 0 && transaction.Income == 0 {
			stat.OutCome = stat.OutCome + transaction.Outcome
		}

		if transaction.Income > 0 && transaction.Outcome == 0 {
			stat.Income = stat.Income + transaction.Income
		}

		stats[key] = stat
	}

	var valuesSlice []MonthStat
	for _, value := range stats {
		valuesSlice = append(valuesSlice, value)
	}

	sort.Slice(valuesSlice, func(i, j int) bool {
		return valuesSlice[i].Month < valuesSlice[j].Month
	})

	return valuesSlice

}
