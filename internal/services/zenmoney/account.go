package zenmoney

type Account struct {
	Id         string  `json:"id"`
	Title      string  `json:"title"`
	Balance    float64 `json:"balance"`
	Instrument int     `json:"instrument"`
}

func (a *Account) IsRuble() bool {
	return a.Instrument == 2
}

func (a *Account) IsDollar() bool {
	return a.Instrument == 1
}
