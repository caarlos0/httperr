// Package httperr provides a middleware to make it easier to handle http
// errors in a common way.
package httperr

import (
	"encoding/json"
	"errors"
	"fmt"
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
	return Error{ // nolint: wsl
		Err:    err,
		Status: status,
	}
}

// Errorf creates a new error and wraps it with the given status
func Errorf(status int, format string, args ...interface{}) error {
	return Wrap(fmt.Errorf(format, args...), status)
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

// ErrorHandler handles an error.
type ErrorHandler func(w http.ResponseWriter, err error, status int)

// DefaultErrorHandler is the default error handler.
// It converts the error to JSON and prints writes it to the response.
func DefaultErrorHandler(w http.ResponseWriter, err error, status int) {
	msg := err.Error()
	bts, _ := json.Marshal(errorResponse{
		Error: msg,
	})
	http.Error(w, string(bts), status)
}

// NewWithHandler() wraps a given http.Handler and returns a http.Handler.
// You can also customize how the error is handled.
func NewWithHandler(next Handler, eh ErrorHandler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := next.ServeHTTP(w, r)
		if err == nil {
			return
		}

		herr := Error{}
		if errors.As(err, &herr) {
			eh(w, herr, herr.Status)
		} else {
			eh(w, err, http.StatusInternalServerError)
		}
	})
}

// New wraps a given http.Handler and returns a http.Handler.
func New(next Handler) http.Handler {
	return NewWithHandler(next, DefaultErrorHandler)
}

// NewFWithHandler wraps a given http.HandlerFunc and return a http.Handler.
// You can also customize how the error is handled.
func NewFWithHandler(next HandlerFunc, eh ErrorHandler) http.Handler { // nolint: interfacer
	return NewWithHandler(next, eh)
}

// NewF wraps a given http.HandlerFunc and return a http.Handler.
func NewF(next HandlerFunc) http.Handler { // nolint: interfacer
	return New(next)
}

type errorResponse struct {
	Error string `json:"error"`
}
