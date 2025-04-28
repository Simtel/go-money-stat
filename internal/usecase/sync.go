package usecase

import (
	"fmt"
	"github.com/pterm/pterm"
	"gorm.io/gorm"
	"log"
	"money-stat/internal/model"
	"money-stat/internal/services/zenmoney"
	"strconv"
)

type Sync struct {
	db  *gorm.DB
	api *zenmoney.Api
}

func NewSync(db *gorm.DB, api *zenmoney.Api) *Sync {
	return &Sync{db: db, api: api}
}

func (s *Sync) FullSync() {
	loadSpin := s.startSpinner("Загрузка данных")

	s.ClearTables()

	diff, _ := s.api.Diff()

	s.stopSpinner(loadSpin, "Загрузка завершена")

	tagsSpin := s.startSpinner("Сохранение тэгов")
	var cntTags int
	for _, tag := range diff.Tag {
		s.db.Create(&model.Tag{
			Id:    tag.Id,
			Title: tag.Title,
		})
		cntTags++
	}
	tagsSpin.Success("Сохранение завершено!")

	intSpinner := s.startSpinner("Сохранение валют")
	var cntInstruments int
	for _, inc := range diff.Instrument {
		s.db.Create(&model.Instrument{
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
		s.db.Create(&model.Account{
			Id:           account.Id,
			Title:        account.Title,
			Balance:      account.Balance,
			Instrument:   account.Instrument,
			StartBalance: account.StartBalance,
		})
		cntAccounts++
	}
	accSpinner.Success("Сохранение завершено!")

	p, _ := pterm.DefaultProgressbar.WithTotal(len(diff.Transaction)).WithTitle("Сохранение транзакций").Start()
	var cntTransactions int
	for _, transaction := range diff.Transaction {
		p.UpdateTitle("Сохранение транзакции " + strconv.Itoa(cntTransactions))
		var tags []model.Tag
		for _, tag := range transaction.Tag {
			tags = append(tags, model.Tag{Id: tag})
		}
		s.db.Create(&model.Transaction{
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
		p.Increment()

	}

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

func (s *Sync) ClearTables() {
	s.db.Where("`id` != ?", "").Delete(&model.Transaction{})
	s.db.Where("`id` != ?", "").Delete(&model.Account{})
	s.db.Where("`id` != ?", "").Delete(&model.Tag{})
	s.db.Where("`id` != ?", "").Delete(&model.Instrument{})
}
