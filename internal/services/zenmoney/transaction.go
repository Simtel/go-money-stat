package zenmoney

type Transaction struct {
	Id                string  `json:"id"`
	Changed           int64   `json:"changed"`
	Created           int64   `json:"created"`
	IncomeInstrument  int64   `json:"incomeInstrument"`
	Income            float64 `json:"income"`
	OutcomeInstrument int64   `json:"outcomeInstrument"`
	Outcome           float64 `json:"outcome"`
	Date              string  `json:"date"`
}

func (t Transaction) FormatAmount() float64 {
	if t.Income == 0 && t.Outcome > 0 {
		return -t.Outcome
	}

	if t.Income > 0 && t.Outcome == 0 {
		return t.Income
	}

	return 0
}
