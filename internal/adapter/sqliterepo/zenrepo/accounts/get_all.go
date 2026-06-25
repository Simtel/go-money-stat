package accounts

import (
	"money-stat/internal/model"
)

func (r *Repository) GetAll() ([]model.Account, error) {
	var accounts []model.Account

	err := r.db.Joins("Currency").Find(&accounts).Error
	if err != nil {
		return nil, err
	}

	return accounts, nil
}
