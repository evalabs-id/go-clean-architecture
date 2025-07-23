package providers

import (
	"net/http"

	"github.com/evalabs-id/go-clean-architecture/internal/middlewares"
	"github.com/evalabs-id/go-clean-architecture/pkg/helper/httphelper"
	"github.com/evalabs-id/go-clean-architecture/pkg/helper/jwthelper"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/rs/zerolog"
	"go.uber.org/fx"
)

type Router struct {
	Public    bool
	Pattern   string
	SubRouter chi.Router
}

var corsOptions = cors.Options{
	AllowedOrigins:   []string{"*"},
	AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
	AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "x-api-key"},
	ExposedHeaders:   []string{"Link"},
	AllowCredentials: true,
}

type routesParams struct {
	fx.In
	ModuleRouters []Router `group:"routers"`
	JWTHelper     *jwthelper.JWTHelper
	Logger        zerolog.Logger
}

func Routes(params routesParams) *chi.Mux {
	r := chi.NewRouter()
	httpMiddleware := middlewares.ProvideHttpMiddleware(params.Logger)
	authMiddleware := middlewares.ProvideAuthenticationMiddleware(params.JWTHelper)

	// default middleware mounted
	r.Use(cors.Handler(corsOptions))
	r.Use(middleware.Compress(5))
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(httpMiddleware.LoggerContext)
	r.Use(httpMiddleware.Handler)
	r.Use(httpMiddleware.Recoverer)

	r.Handle("GET /health", httphelper.OK(map[string]string{
		"status":  "OK",
		"message": "I'm alive!",
	}))

	// Mount routes based on protection level
	r.Group(func(r chi.Router) {
		// Content-Type validation middleware
		r.Use(middlewares.ContentTypeValidator)

		r.Group(func(r chi.Router) {
			r.Use(authMiddleware.Authenticate)
			// Protected routes
			for _, router := range params.ModuleRouters {
				if !router.Public {
					r.Mount(router.Pattern, router.SubRouter)
				}
			}
		})

		// Public routes
		for _, router := range params.ModuleRouters {
			if router.Public {
				r.Mount(router.Pattern, router.SubRouter)
			}
		}
	})

	// TODO: add fallbacks
	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		httphelper.SimpleError(
			"Route Not Found",
			http.StatusNotFound,
			"ROUTE_NOT_FOUND",
		).ServeHTTP(w, r)
	})
	r.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		httphelper.SimpleError(
			"Method Not Allowed",
			http.StatusMethodNotAllowed,
			"METHOD_NOT_ALLOWED",
		).ServeHTTP(w, r)
	})

	return r
}
