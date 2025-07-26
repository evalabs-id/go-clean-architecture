package middlewares

import (
	"net/http"
	"strings"

	"github.com/evalabs-id/go-clean-architecture/pkg/httphelper"
)

// ContentTypeValidator validates the Content-Type header for requests with body
func ContentTypeValidator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Only check content-type for requests that have a body
		if r.Method == "POST" || r.Method == "PUT" || r.Method == "PATCH" {
			contentType := r.Header.Get("Content-Type")
			if contentType == "" {
				httphelper.SimpleError(
					"Content-Type header is missing. Please set it to a supported value.",
					http.StatusUnsupportedMediaType,
					"MISSING_CONTENT_TYPE",
				).ServeHTTP(w, r)
				return
			}
			// Allow any application/json variant (with or without charset)
			if !strings.HasPrefix(contentType, "application/json") &&
				!strings.HasPrefix(contentType, "application/x-www-form-urlencoded") {
				httphelper.SimpleError(
					"Please ensure your request uses a supported content type.",
					http.StatusUnsupportedMediaType,
					"UNSUPPORTED_CONTENT_TYPE",
				).ServeHTTP(w, r)
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}
