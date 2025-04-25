package transactions

import "money-stat/internal/model"

func (r *Repository) GetAll() ([]model.Transaction, error) {
	db := r.db
	var transactions []model.Transaction
	err := db.Model(&model.Transaction{}).
		Where(" deleted = ?", 0).
		Joins("InAccount").
		Joins("OutAccount").
		Joins("InAccount.Currency").
		Joins("OutAccount.Currency").
		Order("date ASC").
		Find(&transactions).
		Error
	if err != nil {
		return nil, err
	}

	return transactions, nil
}
