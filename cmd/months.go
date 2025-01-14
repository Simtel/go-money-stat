package cmd

import (
	"fmt"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"log"
	"money-stat/internal/services/zenmoney"
	"net/http"
	"strconv"
	"time"
)

func RunMonths() *cobra.Command {
	cmd := &cobra.Command{
		Use:       "months",
		Short:     "Показать транзакции за месяц (текущий, прошлый)",
		ValidArgs: []string{"current", "last"},
		Args:      cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	}

	cmd.RunE = func(cmd *cobra.Command, args []string) error {

		month := args[0]
		log.Printf("Show %s months transactions", month)

		api := zenmoney.NewApi(&http.Client{})

		diff, err := api.Diff()
		if err != nil {
			log.Fatal(err)
		}

		now := time.Now()
		var timestamp int64
		if month == "current" {

			firstDayOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())

			timestamp = firstDayOfMonth.Unix()
		}

		var outComeSumm float64
		var inComeSumm float64
		var cnt int

		tableData := pterm.TableData{
			{"Дата", "Категория", "Сумма"},
			{" ", " ", " "},
		}

		for _, transaction := range diff.Transaction {
			if transaction.Created < timestamp {
				continue
			}

			cnt++
			t := time.Unix(transaction.Changed, 0)
			tableData = append(tableData, []string{t.Format("2006-01-02 15:04:05"), "Категория", strconv.FormatFloat(transaction.Outcome, 'f', 2, 64)})
			if transaction.Outcome > 0 && transaction.Income == 0 {
				outComeSumm = outComeSumm + transaction.Outcome
			}

			if transaction.Income > 0 && transaction.Outcome == 0 {
				inComeSumm = inComeSumm + transaction.Income
			}

		}
		errTable := pterm.DefaultTable.WithHasHeader().WithBoxed().WithRowSeparator("-").WithData(tableData).Render()
		if errTable != nil {
			fmt.Println(errTable)
		}

		summData := pterm.TableData{
			{
				"Транзакций",
				"Доходов в рублях",
				"Расходов в рублях",
			},
			{" ", " "},
			{
				strconv.Itoa(cnt),
				strconv.FormatFloat(inComeSumm, 'f', 2, 64),
				strconv.FormatFloat(outComeSumm, 'f', 2, 64),
			},
		}

		errSummTable := pterm.DefaultTable.WithHasHeader().WithBoxed().WithRowSeparator("-").WithData(summData).Render()
		if errSummTable != nil {
			fmt.Println(errSummTable)
		}

		return nil
	}

	return cmd
}
