package transactions

import (
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"money-stat/internal/model"
	"testing"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)

	err = db.AutoMigrate(&model.Transaction{}, &model.Account{}, &model.Instrument{})
	assert.NoError(t, err)

	// Создаём тестовые валюты и счета
	db.Create(&model.Instrument{Id: 1, Title: "USD", ShortTitle: "USD", Symbol: "$", Rate: 75.0})
	db.Create(&model.Account{Id: "acc1", Title: "Счёт 1", Balance: 0, StartBalance: 0, Instrument: 1})
	db.Create(&model.Account{Id: "acc2", Title: "Счёт 2", Balance: 0, StartBalance: 0, Instrument: 1})

	return db
}

func TestRepository_GetByYear(t *testing.T) {
	db := setupTestDB(t)
	repo := &Repository{db: db}

	db.Create(&model.Transaction{
		Id: "1", Date: "2021-09-01", Income: 100, Outcome: 0,
		IncomeAccount: "acc1", OutcomeAccount: "acc2", Deleted: false,
	})
	db.Create(&model.Transaction{
		Id: "2", Date: "2021-09-02", Income: 200, Outcome: 0,
		IncomeAccount: "acc1", OutcomeAccount: "acc2", Deleted: false,
	})
	// Транзакция вне диапазона
	db.Create(&model.Transaction{
		Id: "3", Date: "2022-01-01", Income: 300, Outcome: 0,
		IncomeAccount: "acc1", OutcomeAccount: "acc2", Deleted: false,
	})

	transactions, err := repo.GetByYear(2021)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(transactions))
	// Проверяем, что InAccount/OutAccount загружены
	assert.Equal(t, "Счёт 1", transactions[0].InAccount.Title)
	assert.Equal(t, "Счёт 2", transactions[0].OutAccount.Title)
}

func TestRepository_GetAll(t *testing.T) {
	db := setupTestDB(t)
	repo := &Repository{db: db}

	// Создаём обычную и удалённую транзакции
	db.Create(&model.Transaction{
		Id: "1", Date: "2021-09-01", Income: 100, Outcome: 0,
		IncomeAccount: "acc1", OutcomeAccount: "acc2", Deleted: false,
	})
	db.Create(&model.Transaction{
		Id: "2", Date: "2021-09-02", Income: 200, Outcome: 0,
		IncomeAccount: "acc1", OutcomeAccount: "acc2", Deleted: true,
	})

	transactions, err := repo.GetAll()
	assert.NoError(t, err)
	assert.Equal(t, 1, len(transactions))
	assert.Equal(t, "Счёт 1", transactions[0].InAccount.Title)
	assert.Equal(t, 75.0, transactions[0].InAccount.Currency.Rate)
}

func TestRepository_GetCurrentMonth(t *testing.T) {
	db := setupTestDB(t)
	repo := &Repository{db: db}

	// Этот тест использует текущую дату, поэтому просто проверяем, что метод не паникует
	transactions := repo.GetCurrentMonth()
	assert.NotNil(t, transactions)
}

func TestRepository_GetPreviousMonth(t *testing.T) {
	db := setupTestDB(t)
	repo := &Repository{db: db}

	transactions := repo.GetPreviousMonth()
	assert.NotNil(t, transactions)
}
