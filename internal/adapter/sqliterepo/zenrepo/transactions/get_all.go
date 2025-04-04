package transactions

import "money-stat/internal/model"

func (r *Repository) GetAll(includeDeleted bool) []model.Transaction {
	db := r.db
	var transactions []model.Transaction
	err := db.Model(&model.Transaction{}).
		Where(" deleted = ?", includeDeleted).
		Order("date ASC").
		Find(&transactions).
		Error
	if err != nil {
		panic(err)
	}

	return transactions
}
