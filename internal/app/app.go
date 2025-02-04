package app

import "errors"

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

func GetGlobalApp() (*App, error) {
	if a == nil {
		return nil, errors.New("global app is not initialized")
	}

	return a, nil
}

func (a *App) GetContainer() *Container {
	return a.c
}
