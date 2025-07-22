package constant

import "time"

type ContextKey string

const (
	ApiV1 = "/api/v1"

	UserCacheKey      = "users:%s"
	UserCacheDuration = time.Hour * 2

	// Context keys for middleware
	UserEmailContextKey    ContextKey = "user.email"
	UserFullnameContextKey ContextKey = "user.fullname"
	UserClaimsContextKey   ContextKey = "user.claims"
)
