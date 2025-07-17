package usecase

import (
	"fmt"
	"money-stat/internal/adapter/sqliterepo/zenrepo/accounts"
	transactionsRepo "money-stat/internal/adapter/sqliterepo/zenrepo/transactions"
	"money-stat/internal/model"
	"sort"
	"time"
)

const (
	dateLayout     = "2006-01-02"
	monthKeyLayout = "2006-01"
	baseMonth      = "2006-01"
)

type CapitalDto struct {
	Month   string  `json:"month"`
	Balance float64 `json:"balance"`
}

type Capital struct {
	transactionRepo transactionsRepo.RepositoryInterface
	accountRepo     accounts.RepositoryInterface
}

func NewCapital(transactionRepo transactionsRepo.RepositoryInterface, accountRepo accounts.RepositoryInterface) *Capital {
	return &Capital{
		transactionRepo: transactionRepo,
		accountRepo:     accountRepo,
	}
}

func (c *Capital) GetCapital() ([]CapitalDto, error) {
	monthlyStats := make(map[string]CapitalDto)

	if err := c.processTransactions(monthlyStats); err != nil {
		return nil, fmt.Errorf("failed to process transactions: %w", err)
	}

	return c.convertToSortedSlice(monthlyStats), nil
}

func (c *Capital) processTransactions(monthlyStats map[string]CapitalDto) error {
	transactions, err := c.transactionRepo.GetAll()
	if err != nil {
		return fmt.Errorf("failed to fetch transactions: %w", err)
	}

	for _, transaction := range transactions {
		monthKey, err := c.extractMonthKey(transaction.Date)
		if err != nil {
			return fmt.Errorf("failed to parse transaction date %s: %w", transaction.Date, err)
		}

		stat := c.getOrCreateMonthlyStat(monthlyStats, monthKey)
		updatedStat := c.applyTransactionToBalance(stat, transaction)
		monthlyStats[monthKey] = updatedStat
	}

	return nil
}

func (c *Capital) extractMonthKey(dateStr string) (string, error) {
	transactionTime, err := time.Parse(dateLayout, dateStr)
	if err != nil {
		return "", err
	}
	return transactionTime.Format(monthKeyLayout), nil
}

func (c *Capital) getOrCreateMonthlyStat(monthlyStats map[string]CapitalDto, monthKey string) CapitalDto {
	if stat, exists := monthlyStats[monthKey]; exists {
		return stat
	}
	return CapitalDto{Month: monthKey, Balance: 0}
}

func (c *Capital) applyTransactionToBalance(stat CapitalDto, transaction model.Transaction) CapitalDto {
	switch {
	case transaction.IsOutcome():
		stat.Balance -= c.convertToRubles(transaction.Outcome, transaction.OutAccount)
	case transaction.IsIncome():

		stat.Balance += c.convertToRubles(transaction.Income, transaction.InAccount)
	case transaction.IsTransfer():
		{
			diff := transaction.Income - transaction.Outcome
			stat.Balance += c.convertToRubles(diff, transaction.InAccount)
		}
	}
	return stat
}

func (c *Capital) convertToRubles(amount float64, account model.Account) float64 {
	if account.IsRuble() {
		return amount
	}
	return amount * account.Currency.Rate
}

func (c *Capital) convertToSortedSlice(monthlyStats map[string]CapitalDto) []CapitalDto {
	result := make([]CapitalDto, 0, len(monthlyStats))

	for _, stat := range monthlyStats {
		result = append(result, stat)
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Month < result[j].Month
	})

	return result
}
