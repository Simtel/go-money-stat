package app

type Container struct {
	db *DB
}

func NewContainer(db *DB) *Container {
	return &Container{db: db}
}

func (c *Container) GetDb() *DB {
	return c.db
}
