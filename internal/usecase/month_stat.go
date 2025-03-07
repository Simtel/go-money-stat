package usecase

type MonthStatDto struct {
	Transactions []MonthStatTransactionDto
	OutComeSumm  float64
	InComeSumm   float64
	Count        int
}

type MonthStatTransactionDto struct {
	Date         string
	Tags         string
	FormatAmount string
	Account      string
	CreatedAt    string
}
