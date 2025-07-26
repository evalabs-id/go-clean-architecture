package app

import (
	"github.com/evalabs-id/go-clean-architecture/internal/providers/logger"
	"github.com/rs/zerolog"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
)

type App struct {
	Modules      []fx.Option
	Servers      []any
	Dependencies []any
	Invokers     []any
}

func New() *App {
	return &App{
		Modules:      []fx.Option{},
		Servers:      []any{},
		Dependencies: []any{},
		Invokers:     []any{},
	}
}

func (app *App) AddModules(modules ...fx.Option) {
	app.Modules = append(app.Modules, modules...)
}

func (app *App) AddServers(servers ...any) {
	app.Servers = append(app.Servers, servers...)
}

func (app *App) AddDependencies(dependencies ...any) {
	app.Dependencies = append(app.Dependencies, dependencies...)
}

func (app *App) AddInvokers(invokers ...any) {
	app.Invokers = append(app.Invokers, invokers...)
}

func (app *App) Run() {
	var opts = []fx.Option{}

	opts = append(opts, app.Modules...)
	opts = append(opts, fx.Provide(app.Dependencies...))
	opts = append(opts, fx.Provide(app.Servers...))
	opts = append(opts, fx.Invoke(app.Invokers...))
	opts = append(opts, fx.WithLogger(func(l zerolog.Logger) fxevent.Logger {
		return logger.Default(l)
	}))

	fx.New(opts...).Run()
}
