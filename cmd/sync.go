package cmd

import (
	"github.com/spf13/cobra"
	app "money-stat/internal/app"
	"money-stat/internal/services/zenmoney"
	"money-stat/internal/usecase"
	"net/http"
)

func RunSync() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sync",
		Short: "Синхронизировать данные",
	}

	cmd.RunE = func(cmd *cobra.Command, args []string) error {

		application, _ := app.GetGlobalApp()

		db := application.GetContainer().GetDb().GetGorm()
		api := zenmoney.NewApi(&http.Client{})
		sync := usecase.NewSync(db, api)

		sync.FullSync()

		return nil
	}

	return cmd
}
