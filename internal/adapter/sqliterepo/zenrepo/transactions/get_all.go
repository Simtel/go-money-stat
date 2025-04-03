package transactions

import "money-stat/internal/model"

func (r *Repository) GetAll() []model.Transaction {
	db := r.db
	var transactions []model.Transaction
	err := db.Model(&model.Transaction{}).Order("date ASC").Find(&transactions).Error
	if err != nil {
		panic(err)
	}

	return transactions
}
