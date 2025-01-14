package zenmoney

type Response struct {
	Account     []Account     `json:"account"`
	Instrument  []Instrument  `json:"instrument"`
	Transaction []Transaction `json:"transaction"`
}
