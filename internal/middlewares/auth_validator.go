package middlewares

import (
	"context"
	"net/http"

	"github.com/evalabs-id/go-clean-architecture/pkg/constant"
	"github.com/evalabs-id/go-clean-architecture/pkg/helper/httphelper"
	"github.com/evalabs-id/go-clean-architecture/pkg/helper/jwthelper"
)

type AuthenticationMiddleware struct {
	JWTHelper *jwthelper.JWTHelper
}

func ProvideAuthenticationMiddleware(j *jwthelper.JWTHelper) *AuthenticationMiddleware {
	return &AuthenticationMiddleware{
		JWTHelper: j,
	}
}

func (a *AuthenticationMiddleware) Authenticate(next http.Handler) http.Handler {
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

		claims, err := a.JWTHelper.ValidateAccessToken(tokenString)
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
