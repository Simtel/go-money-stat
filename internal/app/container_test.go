package app

import (
	"github.com/golang/mock/gomock"
	"money-stat/internal/adapter/sqliterepo/zenrepo/accounts"
	"money-stat/internal/adapter/sqliterepo/zenrepo/transactions"
	"money-stat/internal/app/mocks"
	"reflect"
	"testing"
)

func TestContainer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mocks.NewMockDbInterface(ctrl)
	mockDB.EXPECT().GetGorm().Return(nil).AnyTimes()
	container := NewContainer(mockDB)

	if container.GetDb() != mockDB {
		t.Errorf("Ожидалось, что метод GetDb вернет %v, но получили %v", mockDB, container.GetDb())
	}

	mockGorm := mockDB.GetGorm()
	mockTransactionRepository := transactions.NewRepository(mockGorm)
	if reflect.TypeOf(container.GetTransactionRepository()) != reflect.TypeOf(mockTransactionRepository) {
		t.Errorf("Ожидалось, что метод GetTransactionRepository вернет %v, но получили %v", mockTransactionRepository, container.GetTransactionRepository())
	}

	mockAccountRepository := accounts.NewRepository(mockGorm)
	if reflect.TypeOf(container.GetAccountRepository()) != reflect.TypeOf(mockAccountRepository) {
		t.Errorf("Ожидалось, что метод GetAccountRepository вернет %v, но получили %v", mockAccountRepository, container.GetAccountRepository())
	}
}
