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

// GetDB возвращает экземпляр DBServiceInterface
func (s *Sync) GetDB() db.DBServiceInterface {
	return s.db
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

	// Конвертируем теги из zenmoney.Tag в model.Tag
	tags := make([]model.Tag, len(diff.Tag))
	for i, t := range diff.Tag {
		tags[i] = model.Tag{Id: t.Id}
	}

	// Конвертируем валюты из zenmoney.Instrument в model.Instrument
	instruments := make([]model.Instrument, len(diff.Instrument))
	for i, inst := range diff.Instrument {
		instruments[i] = model.Instrument{
			Id:         inst.Id,
			Title:      inst.Title,
			ShortTitle: inst.ShortTitle,
			Symbol:     inst.Symbol,
			Rate:       inst.Rate,
		}
	}

	// Сохраняем теги
	cntTags := 0
	if err := s.saveTags(tags); err != nil {
		log.Printf("Ошибка сохранения тегов: %v", err)
	}
	cntTags = len(tags)

	// Сохраняем валюты
	cntInstruments := 0
	if err := s.saveInstruments(instruments); err != nil {
		log.Printf("Ошибка сохранения валют: %v", err)
	}
	cntInstruments = len(instruments)

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
	_, _ = p.Stop()

	fmt.Println("Загружено транзакций:", cntTransactions)
	fmt.Println("Загружено тегов:", cntTags)
	fmt.Println("Загружено счетов:", cntAccounts)
	fmt.Println("Загружено валют:", cntInstruments)
	log.Printf("ServerTimestamp из API: %d", diff.ServerTimestamp)

	// Обновляем состояние синхронизации
	if diff.ServerTimestamp > 0 {
		s.updateSyncState(diff.ServerTimestamp)
	} else {
		log.Printf("ServerTimestamp равен 0, не сохраняем sync_state")
	}
}

// saveTags сохраняет теги
func (s *Sync) saveTags(tags []model.Tag) error {
	for i := range tags {
		var existing model.Tag
		result := s.db.Where("id = ?", tags[i].Id).First(&existing)
		if result.GetDB().Error == nil {
			tx := s.db.Save(&tags[i])
			if tx.GetDB().Error != nil {
				return fmt.Errorf("ошибка обновления тега %s: %w", tags[i].Id, tx.GetDB().Error)
			}
		} else {
			tx := s.db.Create(&tags[i])
			if tx.GetDB().Error != nil {
				return fmt.Errorf("ошибка создания тега %s: %w", tags[i].Id, tx.GetDB().Error)
			}
		}
	}
	return nil
}

// saveInstruments сохраняет валюты
func (s *Sync) saveInstruments(instruments []model.Instrument) error {
	for i := range instruments {
		var existing model.Instrument
		result := s.db.Where("id = ?", instruments[i].Id).First(&existing)
		if result.GetDB().Error == nil {
			tx := s.db.Save(&instruments[i])
			if tx.GetDB().Error != nil {
				return fmt.Errorf("ошибка обновления валюты %d: %w", instruments[i].Id, tx.GetDB().Error)
			}
		} else {
			tx := s.db.Create(&instruments[i])
			if tx.GetDB().Error != nil {
				return fmt.Errorf("ошибка создания валюты %d: %w", instruments[i].Id, tx.GetDB().Error)
			}
		}
	}
	return nil
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

	// Преобразуем теги в строку через запятую
	var tagIds string
	for i, tagId := range transaction.Tag {
		if i > 0 {
			tagIds += ","
		}
		tagIds += tagId
	}

	// Создаем объект транзакции без ассоциаций
	txn := model.Transaction{
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
		TagIds:            tagIds,
		Comment:           transaction.Comment,
	}

	if result.GetDB().Error == nil {
		// Транзакция существует, обновляем через Save
		s.db.Save(&txn)
	} else {
		// Транзакция не существует, создаем
		s.db.Create(&txn)
	}
}

// updateSyncState обновляет состояние последней синхронизации
func (s *Sync) updateSyncState(serverTimestamp int64) {
	now := time.Now().Unix()
	syncState := model.SyncState{
		ID:              "main",
		LastSyncedAt:    now,
		ServerTimestamp: serverTimestamp,
		UpdatedAt:       now,
	}
	// Используем Save для обновления или создания записи
	tx := s.db.Save(&syncState)

	// Проверяем ошибку
	if tx.GetDB().Error != nil {
		log.Printf("Ошибка сохранения sync_state: %v", tx.GetDB().Error)
	} else {
		log.Printf("Состояние синхронизации сохранено: lastSyncedAt=%d, serverTimestamp=%d", now, serverTimestamp)
	}
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