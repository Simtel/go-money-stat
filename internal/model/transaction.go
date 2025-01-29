package model

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
