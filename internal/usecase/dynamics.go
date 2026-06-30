package usecase

import (
	"fmt"
	"money-stat/internal/adapter/sqliterepo/zenrepo/transactions"
	"sort"
	"time"
)

// MonthDynamics — данные о доходах/расходах за месяц с динамикой относительно предыдущего
type MonthDynamics struct {
	Month            string
	Income           float64
	Outcome          float64
	IncomeChange     float64 // абсолютное изменение дохода к прошлому месяцу
	OutcomeChange    float64 // абсолютное изменение расхода к прошлому месяцу
	IncomeChangePct  float64 // изменение дохода в процентах
	OutcomeChangePct float64 // изменение расхода в процентах
}

type Dynamics struct {
	repository transactions.RepositoryInterface
}

func NewDynamics(repository transactions.RepositoryInterface) *Dynamics {
	return &Dynamics{repository: repository}
}

// GetDynamics возвращает помесячную динамику доходов и расходов за указанный год
func (d *Dynamics) GetDynamics(selectYear int) ([]MonthDynamics, error) {
	stats := make(map[string]MonthDynamics)

	transactions, err := d.repository.GetByYear(selectYear)
	if err != nil {
		return nil, fmt.Errorf("получение транзакций за год %d: %w", selectYear, err)
	}

	for _, transaction := range transactions {
		tTime, err := time.Parse("2006-01-02", transaction.Date)
		if err != nil {
			continue
		}
		key := tTime.Format("2006-01")
		stat, exists := stats[key]
		if !exists {
			stat = MonthDynamics{Month: key}
		}
		if transaction.Outcome > 0 && transaction.Income == 0 {
			stat.Outcome += transaction.Outcome
		}
		if transaction.Income > 0 && transaction.Outcome == 0 {
			stat.Income += transaction.Income
		}
		stats[key] = stat
	}

	var valuesSlice []MonthDynamics
	for _, value := range stats {
		valuesSlice = append(valuesSlice, value)
	}

	sort.Slice(valuesSlice, func(i, j int) bool {
		return valuesSlice[i].Month < valuesSlice[j].Month
	})

	// Вычисляем изменения относительно предыдущего месяца
	for i := 1; i < len(valuesSlice); i++ {
		prev := valuesSlice[i-1]
		valuesSlice[i].IncomeChange = valuesSlice[i].Income - prev.Income
		valuesSlice[i].OutcomeChange = valuesSlice[i].Outcome - prev.Outcome

		if prev.Income != 0 {
			valuesSlice[i].IncomeChangePct = (valuesSlice[i].Income - prev.Income) / prev.Income * 100
		}
		if prev.Outcome != 0 {
			valuesSlice[i].OutcomeChangePct = (valuesSlice[i].Outcome - prev.Outcome) / prev.Outcome * 100
		}
	}

	return valuesSlice, nil
}
