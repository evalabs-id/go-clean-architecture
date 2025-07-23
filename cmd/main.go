package main

import (
	"net/http"

	"github.com/evalabs-id/go-clean-architecture/cmd/app"
	"github.com/evalabs-id/go-clean-architecture/internal/configs"
	"github.com/evalabs-id/go-clean-architecture/internal/middlewares"
	"github.com/evalabs-id/go-clean-architecture/internal/providers"
	generalv1 "github.com/evalabs-id/go-clean-architecture/modules/general/v1"
	_ "github.com/lib/pq"
)

func main() {
	apps := app.New()

	apps.AddDependencies(
		configs.ProvideConfig,
		providers.ProvideZerolog,
		providers.SqlxDB,
		middlewares.ProvideHttpMiddleware,
		middlewares.ProvideAuthenticationMiddleware,
		providers.Routes,
	)

	apps.AddModules(
		generalv1.Module,
	)

	apps.AddServers(providers.HttpServer)

	apps.AddInvokers(func(*http.Server) {})

	apps.Run()
}
