package model

type Account struct {
	Id         string
	Title      string
	Balance    float64
	Instrument int
	Currency   Instrument `gorm:"foreignKey:Instrument"`
}

func (a *Account) IsRuble() bool {
	return a.Instrument == 2
}

func (a *Account) IsDollar() bool {
	return a.Instrument == 1
}
