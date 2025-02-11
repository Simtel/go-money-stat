package app

import (
	"money-stat/internal/adapter/sqliterepo/zenrepo/accounts"
	"money-stat/internal/adapter/sqliterepo/zenrepo/transactions"
)

type Container struct {
	db *DB
}

func NewContainer(db *DB) *Container {
	return &Container{db: db}
}

func (c *Container) GetDb() *DB {
	return c.db
}

func (c *Container) GetTransactionRepository() *transactions.Repository {
	return transactions.NewRepository(c.db.GetGorm())
}

func (c *Container) GetAccountRepository() *accounts.Repository {
	return accounts.NewRepository(c.db.GetGorm())
}
