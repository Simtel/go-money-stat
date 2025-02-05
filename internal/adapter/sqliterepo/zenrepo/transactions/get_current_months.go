package transactions

import (
	"fmt"
	"log"
	"money-stat/internal/model"
	"strconv"
)

func (r *Repository) GetCurrentMonths() {
	db := r.db
	var transactions []model.Transaction
	err := db.Model(&model.Transaction{}).Where("date BETWEEN ? and ? ", "2024-11-11", "2024-11-12").Preload("Tag").Find(&transactions).Error
	if err != nil {
		log.Println(err)
	}
	fmt.Println(transactions)
	fmt.Println("В базе транзакций: " + strconv.Itoa(len(transactions)))
}
