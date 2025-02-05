package usecase

import (
	"fmt"
	"github.com/pterm/pterm"
	"log"
	"money-stat/internal/adapter/sqliterepo/zenrepo/transactions"
	"money-stat/internal/services/zenmoney"
	"net/http"
	"sort"
	"strconv"
	"time"
)

type Month struct {
	api  *zenmoney.Api
	repo *transactions.Repository
}

func NewMonth(api *zenmoney.Api, repo *transactions.Repository) *Month {
	return &Month{api: api, repo: repo}
}

func (m *Month) GetMonthStat(month string) {
	multi := pterm.DefaultMultiPrinter

	loadSpinner, _ := pterm.DefaultSpinner.WithWriter(multi.NewWriter()).Start("Загрузка данных")
	_, err := multi.Start()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Show %s months transactions", month)

	m.repo.GetCurrentMonths()

	api := zenmoney.NewApi(&http.Client{})

	diff, err := api.Diff()
	if err != nil {
		log.Fatal(err)
	}

	loadSpinner.Success("Загрузка завершена!")
	_, errStop := multi.Stop()
	if errStop != nil {
		log.Fatal(err)
	}
	now := time.Now()
	var firstDayTimestamp, lastDayTimestamp int64
	if month == "current" {

		firstDayOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
		firstDayTimestamp = firstDayOfMonth.Unix()
		firstOfNextMonth := time.Date(now.Year(), now.Month()+1, 1, 0, 0, 0, 0, now.Location())
		lastOfCurrentMonth := firstOfNextMonth.AddDate(0, 0, -1)
		lastDayTimestamp = lastOfCurrentMonth.Unix()
	}

	if month == "last" {

		firstDayOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
		previousMonth := firstDayOfMonth.AddDate(0, -1, 0)
		firstDayTimestamp = previousMonth.Unix()

		firstOfNextMonth := time.Date(previousMonth.Year(), previousMonth.Month()+1, 1, 23, 59, 59, 0, now.Location())
		lastDayMonth := firstOfNextMonth.AddDate(0, 0, -1)
		lastDayTimestamp = lastDayMonth.Unix()
	}

	var outComeSumm, inComeSumm float64

	var cnt int

	tags := diff.GetIndexedTags()
	accounts := diff.GetIndexedAccounts()

	tableData := pterm.TableData{
		{"Дата", "Категория", "Сумма", "Счет", "Дата создания"},
		{" ", " ", " ", " ", " "},
	}

	var transactions []zenmoney.Transaction

	for _, t := range diff.Transaction {
		layout := "2006-01-02"
		tTime, _ := time.Parse(layout, t.Date)
		if tTime.Unix() < firstDayTimestamp || tTime.Unix() > lastDayTimestamp || t.IsDeleted() {
			continue
		}

		transactions = append(transactions, t)
	}

	sort.Slice(transactions, func(i, j int) bool {
		return transactions[i].Created > transactions[j].Created
	})

	for _, transaction := range transactions {
		cnt++

		var transactionTags string
		for _, tag := range transaction.Tag {
			transactionTags += tags[tag].Title + " "
		}

		if transactionTags == "" {
			transactionTags = "Перевод"
		}

		var account string
		if transaction.IsIncome() {
			account = accounts[transaction.IncomeAccount].Title
		}

		if transaction.IsOutcome() {
			account = accounts[transaction.OutcomeAccount].Title
		}

		if transaction.IsTransfer() {
			account = accounts[transaction.OutcomeAccount].Title + "->" + accounts[transaction.IncomeAccount].Title
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
