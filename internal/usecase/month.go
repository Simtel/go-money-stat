package usecase

import (
	"fmt"
	"github.com/pterm/pterm"
	"log"
	transactionsRepo "money-stat/internal/adapter/sqliterepo/zenrepo/transactions"
	"money-stat/internal/model"
	"money-stat/internal/services/zenmoney"
	"strconv"
	"time"
)

type Month struct {
	api  *zenmoney.Api
	repo *transactionsRepo.Repository
}

func NewMonth(api *zenmoney.Api, repo *transactionsRepo.Repository) *Month {
	return &Month{api: api, repo: repo}
}

func (m *Month) GetMonthStat(month string) {

	log.Printf("Show %s months transactions", month)

	var transactions []model.Transaction
	if month == "current" {
		transactions = m.repo.GetCurrentMonth()
	}

	if month == "last" {
		transactions = m.repo.GetPreviousMonth()
	}

	var outComeSumm, inComeSumm float64

	var cnt int

	tableData := pterm.TableData{
		{"Дата", "Категория", "Сумма", "Счет", "Дата создания"},
		{" ", " ", " ", " ", " "},
	}

	for _, transaction := range transactions {
		cnt++

		var transactionTags string
		for _, tag := range transaction.Tag {
			transactionTags += tag.Title + " "
		}

		if transactionTags == "" {
			transactionTags = "Перевод"
		}

		var account string
		if transaction.IsIncome() {
			account = transaction.InAccount.Title
		}

		if transaction.IsOutcome() {
			account = transaction.OutAccount.Title
		}

		if transaction.IsTransfer() {
			account = transaction.OutAccount.Title + "->" + transaction.InAccount.Title
		}

		tCreatedDate := time.Unix(transaction.Created, 0)
		tableData = append(
			tableData,
			[]string{
				transaction.Date,
				transactionTags,
				transaction.FormatAmount(),
				account,
				tCreatedDate.Format("2006-01-02 15:04:05"),
			},
		)
		if transaction.Outcome > 0 && transaction.Income == 0 {
			outComeSumm = outComeSumm + transaction.Outcome
		}

		if transaction.Income > 0 && transaction.Outcome == 0 {
			inComeSumm = inComeSumm + transaction.Income
		}
	}

	errTable := pterm.DefaultTable.WithHasHeader().WithBoxed().WithRowSeparator("-").WithData(tableData).Render()
	if errTable != nil {
		fmt.Println(errTable)
	}

	summData := pterm.TableData{
		{
			"Транзакций",
			"Доходов в рублях",
			"Расходов в рублях",
			"Чистыми",
		},
		{" ", " ", ""},
		{
			strconv.Itoa(cnt),
			strconv.FormatFloat(inComeSumm, 'f', 2, 64),
			strconv.FormatFloat(outComeSumm, 'f', 2, 64),
			strconv.FormatFloat(inComeSumm-outComeSumm, 'f', 2, 64),
		},
	}

	errSummTable := pterm.DefaultTable.WithHasHeader().WithBoxed().WithRowSeparator("-").WithData(summData).Render()
	if errSummTable != nil {
		fmt.Println(errSummTable)
	}

}
