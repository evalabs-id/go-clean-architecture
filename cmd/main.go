package main

import (
	"net/http"

	"github.com/evalabs-id/go-clean-architecture/cmd/app"
	"github.com/evalabs-id/go-clean-architecture/internal/configs"
	"github.com/evalabs-id/go-clean-architecture/internal/middlewares"
	"github.com/evalabs-id/go-clean-architecture/internal/providers/database"
	httpserver "github.com/evalabs-id/go-clean-architecture/internal/providers/http_server"
	"github.com/evalabs-id/go-clean-architecture/internal/providers/logger"
	authv1 "github.com/evalabs-id/go-clean-architecture/modules/auth/v1"
	generalv1 "github.com/evalabs-id/go-clean-architecture/modules/general/v1"

	_ "github.com/lib/pq"
)

func main() {
	apps := app.New()

	apps.AddDependencies(
		configs.ProvideConfig,
		logger.ProvideZerolog,
		database.ProvideDB,
		middlewares.ProvideHttpMiddleware,
		middlewares.ProvideAuthMiddleware,
		httpserver.ProvideRoutes,
	)

	apps.AddModules(
		authv1.Module,
		generalv1.Module,
	)

	apps.AddServers(httpserver.HttpServer)

	apps.AddInvokers(func(*http.Server) {})

	apps.Run()
}
