package usecase

import (
	"fmt"
	"github.com/pterm/pterm"
	"money-stat/internal/adapter/sqliterepo/zenrepo/accounts"
	"strconv"
)

type Accounts struct {
	repo *accounts.Repository
}

func NewAccounts(repo *accounts.Repository) *Accounts {
	return &Accounts{repo: repo}
}

func (a *Accounts) GetAccounts() {

	accountsList := a.repo.GetAll()

	tableData := pterm.TableData{
		{"Счет", "Баланс", "Валюта"},
		{" ", " ", " "},
	}

	rateDollar := 0.0

	var summRuble float64
	var summDollar float64

	for _, account := range accountsList {
		tableData = append(
			tableData,
			[]string{
				account.Title,
				strconv.FormatFloat(account.Balance, 'f', 2, 64),
				account.Currency.ShortTitle,
			},
		)
		if account.IsRuble() {
			summRuble = summRuble + account.Balance
		}
		if account.IsDollar() {
			summDollar = summDollar + account.Balance
			rateDollar = account.Currency.Rate
		}
	}

	errTable := pterm.DefaultTable.WithHasHeader().WithBoxed().WithRowSeparator("-").WithData(tableData).Render()
	if errTable != nil {
		fmt.Println(errTable)
	}

	summData := pterm.TableData{
		{
			"Итого в рублях",
			"Итого в долларах",
			"Общая сумма в рублях",
		},
		{" ", " "},
		{
			strconv.FormatFloat(summRuble, 'f', 2, 64),
			strconv.FormatFloat(summDollar, 'f', 2, 64),
			strconv.FormatFloat(summRuble+(summDollar*rateDollar), 'f', 2, 64),
		},
	}

	errSummTable := pterm.DefaultTable.WithHasHeader().WithBoxed().WithRowSeparator("-").WithData(summData).Render()
	if errSummTable != nil {
		fmt.Println(errSummTable)
	}
}
