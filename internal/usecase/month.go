package usecase

import (
	"fmt"
	"log"
	"sort"
	"time"

	transactionsRepo "money-stat/internal/adapter/sqliterepo/zenrepo/transactions"
	"money-stat/internal/model"
)

type MonthStatDto struct {
	Transactions []MonthStatTransactionDto `json:"transactions"`
	OutcomeSumm  float64                   `json:"outcome_summ"`
	IncomeSumm   float64                   `json:"income_summ"`
	Count        int                       `json:"count"`
}

type MonthStatTransactionDto struct {
	Date         string `json:"date"`
	Tags         string `json:"tags"`
	FormatAmount string `json:"format_amount"`
	Account      string `json:"account"`
	CreatedAt    string `json:"created_at"`
	Comment      string `json:"comment"`
}

type MonthType string

const (
	CurrentMonth  MonthType = "current"
	PreviousMonth MonthType = "previous"
)

type Month struct {
	repo transactionsRepo.RepositoryInterface
}

func NewMonth(repo transactionsRepo.RepositoryInterface) *Month {
	return &Month{repo: repo}
}

func (m *Month) GetMonthStat(month string) (MonthStatDto, error) {
	monthType := MonthType(month)
	log.Printf("Generating statistics for %s month", month)

	transactions, err := m.getTransactionsByMonth(monthType)
	if err != nil {
		return MonthStatDto{}, fmt.Errorf("failed to get transactions: %w", err)
	}

	return m.buildMonthStat(transactions), nil
}

func (m *Month) getTransactionsByMonth(monthType MonthType) ([]model.Transaction, error) {
	switch monthType {
	case CurrentMonth:
		return m.repo.GetCurrentMonth(), nil
	case PreviousMonth:
		return m.repo.GetPreviousMonth(), nil
	default:
		return nil, fmt.Errorf("unsupported month type: %s", monthType)
	}
}

func (m *Month) buildMonthStat(transactions []model.Transaction) MonthStatDto {
	monthStat := MonthStatDto{
		Transactions: make([]MonthStatTransactionDto, 0, len(transactions)),
	}

	var outcomeSumm, incomeSumm float64

	for _, transaction := range transactions {
		transactionDto := m.convertTransactionToDto(transaction)
		monthStat.Transactions = append(monthStat.Transactions, transactionDto)

		outcomeSumm += transaction.Outcome
		incomeSumm += transaction.Income
	}

	monthStat.OutcomeSumm = outcomeSumm
	monthStat.IncomeSumm = incomeSumm
	monthStat.Count = len(transactions)

	m.sortTransactionsByCreatedAt(monthStat.Transactions)

	return monthStat
}

func (m *Month) convertTransactionToDto(transaction model.Transaction) MonthStatTransactionDto {
	createdDate := time.Unix(transaction.Created, 0)

	return MonthStatTransactionDto{
		Date:         transaction.Date,
		Tags:         transaction.GetTagsTitle(),
		FormatAmount: transaction.FormatAmount(),
		Account:      m.getAccountTitle(transaction),
		CreatedAt:    createdDate.Format("2006-01-02 15:04:05"),
		Comment:      transaction.Comment,
	}
}

func (m *Month) sortTransactionsByCreatedAt(transactions []MonthStatTransactionDto) {
	sort.Slice(transactions, func(i, j int) bool {
		return transactions[i].CreatedAt < transactions[j].CreatedAt
	})
}

func (m *Month) getAccountTitle(transaction model.Transaction) string {
	switch {
	case transaction.IsTransfer():
		return fmt.Sprintf("%s->%s", transaction.OutAccount.Title, transaction.InAccount.Title)
	case transaction.IsIncome():
		return transaction.InAccount.Title
	case transaction.IsOutcome():
		return transaction.OutAccount.Title
	default:
		return ""
	}
}
