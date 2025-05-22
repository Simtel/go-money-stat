package app

import (
	"money-stat/internal/adapter/sqliterepo/zenrepo/accounts"
	"money-stat/internal/adapter/sqliterepo/zenrepo/transactions"
)

type Container struct {
	db DbInterface
}

func NewContainer(db DbInterface) *Container {
	return &Container{db: db}
}

func (c *Container) GetDb() DbInterface {
	return c.db
}

func (c *Container) GetTransactionRepository() transactions.RepositoryInterface {
	return transactions.NewRepository(c.db.GetGorm())
}

func (c *Container) GetAccountRepository() accounts.RepositoryInterface {
	return accounts.NewRepository(c.db.GetGorm())
}
