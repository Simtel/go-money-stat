package app

type App struct {
	c *Container
}

func NewApp(c *Container) *App {
	return &App{c: c}
}

func (a *App) GetContainer() *Container {
	return a.c
}
