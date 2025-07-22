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
	dateLayout = "2006-01-02"
	baseMonth  = "2006-01"
)

type Capital struct {
	transactionRepo transactionsRepo.RepositoryInterface
	accountRepo     accounts.RepositoryInterface
}

type MonthlyBalance struct {
	Month   string
	Balance float64
}

func NewCapital(transactionRepo transactionsRepo.RepositoryInterface, accountRepo accounts.RepositoryInterface) *Capital {
	return &Capital{
		transactionRepo: transactionRepo,
		accountRepo:     accountRepo,
	}
}

func (c *Capital) GetCapital() ([]MonthlyBalance, error) {

	transactions, _ := c.transactionRepo.GetAll()

	return c.calculateMonthlyBalances(transactions), nil

}

func (c *Capital) calculateMonthlyBalances(transactions []model.Transaction) []MonthlyBalance {
	if len(transactions) == 0 {
		return []MonthlyBalance{}
	}

	validTx := c.getValidTransactions(transactions)

	if len(validTx) == 0 {
		return []MonthlyBalance{}
	}

	sort.Slice(validTx, func(i, j int) bool {
		dateI, _ := time.Parse(dateLayout, validTx[i].Date)
		dateJ, _ := time.Parse(dateLayout, validTx[j].Date)
		return dateI.Before(dateJ)
	})

	firstDate, err := time.Parse(dateLayout, validTx[0].Date)
	if err != nil {
		fmt.Printf("Ошибка парсинга даты: %v\n", err)
		return []MonthlyBalance{}
	}

	monthlyData := c.getMonthlyBalance(validTx)

	lastDate, _ := time.Parse(dateLayout, validTx[len(validTx)-1].Date)

	var result []MonthlyBalance
	currentBalance := 0.0

	current := time.Date(firstDate.Year(), firstDate.Month(), 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(lastDate.Year(), lastDate.Month(), 1, 0, 0, 0, 0, time.UTC)

	for current.Before(end.AddDate(0, 1, 0)) {
		monthKey := current.Format(baseMonth)
		monthChange := monthlyData[monthKey]
		currentBalance += monthChange

		result = append(result, MonthlyBalance{
			Month:   monthKey,
			Balance: currentBalance,
		})

		current = current.AddDate(0, 1, 0)
	}

	return result
}

func (c *Capital) convertToRubles(amount float64, account model.Account) float64 {
	if account.IsRuble() {
		return amount
	}
	return amount * account.Currency.Rate
}

func (c *Capital) getValidTransactions(transactions []model.Transaction) []model.Transaction {
	var validTx []model.Transaction
	for _, tx := range transactions {
		if !tx.Deleted {
			validTx = append(validTx, tx)
		}
	}

	return validTx
}

func (c *Capital) getMonthlyBalance(transactions []model.Transaction) map[string]float64 {
	monthlyData := make(map[string]float64)

	for _, tx := range transactions {
		txDate, err := time.Parse(dateLayout, tx.Date)
		if err != nil {
			continue
		}

		monthKey := txDate.Format(baseMonth)

		monthlyData[monthKey] += c.convertToRubles(tx.Income, tx.InAccount) - c.convertToRubles(tx.Outcome, tx.OutAccount)
	}

	return monthlyData
}
