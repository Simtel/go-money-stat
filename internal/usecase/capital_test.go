package usecase

import (
	accountsRepo "money-stat/internal/adapter/sqliterepo/zenrepo/accounts/mocks"
	transactionsRepo "money-stat/internal/adapter/sqliterepo/zenrepo/transactions/mocks"
	"money-stat/internal/model"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestGetCapital_EmptyTransactions(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransactionRepo := transactionsRepo.NewMockRepositoryInterface(ctrl)
	mockAccountRepo := accountsRepo.NewMockRepositoryInterface(ctrl)

	mockTransactionRepo.EXPECT().GetAll().Return([]model.Transaction{}, nil)
	mockAccountRepo.EXPECT().GetAll().Return([]model.Account{}, nil)

	capital := NewCapital(mockTransactionRepo, mockAccountRepo)

	result, err := capital.GetCapital(2023)
	assert.NoError(t, err)
	assert.Len(t, result, 12)

	// Без счетов и транзакций — все месяцы с нулевым балансом
	for _, m := range result {
		assert.Equal(t, 0.0, m.Balance)
	}
}

func TestGetCapital_StartBalanceOnly(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransactionRepo := transactionsRepo.NewMockRepositoryInterface(ctrl)
	mockAccountRepo := accountsRepo.NewMockRepositoryInterface(ctrl)

	accounts := []model.Account{
		{
			Id:           "acc1",
			Title:        "Счёт RUB",
			StartBalance: 100000,
			Instrument:   2,
			Currency: model.Instrument{
				Title: "RUB",
				Rate:  1.0,
			},
		},
		{
			Id:           "acc2",
			Title:        "Счёт USD",
			StartBalance: 1000,
			Instrument:   1,
			Currency: model.Instrument{
				Title: "USD",
				Rate:  75.0,
			},
		},
	}

	mockTransactionRepo.EXPECT().GetAll().Return([]model.Transaction{}, nil)
	mockAccountRepo.EXPECT().GetAll().Return(accounts, nil)

	capital := NewCapital(mockTransactionRepo, mockAccountRepo)

	result, err := capital.GetCapital(2023)
	assert.NoError(t, err)
	assert.Len(t, result, 12)

	expectedBalance := 100000.0 + 1000.0*75.0 // 175000
	for _, m := range result {
		assert.Equal(t, expectedBalance, m.Balance)
	}
}

func TestGetCapital_TransactionsInYear(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransactionRepo := transactionsRepo.NewMockRepositoryInterface(ctrl)
	mockAccountRepo := accountsRepo.NewMockRepositoryInterface(ctrl)

	accounts := []model.Account{
		{
			Id:           "acc1",
			Title:        "Счёт",
			StartBalance: 100000,
			Instrument:   2,
			Currency: model.Instrument{
				Title: "RUB",
				Rate:  1.0,
			},
		},
	}

	tx1 := model.Transaction{
		Id:      "1",
		Date:    "2023-01-15",
		Income:  50000,
		Outcome: 0,
		InAccount: model.Account{
			Currency: model.Instrument{Rate: 1.0},
		},
		OutAccount: model.Account{
			Currency: model.Instrument{Rate: 1.0},
		},
		Deleted: false,
	}

	tx2 := model.Transaction{
		Id:      "2",
		Date:    "2023-02-15",
		Income:  0,
		Outcome: 20000,
		InAccount: model.Account{
			Currency: model.Instrument{Rate: 1.0},
		},
		OutAccount: model.Account{
			Currency: model.Instrument{Rate: 1.0},
		},
		Deleted: false,
	}

	mockTransactionRepo.EXPECT().GetAll().Return([]model.Transaction{tx1, tx2}, nil)
	mockAccountRepo.EXPECT().GetAll().Return(accounts, nil)

	capital := NewCapital(mockTransactionRepo, mockAccountRepo)

	result, err := capital.GetCapital(2023)
	assert.NoError(t, err)
	assert.Len(t, result, 12)

	// Январь: 100000 + 50000 = 150000
	assert.Equal(t, 150000.0, result[0].Balance)
	assert.Equal(t, "2023-01", result[0].Month)

	// Февраль: 150000 - 20000 = 130000
	assert.Equal(t, 130000.0, result[1].Balance)
	assert.Equal(t, "2023-02", result[1].Month)

	// Март-декабрь: без изменений = 130000
	for i := 2; i < 12; i++ {
		assert.Equal(t, 130000.0, result[i].Balance)
	}
}

func TestGetCapital_TransactionsBeforeYear(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransactionRepo := transactionsRepo.NewMockRepositoryInterface(ctrl)
	mockAccountRepo := accountsRepo.NewMockRepositoryInterface(ctrl)

	accounts := []model.Account{
		{
			Id:           "acc1",
			Title:        "Счёт",
			StartBalance: 100000,
			Instrument:   2,
			Currency: model.Instrument{
				Title: "RUB",
				Rate:  1.0,
			},
		},
	}

	// Транзакция до запрашиваемого года
	txBefore := model.Transaction{
		Id:      "1",
		Date:    "2022-06-15",
		Income:  30000,
		Outcome: 0,
		InAccount: model.Account{
			Currency: model.Instrument{Rate: 1.0},
		},
		OutAccount: model.Account{
			Currency: model.Instrument{Rate: 1.0},
		},
		Deleted: false,
	}

	// Транзакция в запрашиваемом году
	txIn := model.Transaction{
		Id:      "2",
		Date:    "2023-03-10",
		Income:  0,
		Outcome: 10000,
		InAccount: model.Account{
			Currency: model.Instrument{Rate: 1.0},
		},
		OutAccount: model.Account{
			Currency: model.Instrument{Rate: 1.0},
		},
		Deleted: false,
	}

	mockTransactionRepo.EXPECT().GetAll().Return([]model.Transaction{txBefore, txIn}, nil)
	mockAccountRepo.EXPECT().GetAll().Return(accounts, nil)

	capital := NewCapital(mockTransactionRepo, mockAccountRepo)

	result, err := capital.GetCapital(2023)
	assert.NoError(t, err)
	assert.Len(t, result, 12)

	// Капитал на начало 2023: 100000 + 30000 = 130000
	// Январь: 130000
	assert.Equal(t, 130000.0, result[0].Balance)
	// Февраль: 130000
	assert.Equal(t, 130000.0, result[1].Balance)
	// Март: 130000 - 10000 = 120000
	assert.Equal(t, 120000.0, result[2].Balance)
	// Апрель-декабрь: 120000
	for i := 3; i < 12; i++ {
		assert.Equal(t, 120000.0, result[i].Balance)
	}
}

func TestGetCapital_TransferDoesNotChangeCapital(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransactionRepo := transactionsRepo.NewMockRepositoryInterface(ctrl)
	mockAccountRepo := accountsRepo.NewMockRepositoryInterface(ctrl)

	accounts := []model.Account{
		{
			Id:           "acc1",
			Title:        "Счёт",
			StartBalance: 100000,
			Instrument:   2,
			Currency: model.Instrument{
				Title: "RUB",
				Rate:  1.0,
			},
		},
	}

	// Перевод между счетами в одной валюте
	transfer := model.Transaction{
		Id:      "1",
		Date:    "2023-05-10",
		Income:  10000,
		Outcome: 10000,
		InAccount: model.Account{
			Currency: model.Instrument{Rate: 1.0},
		},
		OutAccount: model.Account{
			Currency: model.Instrument{Rate: 1.0},
		},
		Deleted: false,
	}

	mockTransactionRepo.EXPECT().GetAll().Return([]model.Transaction{transfer}, nil)
	mockAccountRepo.EXPECT().GetAll().Return(accounts, nil)

	capital := NewCapital(mockTransactionRepo, mockAccountRepo)

	result, err := capital.GetCapital(2023)
	assert.NoError(t, err)
	assert.Len(t, result, 12)

	// Все месяцы должны иметь одинаковый баланс (перевод не меняет капитал)
	expectedBalance := 100000.0
	for _, m := range result {
		assert.Equal(t, expectedBalance, m.Balance)
	}
}

func TestGetCapital_MultiCurrency(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransactionRepo := transactionsRepo.NewMockRepositoryInterface(ctrl)
	mockAccountRepo := accountsRepo.NewMockRepositoryInterface(ctrl)

	accounts := []model.Account{
		{
			Id:           "acc1",
			Title:        "Счёт USD",
			StartBalance: 1000,
			Instrument:   1,
			Currency: model.Instrument{
				Title: "USD",
				Rate:  75.0,
			},
		},
	}

	// Доход в USD
	tx := model.Transaction{
		Id:      "1",
		Date:    "2023-01-15",
		Income:  500,
		Outcome: 0,
		InAccount: model.Account{
			Currency: model.Instrument{Rate: 75.0},
		},
		OutAccount: model.Account{
			Currency: model.Instrument{Rate: 1.0},
		},
		Deleted: false,
	}

	mockTransactionRepo.EXPECT().GetAll().Return([]model.Transaction{tx}, nil)
	mockAccountRepo.EXPECT().GetAll().Return(accounts, nil)

	capital := NewCapital(mockTransactionRepo, mockAccountRepo)

	result, err := capital.GetCapital(2023)
	assert.NoError(t, err)
	assert.Len(t, result, 12)

	// Начальный: 1000 * 75 = 75000
	// Январь: 75000 + 500 * 75 = 75000 + 37500 = 112500
	assert.Equal(t, 112500.0, result[0].Balance)
}

func TestGetCapital_DeletedTransactionsIgnored(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransactionRepo := transactionsRepo.NewMockRepositoryInterface(ctrl)
	mockAccountRepo := accountsRepo.NewMockRepositoryInterface(ctrl)

	accounts := []model.Account{
		{
			Id:           "acc1",
			Title:        "Счёт",
			StartBalance: 50000,
			Instrument:   2,
			Currency: model.Instrument{
				Title: "RUB",
				Rate:  1.0,
			},
		},
	}

	txValid := model.Transaction{
		Id:      "2",
		Date:    "2023-01-15",
		Income:  10000,
		Outcome: 0,
		InAccount: model.Account{
			Currency: model.Instrument{Rate: 1.0},
		},
		OutAccount: model.Account{
			Currency: model.Instrument{Rate: 1.0},
		},
		Deleted: false,
	}

	mockTransactionRepo.EXPECT().GetAll().Return([]model.Transaction{txValid}, nil)
	mockAccountRepo.EXPECT().GetAll().Return(accounts, nil)

	capital := NewCapital(mockTransactionRepo, mockAccountRepo)

	result, err := capital.GetCapital(2023)
	assert.NoError(t, err)
	assert.Len(t, result, 12)

	// Репозиторий уже отфильтровал deleted, поэтому только валидная транзакция учтена
	// 50000 + 10000 = 60000
	assert.Equal(t, 60000.0, result[0].Balance)
}

func TestGetCapital_InvalidDateIgnored(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransactionRepo := transactionsRepo.NewMockRepositoryInterface(ctrl)
	mockAccountRepo := accountsRepo.NewMockRepositoryInterface(ctrl)

	accounts := []model.Account{
		{
			Id:           "acc1",
			Title:        "Счёт",
			StartBalance: 100000,
			Instrument:   2,
			Currency: model.Instrument{
				Title: "RUB",
				Rate:  1.0,
			},
		},
	}

	invalidTx := model.Transaction{
		Id:      "1",
		Date:    "2023-02-30", // невалидная дата
		Income:  50000,
		Outcome: 0,
		InAccount: model.Account{
			Currency: model.Instrument{Rate: 1.0},
		},
		OutAccount: model.Account{
			Currency: model.Instrument{Rate: 1.0},
		},
		Deleted: false,
	}

	mockTransactionRepo.EXPECT().GetAll().Return([]model.Transaction{invalidTx}, nil)
	mockAccountRepo.EXPECT().GetAll().Return(accounts, nil)

	capital := NewCapital(mockTransactionRepo, mockAccountRepo)

	result, err := capital.GetCapital(2023)
	assert.NoError(t, err)
	assert.Len(t, result, 12)

	// Невалидная транзакция игнорируется, баланс не меняется
	for _, m := range result {
		assert.Equal(t, 100000.0, m.Balance)
	}
}

func TestGetCapital_CacheUsed(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransactionRepo := transactionsRepo.NewMockRepositoryInterface(ctrl)
	mockAccountRepo := accountsRepo.NewMockRepositoryInterface(ctrl)

	accounts := []model.Account{
		{
			Id:           "acc1",
			Title:        "Счёт",
			StartBalance: 50000,
			Instrument:   2,
			Currency: model.Instrument{
				Title: "RUB",
				Rate:  1.0,
			},
		},
	}

	tx := model.Transaction{
		Id:      "1",
		Date:    "2023-01-15",
		Income:  10000,
		Outcome: 0,
		InAccount: model.Account{
			Currency: model.Instrument{Rate: 1.0},
		},
		OutAccount: model.Account{
			Currency: model.Instrument{Rate: 1.0},
		},
		Deleted: false,
	}

	// Методы должны вызваться только один раз
	mockTransactionRepo.EXPECT().GetAll().Return([]model.Transaction{tx}, nil).Times(1)
	mockAccountRepo.EXPECT().GetAll().Return(accounts, nil).Times(1)

	capital := NewCapital(mockTransactionRepo, mockAccountRepo)

	// Первый вызов — загрузка данных
	result1, err := capital.GetCapital(2023)
	assert.NoError(t, err)
	assert.Len(t, result1, 12)

	// Второй вызов — должен использоваться кэш
	result2, err := capital.GetCapital(2023)
	assert.NoError(t, err)
	assert.Equal(t, result1, result2)
}
