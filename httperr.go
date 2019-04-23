// Package httperr provides a middleware to make it easier to handle http
// errors in a common way.
package httperr

import (
	"encoding/json"
	"net/http"
)

// Error is a HTTP error with an underlying error and a status code.
type Error struct {
	Err    error
	Status int
}

func (e Error) Error() string {
	return e.Err.Error()
}

// Wrap a given error with the given status.
// Returns nil if the given error is nil.
func Wrap(err error, status int) error {
	if err == nil {
		return nil
	}
	return Error{
		Err:    err,
		Status: status,
	}
}

// Handler is like http.Handler, but the ServeHTTP method can also return
// an error.
type Handler interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request) error
}

// HandlerFunc is just like http.HandlerFunc.
type HandlerFunc func(http.ResponseWriter, *http.Request) error

// ServeHTTP calls f(w, r).
func (f HandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) error {
	return f(w, r)
}

// New wraps a given http.Handler and returns a http.Handler.
func New(next Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := next.ServeHTTP(w, r)
		if err == nil {
			return
		}

		var msg = err.Error()
		var status = http.StatusInternalServerError
		herr, ok := err.(Error)
		if ok {
			status = herr.Status
		}
		bts, _ := json.Marshal(errorResponse{
			Error: msg,
		})
		http.Error(w, string(bts), status)
	})
}

// NewF wraps a given http.HandlerFunc and return a http.Handler
func NewF(next HandlerFunc) http.Handler { // nolint: interfacer
	return New(next)
}

type errorResponse struct {
	Error string `json:"error"`
}
