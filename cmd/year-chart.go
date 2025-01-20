package cmd

import (
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"github.com/spf13/cobra"
	"log"
	"math"
	"money-stat/internal/services/zenmoney"
	"net/http"
	"time"
)

func RunYearChart() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "year-chart",
		Short: "Показать график доходов и расходов  за последние 12 месяцев",
	}

	cmd.RunE = func(cmd *cobra.Command, args []string) error {

		if err := ui.Init(); err != nil {
			log.Fatalf("failed to initialize termui: %v", err)
		}
		defer ui.Close()

		api := zenmoney.NewApi(&http.Client{})

		stats := make(map[string]MonthStat)

		diff, _ := api.Diff()

		for _, transaction := range diff.Transaction {
			layout := "2006-01-02"
			tTime, _ := time.Parse(layout, transaction.Date)
			key := tTime.Format("2006-01")
			if tTime.Format("2006") < "2020" {
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

			if transaction.Outcome > 0 && transaction.Income > 0 {
				stat.Diff = stat.Diff + math.Abs(transaction.Outcome-transaction.Income)
			}
			stats[key] = stat
		}

		sbc := widgets.NewStackedBarChart()
		sbc.Title = "Статистика доходов и расходов по месецам"
		sbc.Labels = []string{}
		sbc.Data = make([][]float64, len(stats))
		index := 0
		for _, row := range stats {
			sbc.Labels = append(sbc.Labels, row.Month)
			sbc.Data[index] = []float64{row.Income, row.OutCome, row.Diff}
			index++
		}

		sbc.SetRect(5, 5, 1000, 50)
		sbc.BarWidth = 6

		ui.Render(sbc)

		uiEvents := ui.PollEvents()
		for {
			e := <-uiEvents
			switch e.ID {
			case "q", "<C-c>":
				return nil
			}
		}

	}

	return cmd
}

type MonthStat struct {
	Month   string
	Income  float64
	OutCome float64
	Diff    float64
}
