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
		ValidArgs: []string{"current", "last"},
		Args:      cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	}

	cmd.RunE = func(cmd *cobra.Command, args []string) error {

		month := args[0]

		api := zenmoney.NewApi(&http.Client{})

		months := usecase.NewMonth(api)

		months.GetMonthStat(month)

		return nil
	}

	return cmd
}
