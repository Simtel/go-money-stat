package app

import (
	"github.com/golang/mock/gomock"
	"money-stat/mock_app"
	"testing"
)

func TestContainer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mock_app.NewMockDbInterface(ctrl)
	//mockDB.EXPECT().GetGorm().Return(nil)
	container := NewContainer(mockDB)

	if container.GetDb() != mockDB {
		t.Errorf("Ожидалось, что метод GetDb вернет %v, но получили %v", mockDB, container.GetDb())
	}

	//mockGorm := mockDB.GetGorm()
	//mockTransactionRepository := transactions.NewRepository(mockGorm)
	/*if container.GetTransactionRepository() != mockTransactionRepository {
		t.Errorf("Ожидалось, что метод GetTransactionRepository вернет %v, но получили %v", mockTransactionRepository, container.GetTransactionRepository())
	}

	mockAccountRepository := accounts.NewRepository(mockGorm)
	if container.GetAccountRepository() != mockAccountRepository {
		t.Errorf("Ожидалось, что метод GetAccountRepository вернет %v, но получили %v", mockAccountRepository, container.GetAccountRepository())
	}*/
}
