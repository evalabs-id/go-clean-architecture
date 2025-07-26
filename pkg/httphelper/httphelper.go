package httphelper

import (
	"net/http"

	"github.com/bytedance/sonic"
)

// Response represents the standardized API response structure
type Response struct {
	Status int    `json:"status"`
	Data   any    `json:"data,omitempty"`
	Error  *Error `json:"error,omitempty"`
}

// Error represents detailed error information
type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details any    `json:"details,omitempty"`
}

func (e *Error) Error() string {
	return e.Message
}

func (e *Error) WithDetails(details any) *Error {
	return &Error{
		Code:    e.Code,
		Message: e.Message,
		Details: details,
	}
}

type HttpHandlerFunc func(w http.ResponseWriter, r *http.Request) http.Handler

func (fn HttpHandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if handler := fn(w, r); handler != nil {
		handler.ServeHTTP(w, r)
	}
}

func OK(data any) http.Handler {
	return jsonResponse(data, http.StatusOK)
}

func Created(data any) http.Handler {
	return jsonResponse(data, http.StatusCreated)
}

func NoContent() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})
}

func SimpleError(message string, statusCode int, code string, details ...any) http.Handler {
	var finalDetails any
	if len(details) == 1 {
		finalDetails = details[0]
	} else if len(details) > 1 {
		finalDetails = details
	}

	return jsonErrorResponse(&Error{
		Code:    code,
		Message: message,
		Details: finalDetails,
	}, statusCode)
}

func jsonResponse(data any, statusCode int) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.WriteHeader(statusCode)

		response := Response{
			Status: statusCode,
			Data:   data,
		}

		jsonBytes, _ := sonic.Marshal(response)
		w.Write(jsonBytes)
	})
}

func toCustomError(err error) *Error {
	if customErr, ok := err.(*Error); ok {
		return customErr
	}

	return &Error{
		Code:    "INTERNAL_SERVER_ERROR",
		Message: "An unexpected error occurred. Please try again later.",
		Details: err.Error(),
	}
}

func jsonErrorResponse(err error, statusCode int) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.WriteHeader(statusCode)

		customErr := toCustomError(err)

		response := Response{
			Status: statusCode,
			Error:  customErr,
		}

		jsonBytes, _ := sonic.Marshal(response)
		w.Write(jsonBytes)
	})
}
