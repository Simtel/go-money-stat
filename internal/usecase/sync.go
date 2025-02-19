package usecase

import (
	"fmt"
	"github.com/pterm/pterm"
	"log"
	app2 "money-stat/internal/app"
	"money-stat/internal/model"
	"money-stat/internal/services/zenmoney"
	"net/http"
	"strconv"
)

type Sync struct{}

func (s *Sync) FullSync() {
	loadSpin := s.startSpinner("Загрузка данных")

	app, _ := app2.GetGlobalApp()

	db := app.GetContainer().GetDb().GetGorm()
	db.Where("`id` != ?", "").Delete(&model.Transaction{})
	db.Where("`id` != ?", "").Delete(&model.Account{})
	db.Where("`id` != ?", "").Delete(&model.Tag{})
	db.Where("`id` != ?", "").Delete(&model.Instrument{})

	api := zenmoney.NewApi(&http.Client{})

	diff, _ := api.Diff()

	s.stopSpinner(loadSpin, "Загрузка завершена")

	tagsSpin := s.startSpinner("Сохранение тэгов")
	var cntTags int
	for _, tag := range diff.Tag {
		db.Create(&model.Tag{
			Id:    tag.Id,
			Title: tag.Title,
		})
		cntTags++
	}
	tagsSpin.Success("Сохранение завершено!")

	intSpinner := s.startSpinner("Сохранение валют")
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

	accSpinner := s.startSpinner("Сохранение счетов")
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

	trSpinner := s.startSpinner("Сохранение транзакций")
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
	s.stopSpinner(trSpinner, "Сохранение завершено!")

	fmt.Println("Загружено транзакций:" + strconv.Itoa(cntTransactions))
	fmt.Println("Загружено тэгов:" + strconv.Itoa(cntTags))
	fmt.Println("Загружено счетов:" + strconv.Itoa(cntAccounts))
	fmt.Println("Загружено валют:" + strconv.Itoa(cntInstruments))

}

func (s *Sync) startSpinner(title string) *pterm.SpinnerPrinter {
	multi := pterm.DefaultMultiPrinter
	loadSpinner, _ := pterm.DefaultSpinner.WithWriter(multi.NewWriter()).Start(title)
	_, err := multi.Start()
	if err != nil {
		log.Fatal(err)
	}
	return loadSpinner
}

func (s *Sync) stopSpinner(spinner *pterm.SpinnerPrinter, title string) {
	spinner.Success(title)
}
