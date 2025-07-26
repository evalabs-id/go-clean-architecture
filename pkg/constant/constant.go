package constant

import "time"

type ContextKey string

const (
	ApiV1 = "/api/v1"

	UserCacheKey      = "users:%s"
	UserCacheDuration = time.Hour * 2

	// Context keys for middleware
	UserIDContextKey     ContextKey = "user.id"
	UserEmailContextKey  ContextKey = "user.email"
	UserClaimsContextKey ContextKey = "user.claims"
)
