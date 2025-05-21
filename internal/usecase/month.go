package usecase

import (
	"log"
	transactionsRepo "money-stat/internal/adapter/sqliterepo/zenrepo/transactions"
	"money-stat/internal/model"
	"time"
)

type MonthStatDto struct {
	Transactions []MonthStatTransactionDto
	OutComeSumm  float64
	InComeSumm   float64
	Count        int
}

type MonthStatTransactionDto struct {
	Date         string
	Tags         string
	FormatAmount string
	Account      string
	CreatedAt    string
}

type Month struct {
	repo transactionsRepo.RepositoryInterface
}

func NewMonth(repo transactionsRepo.RepositoryInterface) *Month {
	return &Month{repo: repo}
}

func (m *Month) GetMonthStat(month string) MonthStatDto {

	var monthStat MonthStatDto
	log.Printf("Show %s months transactions", month)

	var transactions []model.Transaction
	if month == "current" {
		transactions = m.repo.GetCurrentMonth()
	}

	if month == "previous" {
		transactions = m.repo.GetPreviousMonth()
	}

	var outComeSumm, inComeSumm float64

	var cnt int

	for _, transaction := range transactions {
		cnt++

		tCreatedDate := time.Unix(transaction.Created, 0)
		monthStat.Transactions = append(
			monthStat.Transactions,
			MonthStatTransactionDto{
				transaction.Date,
				transaction.GetTagsTitle(),
				transaction.FormatAmount(),
				m.getAccountTitle(transaction),
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
	monthStat.OutComeSumm = outComeSumm
	monthStat.InComeSumm = inComeSumm
	monthStat.Count = cnt

	return monthStat

}

func (m *Month) getAccountTitle(transaction model.Transaction) string {
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
	return account
}
