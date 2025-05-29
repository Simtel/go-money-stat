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

	mock.ExpectQuery("SELECT (.+) FROM `transactions`").
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
