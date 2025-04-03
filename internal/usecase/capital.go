package usecase

import (
	"fmt"
	"github.com/pterm/pterm"
	transactionsRepo "money-stat/internal/adapter/sqliterepo/zenrepo/transactions"
	"sort"
	"strconv"
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

func (c *Capital) GetCapital(year int) {
	tableData := pterm.TableData{
		{"Месяц", "Капитал"},
		{" ", " "},
	}

	stats := make(map[string]CapitalDto)

	transactions := c.repo.GetAll()

	for _, transaction := range transactions {
		layout := "2006-01-02"
		tTime, _ := time.Parse(layout, transaction.Date)
		key := tTime.Format("2006-01")
		stat, exists := stats[key]
		if !exists {
			stat = CapitalDto{Month: tTime.Format("2006-01")}
		}
		if transaction.Outcome > 0 && transaction.Income == 0 {
			stat.Balance = stat.Balance - transaction.Outcome
		}

		if transaction.Income > 0 && transaction.Outcome == 0 {
			stat.Balance = stat.Balance + transaction.Income
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
	var summ float64
	for index, row := range valuesSlice {
		if index == 0 {
			summ = 0
		}
		summ = summ + row.Balance

		tableData = append(
			tableData,
			[]string{
				row.Month,
				strconv.FormatFloat(summ, 'f', 2, 64),
			},
		)

	}

	errTable := pterm.DefaultTable.WithHasHeader().WithBoxed().WithRowSeparator("-").WithData(tableData).Render()
	if errTable != nil {
		fmt.Println(errTable)
	}
}
