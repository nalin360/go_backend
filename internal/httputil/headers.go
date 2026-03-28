package httputil

import "net/http"

// GetHeaders extracts the X-Request-ID and Authorization headers from the request.
func GetHeaders(r *http.Request) (requestID string, auth string) {
	requestID = r.Header.Get("X-Request-ID")
	auth = r.Header.Get("Authorization")
	return requestID, auth
}
