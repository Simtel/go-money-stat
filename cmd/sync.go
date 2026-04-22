package cmd

import (
	"github.com/spf13/cobra"
	"money-stat/internal/adapter/db"
	app "money-stat/internal/app"
	"money-stat/internal/services/zenmoney"
	"money-stat/internal/usecase"
	"net/http"
	"time"
)

func RunSync(app *app.App) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sync",
		Short: "Синхронизировать данные",
	}

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		gorm := app.GetContainer().GetDb().GetGorm()
		client := &http.Client{Timeout: 30 * time.Second}
		api := zenmoney.NewApi(client)
		sync := usecase.NewSync(db.NewDBService(gorm), api)

		sync.FullSync()

		return nil
	}

	return cmd
}
