package usecase

import (
	"github.com/pterm/pterm"
	"log"
	"money-stat/internal/services/zenmoney"
	"strconv"
)

type Accounts struct {
	api *zenmoney.Api
}

func NewAccounts(api *zenmoney.Api) *Accounts {
	return &Accounts{api: api}
}

func (a *Accounts) GetAccounts() {
	result, err := a.api.Diff()

	if err != nil {
		log.Println(err)
	}

	tableData := pterm.TableData{
		{"Счет", "Баланс", "Валюта"},
		{" ", " ", " "},
	}

	instruments := make(map[int]string)

	rateDollar := 0.0
	for _, instrument := range result.Instrument {
		instruments[instrument.Id] = instrument.Symbol
		if instrument.IsDollar() {
			rateDollar = instrument.Rate
		}
	}

	var summRuble float64
	var summDollar float64

	for _, account := range result.Account {
		tableData = append(tableData, []string{account.Title, strconv.FormatFloat(account.Balance, 'f', 2, 64), instruments[account.Instrument]})
		if account.IsRuble() {
			summRuble = summRuble + account.Balance
		}
		if account.IsDollar() {
			summDollar = summDollar + account.Balance
		}
	}

	pterm.DefaultTable.WithHasHeader().WithBoxed().WithRowSeparator("-").WithData(tableData).Render()

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

	pterm.DefaultTable.WithHasHeader().WithBoxed().WithRowSeparator("-").WithData(summData).Render()
}
