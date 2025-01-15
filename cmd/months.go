package cmd

import (
	"fmt"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"log"
	"math"
	"money-stat/internal/services/zenmoney"
	"net/http"
	"sort"
	"strconv"
	"time"
)

func RunMonths() *cobra.Command {
	cmd := &cobra.Command{
		Use:       "months",
		Short:     "Показать транзакции за месяц (текущий, прошлый)",
		ValidArgs: []string{"current", "last"},
		Args:      cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	}

	cmd.RunE = func(cmd *cobra.Command, args []string) error {

		multi := pterm.DefaultMultiPrinter

		loadSpinner, _ := pterm.DefaultSpinner.WithWriter(multi.NewWriter()).Start("Загрузка данных")
		_, err := multi.Start()
		if err != nil {
			return err
		}
		month := args[0]
		log.Printf("Show %s months transactions", month)

		api := zenmoney.NewApi(&http.Client{})

		diff, err := api.Diff()
		if err != nil {
			log.Fatal(err)
		}

		loadSpinner.Success("Загрузка завершена!")
		_, errStop := multi.Stop()
		if errStop != nil {
			return err
		}
		now := time.Now()
		var timestamp int64
		if month == "current" {

			firstDayOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())

			timestamp = firstDayOfMonth.Unix()
		}

		var outComeSumm, inComeSumm, diffAmount float64

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
			if tTime.Unix() < timestamp || t.IsDeleted() {
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

			if transaction.Outcome > 0 && transaction.Income > 0 {
				diffAmount = diffAmount + math.Abs(transaction.Outcome-transaction.Income)
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

		return nil
	}

	return cmd
}
