package generalv1

import (
	"github.com/evalabs-id/go-clean-architecture/internal/providers"
	"github.com/evalabs-id/go-clean-architecture/pkg/constant"
	"github.com/evalabs-id/go-clean-architecture/pkg/helper/httphelper"
	"github.com/go-chi/chi/v5"
	"go.uber.org/fx"
)

func ProvideRoutes(generalHttp *GeneralHttp) providers.Router {
	return providers.Router{
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
