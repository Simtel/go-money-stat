package cmd

import (
	"github.com/spf13/cobra"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
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
		db.AutoMigrate(&Transaction{})

		db.Where("`id` > ?", 0).Delete(&Transaction{})

		api := zenmoney.NewApi(&http.Client{})

		diff, _ := api.Diff()

		for _, transaction := range diff.Transaction {

			db.Create(&Transaction{
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

type Transaction struct {
	gorm.Model
	Id                string
	Changed           int64
	Created           int64
	IncomeInstrument  int64
	Income            float64
	OutcomeInstrument int64
	Outcome           float64
	Date              string
	Deleted           bool
	IncomeAccount     string
	OutcomeAccount    string
}
