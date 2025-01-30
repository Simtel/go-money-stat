package usecase

import (
	"fmt"
	"github.com/pterm/pterm"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"money-stat/internal/model"
	"strconv"
	"time"
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
	result := db.Limit(cnt).Order("Date desc").Find(&transactions)

	if result.RowsAffected == 0 {
		fmt.Println("Нет локальных записей")
	}

	fmt.Println("Вернулось записей:" + strconv.Itoa(len(transactions)))

	tableData := pterm.TableData{
		{"Дата", "Сумма", "Дата создания"},
		{" ", " ", " "},
	}

	for _, transaction := range transactions {
		tCreatedDate := time.Unix(transaction.Created, 0)
		tableData = append(
			tableData,
			[]string{
				transaction.Date,
				transaction.FormatAmount(),
				tCreatedDate.Format("2006-01-02 15:04:05"),
			},
		)
	}

	errTable := pterm.DefaultTable.WithHasHeader().WithBoxed().WithRowSeparator("-").WithData(tableData).Render()
	if errTable != nil {
		fmt.Println(errTable)
	}
}
