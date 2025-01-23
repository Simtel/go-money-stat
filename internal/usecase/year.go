package usecase

import (
	"fmt"
	"github.com/pterm/pterm"
	"money-stat/internal/services/zenmoney"
	"net/http"
	"sort"
	"strconv"
	"time"
)

type MonthStat struct {
	Month   string
	Income  float64
	OutCome float64
}

type Year struct {
	api *zenmoney.Api
}

func NewYear(api *zenmoney.Api) *Year {
	return &Year{api: api}
}

func (y *Year) GetYearStat() {
	api := zenmoney.NewApi(&http.Client{})

	stats := make(map[string]MonthStat)

	diff, _ := api.Diff()

	for _, transaction := range diff.Transaction {
		layout := "2006-01-02"
		tTime, _ := time.Parse(layout, transaction.Date)
		key := tTime.Format("2006-01")
		if tTime.Format("2006") < "2020" || transaction.IsDeleted() {
			continue
		}
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

	tableData := pterm.TableData{
		{"Месяц", "Доход", "Расход", "Чистыми"},
		{" ", " ", " ", " "},
	}

	for _, row := range valuesSlice {
		tableData = append(
			tableData,
			[]string{
				row.Month,
				strconv.FormatFloat(row.Income, 'f', 2, 64),
				strconv.FormatFloat(row.OutCome, 'f', 2, 64),
				strconv.FormatFloat(row.Income-row.OutCome, 'f', 2, 64),
			},
		)

	}

	errTable := pterm.DefaultTable.WithHasHeader().WithBoxed().WithRowSeparator("-").WithData(tableData).Render()
	if errTable != nil {
		fmt.Println(errTable)
	}

}
