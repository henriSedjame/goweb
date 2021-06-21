package core

type ErrorHandledFunc func(fn func(app *App) error)

