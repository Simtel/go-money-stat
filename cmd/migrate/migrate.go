package migrate

import (
	"money-stat/internal/app"
	"money-stat/internal/dbinit"

	"github.com/spf13/cobra"
)

func Run(app *app.App) *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "migrate",
		Short: "Управление миграциями базы данных",
	}

	// Подкоманда для инициализации базы данных
	cmd.AddCommand(&cobra.Command{
		Use:   "init",
		Short: "Инициализировать базу данных (выполнить миграции)",
		RunE: func(cmd *cobra.Command, args []string) error {
			return dbinit.InitializeDB(app.GetContainer().GetDb().GetGorm())
		},
	})

	return cmd
}
