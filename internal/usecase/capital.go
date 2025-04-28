package usecase

import (
	transactionsRepo "money-stat/internal/adapter/sqliterepo/zenrepo/transactions"
	"sort"
	"time"
)

type CapitalDto struct {
	Month   string
	Balance float64
}

type Capital struct {
	repo transactionsRepo.RepositoryInterface
}

func NewCapital(repo transactionsRepo.RepositoryInterface) *Capital {
	return &Capital{repo: repo}
}

func (c *Capital) GetCapital() []CapitalDto {

	stats := make(map[string]CapitalDto)

	var transactions, err = c.repo.GetAll()

	if err != nil {
		panic(err)
	}

	for _, transaction := range transactions {
		layout := "2006-01-02"
		tTime, errTimeParse := time.Parse(layout, transaction.Date)
		if errTimeParse != nil {
			panic(errTimeParse)
		}
		key := tTime.Format("2006-01")
		stat, exists := stats[key]
		if !exists {
			stat = CapitalDto{Month: key}
		}
		if transaction.Outcome > 0 && transaction.Income == 0 {
			if !transaction.OutAccount.IsRuble() {
				stat.Balance = stat.Balance - (transaction.Outcome * transaction.OutAccount.Currency.Rate)
			} else {
				stat.Balance = stat.Balance - transaction.Outcome
			}

		}

		if transaction.Income > 0 && transaction.Outcome == 0 {
			if !transaction.InAccount.IsRuble() {
				stat.Balance = stat.Balance + (transaction.Income * transaction.InAccount.Currency.Rate)
			} else {
				stat.Balance = stat.Balance + transaction.Income
			}
		}

		stats[key] = stat
	}

	var valuesSlice []CapitalDto
	for _, value := range stats {
		valuesSlice = append(valuesSlice, value)
	}

	sort.Slice(valuesSlice, func(i, j int) bool {
		return valuesSlice[i].Month < valuesSlice[j].Month
	})

	return valuesSlice
}
