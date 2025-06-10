package usecase

import (
	"money-stat/internal/adapter/sqliterepo/zenrepo/accounts"
	"money-stat/internal/model"
	"strconv"
)

type Accounts struct {
	repo accounts.RepositoryInterface
}

type AccountDto struct {
	Account  string
	Balance  string
	Currency string
}

type AccountStatDto struct {
	Accounts   []AccountDto
	RateDollar float64
	SummRuble  float64
	SummDollar float64
}

func NewAccounts(repo accounts.RepositoryInterface) *Accounts {
	return &Accounts{repo: repo}
}

func (a *Accounts) GetAccounts() AccountStatDto {

	accountsList := a.repo.GetAll()

	var statDto AccountStatDto

	statDto.RateDollar = 0.0

	for _, account := range accountsList {
		statDto.Accounts = append(
			statDto.Accounts,
			AccountDto{
				account.Title,
				strconv.FormatFloat(account.Balance, 'f', 2, 64),
				account.Currency.ShortTitle,
			},
		)
		statDto.updateSummByAccountType(account)
	}

	return statDto

}
func (a *AccountStatDto) updateSummByAccountType(account model.Account) {
	if account.IsRuble() {
		a.SummRuble = a.SummRuble + account.Balance
	}
	if account.IsDollar() {
		a.SummDollar = a.SummDollar + account.Balance
		a.RateDollar = account.Currency.Rate
	}
}
