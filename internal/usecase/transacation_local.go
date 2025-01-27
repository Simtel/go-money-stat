package usecase

import (
	"fmt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"money-stat/internal/model"
	"strconv"
)

type TransactionsLocal struct {
}

func (t *TransactionsLocal) GetLast(cnt int) {

	fmt.Println("Поиск локальных транзакций")
	db, err := gorm.Open(sqlite.Open("zenmoney.db?cache=shared&mode=rwc"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	var transactions []model.Transaction
	result := db.Limit(cnt).Find(&transactions)

	if result.RowsAffected == 0 {
		fmt.Println("Нет локальных записей")
	}

	fmt.Println("Вернулось записей:" + strconv.Itoa(len(transactions)))

	for _, transaction := range transactions {
		fmt.Println(transaction.Id)
	}
}
