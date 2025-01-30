package model

import "strconv"

type Transaction struct {
	Id                string
	Changed           int64
	Created           int64
	IncomeInstrument  int64
	Income            float64
	OutcomeInstrument int64
	Outcome           float64
	Date              string
	Deleted           bool
	IncomeAccount     string
	OutcomeAccount    string
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
