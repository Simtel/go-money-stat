package cmd

import (
	"github.com/spf13/cobra"
	"money-stat/internal/services/zenmoney"
	"money-stat/internal/usecase"
	"net/http"
)

func RunYear() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "year",
		Short: "Показать таблицу доходов и расходов  за последние 12 месяцев",
	}

	cmd.RunE = func(cmd *cobra.Command, args []string) error {

		api := zenmoney.NewApi(&http.Client{})

		year := usecase.NewYear(api)

		year.GetYearStat()

		return nil
	}

	return cmd
}
