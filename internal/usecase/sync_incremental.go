package usecase

import (
	"fmt"
	"log"
	"money-stat/internal/adapter/db"
	"money-stat/internal/model"
	"money-stat/internal/services/zenmoney"
	"strconv"
	"time"

	"github.com/pterm/pterm"
)

type Sync struct {
	db  db.DBServiceInterface
	api zenmoney.ApiInterface
}

func NewSync(db db.DBServiceInterface, api zenmoney.ApiInterface) *Sync {
	return &Sync{db: db, api: api}
}

// FullSync выполняет полную синхронизацию с очисткой всех данных
func (s *Sync) FullSync() {
	clearSpin := s.startSpinner("Очистка таблиц")
	s.ClearTables()
	s.stopSpinner(clearSpin, "Очистка завершена")

	s.syncData(0, true)
}

// IncrementalSync выполняет инкрементальную синхронизацию только новых/измененных данных
func (s *Sync) IncrementalSync() {
	// Получаем состояние последней синхронизации
	var syncState model.SyncState
	result := s.db.First(&syncState, "id = ?", "main")

	if result.GetDB().Error != nil {
		// Если запись не найдена, выполняем полную синхронизацию
		log.Println("Первая синхронизация, выполняем полный импорт данных")
		s.FullSync()
		return
	}

	lastTimestamp := syncState.LastSyncedAt
	log.Printf("Последняя синхронизация: %s (timestamp: %d)", time.Unix(lastTimestamp, 0).Format("2006-01-02 15:04:05"), lastTimestamp)

	s.syncData(lastTimestamp, false)
}

// syncData выполняет синхронизацию данных с указанным timestamp
func (s *Sync) syncData(lastTimestamp int64, isFull bool) {
	loadSpin := s.startSpinner("Загрузка данных с сервера")

	var diff *zenmoney.Response
	var err error

	if isFull {
		diff, err = s.api.Diff()
	} else {
		diff, err = s.api.DiffSince(lastTimestamp)
	}

	if err != nil {
		s.stopSpinner(loadSpin, "Ошибка загрузки данных: "+err.Error())
		log.Fatal(err)
	}
	s.stopSpinner(loadSpin, "Загрузка завершена")

	// Сохраняем теги
	tagsSpin := s.startSpinner("Сохранение тегов")
	var cntTags int
	for _, tag := range diff.Tag {
		s.db.Create(&model.Tag{
			Id:    tag.Id,
			Title: tag.Title,
		})
		cntTags++
	}
	s.stopSpinner(tagsSpin, fmt.Sprintf("Сохранено тегов: %d", cntTags))

	// Сохраняем валюты
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
	s.stopSpinner(intSpinner, fmt.Sprintf("Сохранено валют: %d", cntInstruments))

	// Сохраняем счета с использованием Upsert
	accSpinner := s.startSpinner("Сохранение счетов")
	var cntAccounts int
	for _, account := range diff.Account {
		s.upsertAccount(&account)
		cntAccounts++
	}
	s.stopSpinner(accSpinner, fmt.Sprintf("Сохранено счетов: %d", cntAccounts))

	// Сохраняем транзакции с использованием Upsert
	p, _ := pterm.DefaultProgressbar.WithTotal(len(diff.Transaction)).WithTitle("Сохранение транзакций").Start()
	var cntTransactions int
	for _, transaction := range diff.Transaction {
		p.UpdateTitle("Сохранение транзакции " + strconv.Itoa(cntTransactions))
		s.upsertTransaction(&transaction)
		cntTransactions++
		p.Increment()
	}
	p.Stop()

	fmt.Println("Загружено транзакций:", cntTransactions)
	fmt.Println("Загружено тегов:", cntTags)
	fmt.Println("Загружено счетов:", cntAccounts)
	fmt.Println("Загружено валют:", cntInstruments)

	// Обновляем состояние синхронизации
	if diff.ServerTimestamp > 0 {
		s.updateSyncState(diff.ServerTimestamp)
	}
}

// upsertAccount создает или обновляет счет
func (s *Sync) upsertAccount(account *zenmoney.Account) {
	var existing model.Account
	result := s.db.First(&existing, "id = ?", account.Id)

	if result.GetDB().Error == nil {
		// Счет существует, обновляем
		s.db.Model(&existing).Updates(map[string]interface{}{
			"Title":        account.Title,
			"Balance":      account.Balance,
			"Instrument":   account.Instrument,
			"StartBalance": account.StartBalance,
		})
	} else {
		// Счет не существует, создаем
		s.db.Create(&model.Account{
			Id:           account.Id,
			Title:        account.Title,
			Balance:      account.Balance,
			Instrument:   account.Instrument,
			StartBalance: account.StartBalance,
		})
	}
}

// upsertTransaction создает или обновляет транзакцию
func (s *Sync) upsertTransaction(transaction *zenmoney.Transaction) {
	var existing model.Transaction
	result := s.db.Where("id = ?", transaction.Id).First(&existing)

	var tags []model.Tag
	for _, tagId := range transaction.Tag {
		tags = append(tags, model.Tag{Id: tagId})
	}

	if result.GetDB().Error == nil {
		// Транзакция существует, обновляем
		s.db.Model(&existing).Updates(map[string]interface{}{
			"Changed":           transaction.Changed,
			"Created":           transaction.Created,
			"IncomeInstrument":  transaction.IncomeInstrument,
			"Income":            transaction.Income,
			"OutcomeInstrument": transaction.OutcomeInstrument,
			"Outcome":           transaction.Outcome,
			"Date":              transaction.Date,
			"Deleted":           transaction.Deleted,
			"IncomeAccount":     transaction.IncomeAccount,
			"OutcomeAccount":    transaction.OutcomeAccount,
			"Comment":           transaction.Comment,
		})
		// Обновляем связи с тегами только если есть теги
		if len(tags) > 0 {
			s.db.Association("Tag").Replace(tags)
		}
	} else {
		// Транзакция не существует, создаем
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
			Comment:           transaction.Comment,
		})
	}
}

// updateSyncState обновляет состояние последней синхронизации
func (s *Sync) updateSyncState(serverTimestamp int64) {
	var syncState model.SyncState
	result := s.db.First(&syncState, "id = ?", "main")

	now := time.Now().Unix()

	if result.GetDB().Error == nil {
		// Запись существует, обновляем
		s.db.GetDB().Model(&syncState).Updates(map[string]interface{}{
			"LastSyncedAt":    now,
			"ServerTimestamp": serverTimestamp,
			"UpdatedAt":       now,
		})
	} else {
		// Запись не существует, создаем
		s.db.Create(&model.SyncState{
			ID:              "main",
			LastSyncedAt:    now,
			ServerTimestamp: serverTimestamp,
			UpdatedAt:       now,
		})
	}

	log.Printf("Состояние синхронизации обновлено: lastSyncedAt=%d, serverTimestamp=%d", now, serverTimestamp)
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
	// Не удаляем sync_state при полной синхронизации
}
