package core

import "context"

type AppCtx = context.Context

type StartAndStoppable interface {
	Start(app *App) error
	Stop(app *App) error
	CanStart() bool
}
