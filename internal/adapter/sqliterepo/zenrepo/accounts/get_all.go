package accounts

import (
	"log"
	"money-stat/internal/model"
)

func (r *Repository) GetAll() []model.Account {
	var accounts []model.Account

	err := r.db.Joins("Currency").Find(&accounts).Error

	if err != nil {
		log.Println(err)
	}

	return accounts
}
