package middlewares

import (
	"context"
	"net/http"

	"github.com/evalabs-id/go-clean-architecture/internal/configs"
	"github.com/evalabs-id/go-clean-architecture/pkg/constant"
	"github.com/evalabs-id/go-clean-architecture/pkg/httphelper"
	"github.com/evalabs-id/go-clean-architecture/pkg/jwthelper"
)

type AuthMiddleware struct {
	Config *configs.Config
}

func ProvideAuthMiddleware(config *configs.Config) *AuthMiddleware {
	return &AuthMiddleware{
		Config: config,
	}
}

func (a *AuthMiddleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			httphelper.SimpleError(
				"Authorization header is missing. Please provide a valid token.",
				http.StatusUnauthorized,
				"MISSING_AUTHORIZATION",
			).ServeHTTP(w, r)
			return
		}

		if len(tokenString) < 7 || tokenString[:7] != "Bearer " {
			httphelper.SimpleError(
				"Invalid token format",
				http.StatusUnauthorized,
				"INVALID_TOKEN_FORMAT",
			).ServeHTTP(w, r)
			return
		}

		tokenString = tokenString[7:]

		claims, err := jwthelper.ValidateAccessToken(a.Config, tokenString)
		if err != nil {
			httphelper.SimpleError(
				"Invalid or expired token",
				http.StatusUnauthorized,
				"INVALID_OR_EXPIRED_TOKEN",
			).ServeHTTP(w, r)
			return
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, constant.UserIDContextKey, claims.UserID)
		ctx = context.WithValue(ctx, constant.UserEmailContextKey, claims.Email)
		ctx = context.WithValue(ctx, constant.UserClaimsContextKey, claims.TokenType)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
