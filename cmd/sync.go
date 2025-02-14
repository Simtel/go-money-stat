package cmd

import (
	"fmt"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"log"
	app2 "money-stat/internal/app"
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

		loadSpin := startSpinner("Загрузка данных")

		app, _ := app2.GetGlobalApp()

		db := app.GetContainer().GetDb().GetGorm()
		db.Where("`id` != ?", "").Delete(&model.Transaction{})
		db.Where("`id` != ?", "").Delete(&model.Account{})
		db.Where("`id` != ?", "").Delete(&model.Tag{})
		db.Where("`id` != ?", "").Delete(&model.Instrument{})

		api := zenmoney.NewApi(&http.Client{})

		diff, _ := api.Diff()

		stopSpinner(loadSpin, "Загрузка завершена")

		tagsSpin := startSpinner("Сохранение тэгов")
		var cntTags int
		for _, tag := range diff.Tag {
			db.Create(&model.Tag{
				Id:    tag.Id,
				Title: tag.Title,
			})
			cntTags++
		}
		tagsSpin.Success("Сохранение завершено!")

		intSpinner := startSpinner("Сохранение валют")
		var cntInstruments int
		for _, inc := range diff.Instrument {
			db.Create(&model.Instrument{
				Id:         inc.Id,
				Title:      inc.Title,
				ShortTitle: inc.ShortTitle,
				Symbol:     inc.Symbol,
				Rate:       inc.Rate,
			})
			cntInstruments++
		}
		intSpinner.Success("Сохранение завершено!")

		accSpinner := startSpinner("Сохранение счетов")
		var cntAccounts int
		for _, account := range diff.Account {
			db.Create(&model.Account{
				Id:         account.Id,
				Title:      account.Title,
				Balance:    account.Balance,
				Instrument: account.Instrument,
			})
			cntAccounts++
		}
		accSpinner.Success("Сохранение завершено!")

		trSpinner := startSpinner("Сохранение транзакций")
		var cntTransactions int
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
			cntTransactions++
		}
		stopSpinner(trSpinner, "Сохранение завершено!")

		fmt.Println("Загружено транзакций:" + strconv.Itoa(cntTransactions))
		fmt.Println("Загружено тэгов:" + strconv.Itoa(cntTags))
		fmt.Println("Загружено счетов:" + strconv.Itoa(cntAccounts))
		fmt.Println("Загружено валют:" + strconv.Itoa(cntInstruments))
		return nil
	}

	return cmd
}

func startSpinner(title string) *pterm.SpinnerPrinter {
	multi := pterm.DefaultMultiPrinter
	loadSpinner, _ := pterm.DefaultSpinner.WithWriter(multi.NewWriter()).Start(title)
	_, err := multi.Start()
	if err != nil {
		log.Fatal(err)
	}
	return loadSpinner
}

func stopSpinner(spinner *pterm.SpinnerPrinter, title string) {
	spinner.Success(title)
}
