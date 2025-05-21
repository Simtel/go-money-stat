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
	Tag               []Tag   `gorm:"many2many:transaction_tags;association_autocreate:false"`
	InAccount         Account `gorm:"foreignKey:IncomeAccount"`
	OutAccount        Account `gorm:"foreignKey:OutcomeAccount"`
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

func (t Transaction) IsDeleted() bool {
	return t.Deleted
}

func (t Transaction) IsIncome() bool {
	return t.Income > 0 && t.Outcome == 0
}

func (t Transaction) IsOutcome() bool {
	return t.Outcome > 0 && t.Income == 0
}

func (t Transaction) IsTransfer() bool {
	return t.Income > 0 && t.Outcome > 0
}

func (t Transaction) GetTagsTitle() string {
	var transactionTags string
	for _, tag := range t.Tag {
		transactionTags += tag.Title + " "
	}

	if transactionTags == "" {
		transactionTags = "Перевод"
	}
	return transactionTags
}
