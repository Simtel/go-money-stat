package app

type App struct {
	c *Container
}

func NewApp(c *Container) *App {
	return &App{c: c}
}

var a *App

func SetGlobalApp(app *App) {
	a = app
}

func GetGlobalApp() *App {
	if a == nil {
		panic("global app is not initialized")
	}

	return a
}

func (a *App) GetContainer() *Container {
	return a.c
}
