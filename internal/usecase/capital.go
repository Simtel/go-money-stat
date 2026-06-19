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

	accountRepo accounts.RepositoryInterface

	cachedResult []MonthlyBalance

	cacheTime time.Time
	cacheTTL  time.Duration
}

type MonthlyBalance struct {
	Month   string
	Balance float64
}

func NewCapital(transactionRepo transactionsRepo.RepositoryInterface, accountRepo accounts.RepositoryInterface) *Capital {
	return &Capital{
		transactionRepo: transactionRepo,
		accountRepo:     accountRepo,
		cacheTTL:        5 * time.Minute,
	}
}

func (c *Capital) GetCapital(year int) ([]MonthlyBalance, error) {
	// Проверяем актуальность кэша
	if c.cachedResult != nil && time.Since(c.cacheTime) < c.cacheTTL {
		return c.cachedResult, nil
	}

	transactions, err := c.transactionRepo.GetAll()
	if err != nil {
		return nil, fmt.Errorf("получение транзакций: %w", err)
	}

	accounts := c.accountRepo.GetAll()

	result := c.calculateMonthlyBalances(transactions, accounts, year)

	// Обновляем кэш
	c.cachedResult = make([]MonthlyBalance, len(result))
	copy(c.cachedResult, result)
	c.cacheTime = time.Now()

	return result, nil
}

func (c *Capital) calculateMonthlyBalances(transactions []model.Transaction, accounts []model.Account, year int) []MonthlyBalance {
	// Вычисляем начальный баланс всех счетов (в базовой валюте)
	startCapital := 0.0
	for _, acc := range accounts {
		rate := 1.0
		if acc.Currency.Rate > 0 {
			rate = acc.Currency.Rate
		}
		startCapital += acc.StartBalance * rate
	}

	// Фильтруем удалённые транзакции
	var validTx []model.Transaction
	for _, tx := range transactions {
		if !tx.Deleted {
			validTx = append(validTx, tx)
		}
	}

	// Сортируем транзакции по дате
	sort.Slice(validTx, func(i, j int) bool {
		dateI, errI := time.Parse(dateLayout, validTx[i].Date)
		dateJ, errJ := time.Parse(dateLayout, validTx[j].Date)
		if errI != nil || errJ != nil {
			return false
		}
		return dateI.Before(dateJ)
	})

	// Группируем изменения по месяцам
	monthlyChanges := make(map[string]float64)
	for _, tx := range validTx {
		txDate, err := time.Parse(dateLayout, tx.Date)
		if err != nil {
			continue
		}

		monthKey := txDate.Format(baseMonth)

		incomeRate := 1.0
		outcomeRate := 1.0

		if tx.InAccount.Currency.Rate > 0 {
			incomeRate = tx.InAccount.Currency.Rate
		}
		if tx.OutAccount.Currency.Rate > 0 {
			outcomeRate = tx.OutAccount.Currency.Rate
		}

		change := tx.Income*incomeRate - tx.Outcome*outcomeRate
		monthlyChanges[monthKey] += change
	}

	// Вычисляем капитал на начало запрошенного года:
	// startCapital + все изменения до начала года
	capitalAtYearStart := startCapital
	yearStart := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)

	// Сортируем ключи месяцев для накопления
	var monthKeys []string
	for k := range monthlyChanges {
		monthKeys = append(monthKeys, k)
	}
	sort.Strings(monthKeys)

	for _, monthKey := range monthKeys {
		monthDate, err := time.Parse(baseMonth, monthKey)
		if err != nil {
			continue
		}
		if monthDate.Before(yearStart) {
			capitalAtYearStart += monthlyChanges[monthKey]
		}
	}

	// Формируем результат — все 12 месяцев года
	currentBalance := capitalAtYearStart
	var result []MonthlyBalance

	for m := 1; m <= 12; m++ {
		monthDate := time.Date(year, time.Month(m), 1, 0, 0, 0, 0, time.UTC)
		monthKey := monthDate.Format(baseMonth)

		currentBalance += monthlyChanges[monthKey]

		result = append(result, MonthlyBalance{
			Month:   monthKey,
			Balance: currentBalance,
		})
	}

	return result
}
