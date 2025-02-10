package cmd

import (
	"github.com/spf13/cobra"
	app2 "money-stat/internal/app"
	"money-stat/internal/usecase"
)

func RunMonths() *cobra.Command {
	cmd := &cobra.Command{
		Use:       "months",
		Short:     "Показать транзакции за месяц (текущий(current), прошлый(previous))",
		ValidArgs: []string{"current", "previous"},
		Args:      cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	}

	cmd.RunE = func(cmd *cobra.Command, args []string) error {

		app, _ := app2.GetGlobalApp()

		month := args[0]

		months := usecase.NewMonth(app.GetContainer().GetTransactionRepository())

		months.GetMonthStat(month)

		return nil
	}

	return cmd
}
