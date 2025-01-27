package cmd

import (
	"github.com/spf13/cobra"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"money-stat/internal/model"
	"money-stat/internal/services/zenmoney"
	"net/http"
)

func RunSync() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sync",
		Short: "Синхронизировать данные",
	}

	cmd.RunE = func(cmd *cobra.Command, args []string) error {

		db, err := gorm.Open(sqlite.Open("zenmoney.db?cache=shared&mode=rwc"), &gorm.Config{})
		if err != nil {
			panic("failed to connect database")
		}

		// Migrate the schema
		errMigrate := db.AutoMigrate(&model.Transaction{})
		if err != nil {
			panic(errMigrate)
		}

		db.Where("`id` > ?", 0).Delete(&model.Transaction{})

		api := zenmoney.NewApi(&http.Client{})

		diff, _ := api.Diff()

		for _, transaction := range diff.Transaction {

			db.Create(&model.Transaction{
				Id:                transaction.Id,
				Changed:           transaction.Changed,
				Created:           transaction.Created,
				IncomeInstrument:  transaction.IncomeInstrument,
				Income:            transaction.Income,
				OutcomeInstrument: transaction.OutcomeInstrument,
				Outcome:           transaction.Outcome,
				Date:              transaction.Date,
				Deleted:           transaction.Deleted,
				IncomeAccount:     transaction.IncomeAccount,
				OutcomeAccount:    transaction.OutcomeAccount,
			})
		}

		return nil
	}

	return cmd
}
