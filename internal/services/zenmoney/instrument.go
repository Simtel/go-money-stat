package zenmoney

type Instrument struct {
	Id         int     `json:"id"`
	Title      string  `json:"title"`
	ShortTitle string  `json:"shortTitle"`
	Symbol     string  `json:"symbol"`
	Rate       float64 `json:"rate"`
}
