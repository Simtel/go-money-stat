package usecase

import (
	"money-stat/internal/adapter/sqliterepo/zenrepo/accounts/mocks"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	accountsRepo "money-stat/internal/adapter/sqliterepo/zenrepo/accounts/mocks"
	transactionsRepo "money-stat/internal/adapter/sqliterepo/zenrepo/transactions/mocks"
	"money-stat/internal/model"
)

func TestGetCapital_EmptyTransactions(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransactionRepo := transactionsRepo.NewMockRepositoryInterface(ctrl)
	mockAccountRepo := accountsRepo.NewMockRepositoryInterface(ctrl)

	mockTransactionRepo.EXPECT().GetAll().Return([]model.Transaction{}, nil)

	capital := NewCapital(mockTransactionRepo, mockAccountRepo)

	result, err := capital.GetCapital()
	assert.NoError(t, err)
	assert.Empty(t, result)
}

func TestGetCapital_InvalidDate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransactionRepo := transactionsRepo.NewMockRepositoryInterface(ctrl)
	mockAccountRepo := mocks.NewMockRepositoryInterface(ctrl)

	invalidTx := model.Transaction{
		Id:      "1",
		Date:    "2023-02-30", // invalid date
		Income:  100.0,
		Outcome: 50.0,
		InAccount: model.Account{
			Currency: model.Instrument{
				Title: "RUB",
				Rate:  1.0,
			},
		},
		OutAccount: model.Account{
			Currency: model.Instrument{
				Title: "USD",
				Rate:  75.0,
			},
		},
		Deleted: false,
	}

	mockTransactionRepo.EXPECT().GetAll().Return([]model.Transaction{invalidTx}, nil)

	capital := NewCapital(mockTransactionRepo, mockAccountRepo)

	result, err := capital.GetCapital()
	assert.NoError(t, err)
	assert.Empty(t, result)
}

func TestGetCapital_ValidTransactions(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransactionRepo := transactionsRepo.NewMockRepositoryInterface(ctrl)
	mockAccountRepo := mocks.NewMockRepositoryInterface(ctrl)

	tx1 := model.Transaction{
		Id:      "1",
		Date:    "2023-01-15",
		Income:  100.0,
		Outcome: 50.0,
		InAccount: model.Account{
			Currency: model.Instrument{
				Title: "RUB",
				Rate:  1.0,
			},
		},
		OutAccount: model.Account{
			Currency: model.Instrument{
				Title: "USD",
				Rate:  75.0,
			},
		},
		Deleted: false,
	}

	tx2 := model.Transaction{
		Id:      2,
		Date:    "2023-02-15",
		Income:  200.0,
		Outcome: 100.0,
		InAccount: model.Account{
			Currency: model.Instrument{
				Title: "RUB",
				Rate:  1.0,
			},
		},
		OutAccount: model.Account{
			Currency: model.Instrument{
				Title: "USD",
				Rate:  75.0,
			},
		},
		Deleted: false,
	}

	mockTransactionRepo.EXPECT().GetAll().Return([]model.Transaction{tx1, tx2}, nil)

	capital := NewCapital(mockTransactionRepo, mockAccountRepo)

	expectedResult := []MonthlyBalance{
		{Month: "2023-01", Balance: 50.0},
		{Month: "2023-02", Balance: 150.0},
	}

	result, err := capital.GetCapital()
	assert.NoError(t, err)
	assert.Equal(t, expectedResult, result)
}

func TestGetCapital_MultipleMonths(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransactionRepo := transactionsRepo.NewMockRepositoryInterface(ctrl)
	mockAccountRepo := mocks.NewMockRepositoryInterface(ctrl)

	tx1 := model.Transaction{
		Id:      1,
		Date:    "2023-01-15",
		Income:  100.0,
		Outcome: 50.0,
		InAccount: model.Account{
			Currency: model.Instrument{
				Title: "RUB",
				Rate:  1.0,
			},
		},
		OutAccount: model.Account{
			Currency: model.Instrument{
				Title: "USD",
				Rate:  75.0,
			},
		},
		Deleted: false,
	}

	tx2 := model.Transaction{
		Id:      2,
		Date:    "2023-02-15",
		Income:  200.0,
		Outcome: 100.0,
		InAccount: model.Account{
			Currency: model.Instrument{
				Title: "RUB",
				Rate:  1.0,
			},
		},
		OutAccount: model.Account{
			Currency: model.Instrument{
				Title: "USD",
				Rate:  75.0,
			},
		},
		Deleted: false,
	}

	tx3 := model.Transaction{
		Id:      3,
		Date:    "2023-03-15",
		Income:  300.0,
		Outcome: 150.0,
		InAccount: model.Account{
			Currency: model.Instrument{
				Title: "RUB",
				Rate:  1.0,
			},
		},
		OutAccount: model.Account{
			Currency: model.Instrument{
				Title: "USD",
				Rate:  75.0,
			},
		},
		Deleted: false,
	}

	mockTransactionRepo.EXPECT().GetAll().Return([]model.Transaction{tx1, tx2, tx3}, nil)

	capital := NewCapital(mockTransactionRepo, mockAccountRepo)

	expectedResult := []MonthlyBalance{
		{Month: "2023-01", Balance: 50.0},
		{Month: "2023-02", Balance: 150.0},
		{Month: "2023-03", Balance: 300.0},
	}

	result, err := capital.GetCapital()
	assert.NoError(t, err)
	assert.Equal(t, expectedResult, result)
}

func TestGetCapital_SingleMonth(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransactionRepo := transactionsRepo.NewMockRepositoryInterface(ctrl)
	mockAccountRepo := mocks.NewMockRepositoryInterface(ctrl)

	tx1 := model.Transaction{
		Id:      1,
		Date:    "2023-01-15",
		Income:  100.0,
		Outcome: 50.0,
		InAccount: model.Account{
			Currency: model.Instrument{
				Title: "RUB",
				Rate:  1.0,
			},
		},
		OutAccount: model.Account{
			Currency: model.Instrument{
				Title: "USD",
				Rate:  75.0,
			},
		},
		Deleted: false,
	}

	mockTransactionRepo.EXPECT().GetAll().Return([]model.Transaction{tx1}, nil)

	capital := NewCapital(mockTransactionRepo, mockAccountRepo)

	expectedResult := []MonthlyBalance{
		{Month: "2023-01", Balance: 50.0},
	}

	result, err := capital.GetCapital()
	assert.NoError(t, err)
	assert.Equal(t, expectedResult, result)
}

func TestGetCapital_NegativeBalance(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransactionRepo := transactionsRepo.NewMockRepositoryInterface(ctrl)
	mockAccountRepo := mocks.NewMockRepositoryInterface(ctrl)

	tx1 := model.Transaction{
		Id:      1,
		Date:    "2023-01-15",
		Income:  50.0,
		Outcome: 100.0,
		InAccount: model.Account{
			Currency: model.Instrument{
				Title: "RUB",
				Rate:  1.0,
			},
		},
		OutAccount: model.Account{
			Currency: model.Instrument{
				Title: "USD",
				Rate:  75.0,
			},
		},
		Deleted: false,
	}

	mockTransactionRepo.EXPECT().GetAll().Return([]model.Transaction{tx1}, nil)

	capital := NewCapital(mockTransactionRepo, mockAccountRepo)

	expectedResult := []MonthlyBalance{
		{Month: "2023-01", Balance: -50.0},
	}

	result, err := capital.GetCapital()
	assert.NoError(t, err)
	assert.Equal(t, expectedResult, result)
}
