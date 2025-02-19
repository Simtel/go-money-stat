package zenmoney

import (
	"strconv"
)

type Transaction struct {
	Id                string   `json:"id"`
	Changed           int64    `json:"changed"`
	Created           int64    `json:"created"`
	IncomeInstrument  int64    `json:"incomeInstrument"`
	Income            float64  `json:"income"`
	OutcomeInstrument int64    `json:"outcomeInstrument"`
	Outcome           float64  `json:"outcome"`
	Date              string   `json:"date"`
	Tag               []string `json:"tag"`
	Deleted           bool     `json:"deleted"`
	IncomeAccount     string   `json:"incomeAccount"`
	OutcomeAccount    string   `json:"outcomeAccount"`
	Comment           string   `json:"comment"`
}

func (t Transaction) FormatAmount() string {
	if t.Income == 0 && t.Outcome > 0 {
		return strconv.FormatFloat(-t.Outcome, 'f', 2, 64)
	}

	if t.Income > 0 && t.Outcome == 0 {
		return strconv.FormatFloat(t.Income, 'f', 2, 64)
	}

	if t.Income > 0 && t.Outcome > 0 {
		return strconv.FormatFloat(t.Outcome, 'f', 2, 64) + " -> " + strconv.FormatFloat(t.Income, 'f', 2, 64)
	}

	return "0"
}

func (t Transaction) IsDeleted() bool {
	return t.Deleted
}

func (t Transaction) IsIncome() bool {
	return t.Income > 0 && t.Outcome == 0
}

func (t Transaction) IsOutcome() bool {
	return t.Outcome > 0 && t.Income == 0
}

func (t Transaction) IsTransfer() bool {
	return t.Income > 0 && t.Outcome > 0
}
