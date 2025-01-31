package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"money-stat/internal/model"
	"money-stat/internal/services/zenmoney"
	"net/http"
	"strconv"
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
		errMigrate := db.AutoMigrate(&model.Transaction{}, &model.Tag{}, &model.Instrument{}, &model.Account{})
		if err != nil {
			panic(errMigrate)
		}

		result := db.Where("`id` != ?", "").Delete(&model.Transaction{})

		if result.Error != nil {
			fmt.Println(result.Error)
		}

		fmt.Println("Удалено записей:" + strconv.FormatInt(result.RowsAffected, 16))

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

		for _, tag := range diff.Tag {
			db.Create(&model.Tag{
				Id:    tag.Id,
				Title: tag.Title,
			})
		}

		for _, inc := range diff.Instrument {
			db.Create(&model.Instrument{
				Id:         inc.Id,
				Title:      inc.Title,
				ShortTitle: inc.ShortTitle,
				Symbol:     inc.Symbol,
				Rate:       inc.Rate,
			})
		}

		for _, account := range diff.Account {
			db.Create(&model.Account{
				Id:         account.Id,
				Title:      account.Title,
				Balance:    account.Balance,
				Instrument: account.Instrument,
			})
		}

		return nil
	}

	return cmd
}
