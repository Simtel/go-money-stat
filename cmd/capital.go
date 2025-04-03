package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"money-stat/internal/app"
	"money-stat/internal/usecase"
	"strconv"
	"time"
)

func RunCapital() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "capital",
		Short: "Показать капитал за указанный год",
		Args: func(cmd *cobra.Command, args []string) error {
			if err := cobra.MinimumNArgs(1)(cmd, args); err != nil {
				return err
			}

			selectYear, errYear := strconv.Atoi(args[0])

			if errYear != nil {
				return errYear
			}

			if selectYear > 2020 && selectYear <= time.Now().Year() {
				return nil
			}
			return fmt.Errorf("указан неверный год: %s", args[0])
		},
	}

	cmd.RunE = func(cmd *cobra.Command, args []string) error {

		selectYear, _ := strconv.Atoi(args[0])

		capital := usecase.NewCapital(app.GetGlobalApp().GetContainer().GetTransactionRepository())

		capital.GetCapital(selectYear)

		return nil
	}

	return cmd
}
