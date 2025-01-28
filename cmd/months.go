package cmd

import (
	"github.com/spf13/cobra"
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

		month := args[0]

		if len(args) > 1 {
			local := args[1]

			if local == "local" {
				transactionsLocal := &usecase.TransactionsLocal{}
				transactionsLocal.GetLast(10)
				return nil
			}
		}

		api := zenmoney.NewApi(&http.Client{})

		months := usecase.NewMonth(api)

		months.GetMonthStat(month)

		return nil
	}

	return cmd
}
