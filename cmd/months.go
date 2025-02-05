package cmd

import (
	"github.com/spf13/cobra"
	app2 "money-stat/internal/app"
	"money-stat/internal/services/zenmoney"
	"money-stat/internal/usecase"
	"net/http"
)

func RunMonths() *cobra.Command {
	cmd := &cobra.Command{
		Use:       "months",
		Short:     "Показать транзакции за месяц (текущий, прошлый)",
		ValidArgs: []string{"current", "last", "local"},
		Args:      cobra.MatchAll(cobra.RangeArgs(1, 2), cobra.OnlyValidArgs),
	}

	cmd.RunE = func(cmd *cobra.Command, args []string) error {

		app, _ := app2.GetGlobalApp()
		month := args[0]

		api := zenmoney.NewApi(&http.Client{})

		months := usecase.NewMonth(api, app.GetContainer().GetTransactionRepository())

		months.GetMonthStat(month)

		return nil
	}

	return cmd
}
