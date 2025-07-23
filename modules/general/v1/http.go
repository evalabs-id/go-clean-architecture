package generalv1

import (
	"net/http"

	"github.com/evalabs-id/go-clean-architecture/pkg/helper/httphelper"
	"github.com/evalabs-id/go-clean-architecture/pkg/logger"
)

type GeneralHttp struct{}

func ProvideHttp() *GeneralHttp {
	return &GeneralHttp{}
}

func (g *GeneralHttp) EmailCheck(w http.ResponseWriter, r *http.Request) http.Handler {
	ctx := r.Context()
	logger := logger.FromContext(ctx)

	logger.Info().Msg("Email check endpoint hit")

	return httphelper.OK(map[string]string{
		"status":  "OK",
		"message": "Email check successful",
	})
}
