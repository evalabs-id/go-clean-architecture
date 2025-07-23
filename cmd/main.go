package main

import (
	"net/http"

	"github.com/evalabs-id/go-clean-architecture/cmd/app"
	"github.com/evalabs-id/go-clean-architecture/internal/configs"
	"github.com/evalabs-id/go-clean-architecture/internal/middlewares"
	"github.com/evalabs-id/go-clean-architecture/internal/providers"
	"github.com/evalabs-id/go-clean-architecture/pkg/helper/jwthelper"
	"github.com/evalabs-id/go-clean-architecture/pkg/logger"
	_ "github.com/lib/pq"
)

func main() {
	apps := app.New()

	apps.AddDependencies(
		configs.ProvideConfig,
		logger.ProvideZerolog,
		providers.SqlxDB,
		jwthelper.ProvideJWTHelper,
		middlewares.ProvideHttpMiddleware,
		middlewares.ProvideAuthenticationMiddleware,
		providers.Routes,
	)

	apps.AddModules()

	apps.AddServers(providers.HttpServer)

	apps.AddInvokers(func(*http.Server) {})

	apps.Run()
}
