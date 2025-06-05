package cmd

import (
	"fmt"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"money-stat/internal/app"
	"money-stat/internal/usecase"
	"strconv"
	"strings"
)

func RunMonths() *cobra.Command {
	cmd := &cobra.Command{
		Use:       "months",
		Short:     "Показать транзакции за месяц (текущий(current), прошлый(previous))",
		ValidArgs: []string{"current", "previous"},
		Args:      cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	}

	cmd.RunE = func(cmd *cobra.Command, args []string) error {

		month := args[0]

		months := usecase.NewMonth(app.GetGlobalApp().GetContainer().GetTransactionRepository())

		stat, err := months.GetMonthStat(month)
		if err != nil {
			return fmt.Errorf("failed to get Month stat: %w", err)
		}

		tableData := pterm.TableData{
			{"Дата", "Категория", "Сумма", "Счет", "Дата создания", "Комментарий"},
			{" ", " ", " ", " ", " "},
		}

		for _, t := range stat.Transactions {
			tableData = append(
				tableData,
				[]string{
					t.Date,
					t.Tags,
					t.FormatAmount,
					t.Account,
					t.CreatedAt,
					t.Comment,
				},
			)
		}

		errTable := pterm.DefaultTable.WithHasHeader().WithBoxed().WithRowSeparator("-").WithData(tableData).Render()
		if errTable != nil {
			fmt.Println(errTable)
		}

		monthDiff := strconv.FormatFloat(stat.IncomeSumm-stat.OutcomeSumm, 'f', 2, 64)
		if strings.HasPrefix(monthDiff, "-") {
			monthDiff = pterm.FgRed.Sprint(monthDiff)
		} else {
			monthDiff = pterm.FgGreen.Sprint(monthDiff)
		}

		summData := pterm.TableData{
			{
				"Транзакций",
				"Доходов в рублях",
				"Расходов в рублях",
				"Чистыми",
			},
			{" ", " ", ""},
			{
				strconv.Itoa(stat.Count),
				strconv.FormatFloat(stat.IncomeSumm, 'f', 2, 64),
				strconv.FormatFloat(stat.OutcomeSumm, 'f', 2, 64),
				monthDiff,
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
