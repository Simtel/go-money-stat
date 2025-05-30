package transactions

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"regexp"
	"testing"
)

func TestRepository_GetByYear(t *testing.T) {
	repository, mock := getRepository()

	rows := sqlmock.NewRows([]string{"date", "income", "outcome", "amount", "comment", "tags", "created_at"}).
		AddRow("2021-09-01", 100, 0, 100, "", "", 0).
		AddRow("2021-09-02", 200, 0, 200, "", "", 0)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `transactions` WHERE date BETWEEN ? and ? AND deleted = ? ORDER BY date ASC")).
		WithArgs("2021-01-01", "2021-12-31", 0).
		WillReturnRows(rows)

	transactions := repository.GetByYear(2021)
	assert.Equal(t, 2, len(transactions))
}

func TestRepository_GetAll(t *testing.T) {
	repository, mock := getRepository()

	rows := sqlmock.NewRows([]string{"date", "income", "outcome", "amount", "comment", "tags", "created_at"}).
		AddRow("2021-09-01", 100, 0, 100, "", "", 0).
		AddRow("2021-09-02", 200, 0, 200, "", "", 0)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT `transactions`.`id`,`transactions`.`changed`,`transactions`.`created`,`transactions`.`income_instrument`,`transactions`.`income`,`transactions`.`outcome_instrument`,`transactions`.`outcome`,`transactions`.`date`,`transactions`.`deleted`,`transactions`.`income_account`,`transactions`.`outcome_account`,`InAccount`.`id` AS `InAccount__id`,`InAccount`.`title` AS `InAccount__title`,`InAccount`.`balance` AS `InAccount__balance`,`InAccount`.`start_balance` AS `InAccount__start_balance`,`InAccount`.`instrument` AS `InAccount__instrument`,`OutAccount`.`id` AS `OutAccount__id`,`OutAccount`.`title` AS `OutAccount__title`,`OutAccount`.`balance` AS `OutAccount__balance`,`OutAccount`.`start_balance` AS `OutAccount__start_balance`,`OutAccount`.`instrument` AS `OutAccount__instrument`,`InAccount__Currency`.`id` AS `InAccount__Currency__id`,`InAccount__Currency`.`title` AS `InAccount__Currency__title`,`InAccount__Currency`.`short_title` AS `InAccount__Currency__short_title`,`InAccount__Currency`.`symbol` AS `InAccount__Currency__symbol`,`InAccount__Currency`.`rate` AS `InAccount__Currency__rate`,`OutAccount__Currency`.`id` AS `OutAccount__Currency__id`,`OutAccount__Currency`.`title` AS `OutAccount__Currency__title`,`OutAccount__Currency`.`short_title` AS `OutAccount__Currency__short_title`,`OutAccount__Currency`.`symbol` AS `OutAccount__Currency__symbol`,`OutAccount__Currency`.`rate` AS `OutAccount__Currency__rate` FROM `transactions` LEFT JOIN `accounts` `InAccount` ON `transactions`.`income_account` = `InAccount`.`id` LEFT JOIN `accounts` `OutAccount` ON `transactions`.`outcome_account` = `OutAccount`.`id` LEFT JOIN `instruments` `InAccount__Currency` ON `InAccount`.`instrument` = `InAccount__Currency`.`id` LEFT JOIN `instruments` `OutAccount__Currency` ON `OutAccount`.`instrument` = `OutAccount__Currency`.`id` WHERE deleted = ? ORDER BY date ASC")).
		WithArgs(0).
		WillReturnRows(rows)

	transactions, _ := repository.GetAll()
	assert.Equal(t, 2, len(transactions))
}

func TestRepository_GetCurrentMonth(t *testing.T) {
	repository, mock := getRepository()

	rows := sqlmock.NewRows([]string{"date", "income", "outcome", "amount", "comment", "tags", "created_at"}).
		AddRow("2021-09-01", 100, 0, 100, "", "", 0).
		AddRow("2021-09-02", 200, 0, 200, "", "", 0)

	mock.ExpectQuery("SELECT (.+) FROM `transactions`").
		WillReturnRows(rows)

	transactions := repository.GetCurrentMonth()
	assert.Equal(t, 2, len(transactions))
}

func getRepository() (RepositoryInterface, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()

	if err != nil {
		log.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
	}

	gormDB, _ := gorm.Open(mysql.New(mysql.Config{
		Conn:                      db,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})
	return &Repository{db: gormDB}, mock
}
