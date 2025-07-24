package authv1

import (
	"github.com/evalabs-id/go-clean-architecture/internal/providers"
	"github.com/evalabs-id/go-clean-architecture/pkg/constant"
	"github.com/evalabs-id/go-clean-architecture/pkg/httphelper"
	"github.com/go-chi/chi/v5"
	"go.uber.org/fx"
)

func ProvideRoutes(authHttp *Http) providers.Router {
	return providers.Router{
		Pattern: constant.ApiV1 + "/auth",
		SubRouter: func(r chi.Router) {
			r.Handle("POST /sign-up", httphelper.HttpHandlerFunc(authHttp.SignUp))
			r.Handle("POST /sign-in", httphelper.HttpHandlerFunc(authHttp.SignIn))
		},
		Public: true, // These are public endpoints
	}
}

var Module = fx.Module(
	"auth",
	fx.Provide(ProvideRepository),
	fx.Provide(ProvideService),
	fx.Provide(ProvideHttp),
	fx.Provide(
		fx.Annotate(ProvideRoutes,
			fx.ResultTags(`group:"routers"`),
		),
	),
)
