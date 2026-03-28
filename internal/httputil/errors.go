package httputil

import (
	"fmt"
	"net/http"
)

// CustomError represents a structured application error with HTTP context.
type CustomError struct {
	BaseErr     error                  `json:"-"`
	StatusCode  int                    `json:"status_code"`
	Message     string                 `json:"message"`
	UserMessage string                 `json:"user_message"`
	ErrType     string                 `json:"error_type"`
	ErrCode     string                 `json:"error_code"`
	Retryable   bool                   `json:"retryable"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// Error implements the error interface.
func (e *CustomError) Error() string {
	return fmt.Sprintf("[%s] %s: %s", e.ErrCode, e.ErrType, e.Message)
}

// Unwrap returns the underlying error for errors.Is/As support.
func (e *CustomError) Unwrap() error {
	return e.BaseErr
}

// WithMetadata adds a key-value pair to the error metadata and returns the error for chaining.
func (e *CustomError) WithMetadata(key string, value interface{}) *CustomError {
	e.Metadata[key] = value
	return e
}

// NewFromError creates a CustomError from an existing error.
func NewFromError(err error, statusCode int, userMessage, errType, errCode string, retryable bool) *CustomError {
	return &CustomError{
		BaseErr:     err,
		StatusCode:  statusCode,
		Message:     err.Error(),
		UserMessage: userMessage,
		ErrType:     errType,
		ErrCode:     errCode,
		Retryable:   retryable,
		Metadata:    make(map[string]interface{}),
	}
}

// --- Common error constructors ---

// NewBadRequest creates a 400 Bad Request error.
func NewBadRequest(err error, userMessage string) *CustomError {
	return NewFromError(err, http.StatusBadRequest, userMessage, "VALIDATION_ERROR", "BAD_REQUEST", false)
}

// NewNotFound creates a 404 Not Found error.
func NewNotFound(err error, userMessage string) *CustomError {
	return NewFromError(err, http.StatusNotFound, userMessage, "NOT_FOUND", "RESOURCE_NOT_FOUND", false)
}

// NewInternalError creates a 500 Internal Server Error.
func NewInternalError(err error) *CustomError {
	return NewFromError(err, http.StatusInternalServerError, "Something went wrong", "INTERNAL_ERROR", "INTERNAL_SERVER_ERROR", true)
}

// NewUnauthorized creates a 401 Unauthorized error.
func NewUnauthorized(err error, userMessage string) *CustomError {
	return NewFromError(err, http.StatusUnauthorized, userMessage, "AUTH_ERROR", "UNAUTHORIZED", false)
}
