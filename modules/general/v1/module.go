package generalv1

import (
	httpserver "github.com/evalabs-id/go-clean-architecture/internal/providers/http_server"
	"github.com/evalabs-id/go-clean-architecture/pkg/constant"
	"github.com/evalabs-id/go-clean-architecture/pkg/httphelper"
	"github.com/go-chi/chi/v5"
	"go.uber.org/fx"
)

func ProvideRoutes(generalHttp *GeneralHttp) httpserver.Router {
	return httpserver.Router{
		Pattern: constant.ApiV1 + "/general",
		SubRouter: func(r chi.Router) {
			r.Handle("POST /email-check", httphelper.HttpHandlerFunc(generalHttp.EmailCheck))
		},
		Public: true,
	}
}

var Module = fx.Module(
	"general",
	fx.Provide(ProvideHttp),
	fx.Provide(
		fx.Annotate(ProvideRoutes,
			fx.ResultTags(`group:"routers"`),
		),
	),
)
