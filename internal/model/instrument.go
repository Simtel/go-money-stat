package model

type Instrument struct {
	Id         int `gorm:"primaryKey"`
	Title      string
	ShortTitle string
	Symbol     string
	Rate       float64
}
