package cmd

import (
	"github.com/spf13/cobra"
	app2 "money-stat/internal/app"
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

		app, _ := app2.GetGlobalApp()

		db := app.GetContainer().GetDb().GetGorm()
		db.Where("`id` != ?", "").Delete(&model.Transaction{})
		db.Where("`id` != ?", "").Delete(&model.Account{})
		db.Where("`id` != ?", "").Delete(&model.Tag{})
		db.Where("`id` != ?", "").Delete(&model.Instrument{})

		api := zenmoney.NewApi(&http.Client{})

		diff, _ := api.Diff()

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

		for _, transaction := range diff.Transaction {

			var tags []model.Tag
			for _, tag := range transaction.Tag {
				tags = append(tags, model.Tag{Id: tag})
			}
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
				Tag:               tags,
			})
		}

		return nil
	}

	return cmd
}
