package zenmoney

type Transaction struct {
	Id                string  `json:"id"`
	Changed           int64   `json:"changed"`
	Created           int64   `json:"created"`
	IncomeInstrument  int64   `json:"incomeInstrument"`
	Income            float64 `json:"income"`
	OutcomeInstrument int64   `json:"outcomeInstrument"`
	Outcome           float64 `json:"outcome"`
}
