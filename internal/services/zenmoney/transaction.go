package zenmoney

import "strconv"

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
}

func (t Transaction) FormatAmount() string {
	if t.Income == 0 && t.Outcome > 0 {
		return strconv.FormatFloat(-t.Outcome, 'f', 2, 64)
	}

	if t.Income > 0 && t.Outcome == 0 {
		return strconv.FormatFloat(-t.Income, 'f', 2, 64)
	}

	if t.Income > 0 && t.Outcome > 0 {
		return strconv.FormatFloat(t.Outcome, 'f', 2, 64) + " -> " + strconv.FormatFloat(t.Income, 'f', 2, 64)
	}

	return "0"
}
