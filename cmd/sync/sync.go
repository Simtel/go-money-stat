package sync

import (
	"github.com/spf13/cobra"
	"money-stat/internal/adapter/db"
	"money-stat/internal/app"
	"money-stat/internal/services/zenmoney"
	"money-stat/internal/usecase"
	"net/http"
	"time"
)

func Run(app *app.App) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sync",
		Short: "Синхронизировать данные с ZenMoney",
	}

	var incremental bool
	var full bool

	cmd.Flags().BoolVarP(&incremental, "incremental", "i", true, "Инкрементальная синхронизация (только новые/измененные данные)")
	cmd.Flags().BoolVarP(&full, "full", "f", false, "Полная синхронизация (очистка и перезагрузка всех данных)")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		gorm := app.GetContainer().GetDb().GetGorm()
		client := &http.Client{Timeout: 30 * time.Second}
		api := zenmoney.NewApi(client)
		sync := usecase.NewSync(db.NewDBService(gorm), api)

		if full {
			println("Выполняется ПОЛНАЯ синхронизация (все данные будут перезагружены)")
			sync.FullSync()
		} else {
			println("Выполняется ИНКРЕМЕНТАЛЬНАЯ синхронизация (только новые/измененные данные)")
			sync.IncrementalSync()
		}

		return nil
	}

	return cmd
}
