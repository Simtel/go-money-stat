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
	IncomeAccount     string `gorm:"column:income_account"`
	OutcomeAccount    string `gorm:"column:outcome_account"`
	TagIds            string `gorm:"column:tag_ids"`
	Comment           string
	Tag               []Tag   `gorm:"-"`
	InAccount  Account `gorm:"foreignKey:IncomeAccount;references:Id"`
	OutAccount Account `gorm:"foreignKey:OutcomeAccount;references:Id"`
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
	if len(t.Tag) == 0 {
		return "Перевод"
	}
	result := ""
	for _, tag := range t.Tag {
		result += tag.Title + " "
	}
	return result
}
