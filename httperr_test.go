package httperr

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

func TestServeNoError(t *testing.T) {
	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	NewF(func(w http.ResponseWriter, r *http.Request) error {
		return nil
	}).ServeHTTP(w, req)
	resp := w.Result()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected http status ok, got %d", resp.StatusCode)
	}
}

func TestServeHTTPError(t *testing.T) {
	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	status := http.StatusBadRequest
	msg := "data was wrong"
	NewF(func(w http.ResponseWriter, r *http.Request) error {
		return Error{
			Status: status,
			Err:    fmt.Errorf(msg),
		}
	}).ServeHTTP(w, req)
	resp := w.Result()
	if resp.StatusCode != status {
		t.Fatalf("expected http status %d, got %d", status, resp.StatusCode)
	}

	bts, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("got error: %v", err)
	}
	if !strings.Contains(string(bts), msg) {
		t.Fatalf("%q does not contain %q", string(bts), msg)
	}
}

func TestServeError(t *testing.T) {
	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	msg := "server is doing funky stuff"
	NewF(func(w http.ResponseWriter, r *http.Request) error {
		return fmt.Errorf(msg)
	}).ServeHTTP(w, req)
	resp := w.Result()

	status := http.StatusInternalServerError
	if resp.StatusCode != status {
		t.Fatalf("expected http status %d, got %d", status, resp.StatusCode)
	}

	bts, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("got error: %v", err)
	}
	if !strings.Contains(string(bts), msg) {
		t.Fatalf("%q does not contain %q", string(bts), msg)
	}
}

func TestWrap(t *testing.T) {
	type args struct {
		err    error
		status int
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "nil err",
			args: args{
				err:    nil,
				status: 400,
			},
			wantErr: false,
		},
		{
			name: "nonnil err",
			args: args{
				err:    fmt.Errorf("errr"),
				status: 404,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Wrap(tt.args.err, tt.args.status)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error, got none")
				}
			} else {
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}
			}
		})
	}
}

func TestErrorIs(t *testing.T) {
	t.Run("outer", func(t *testing.T) {
		if !errors.Is(Wrap(io.EOF, http.StatusGone), Error{}) {
			t.Fatalf("underlying error should be Error")
		}
	})

	t.Run("is", func(t *testing.T) {
		if !errors.Is(Wrap(io.EOF, http.StatusGone), io.EOF) {
			t.Fatalf("underlying error should be io.EOF")
		}
	})

	t.Run("is not", func(t *testing.T) {
		if errors.Is(Wrap(io.ErrUnexpectedEOF, http.StatusGone), io.EOF) {
			t.Fatalf("underlying error should not be io.EOF")
		}
	})
}

func TestErrorf(t *testing.T) {
	err := Errorf(http.StatusConflict, "foo bar %d", 10)
	expectedErr := Error{
		Err:    fmt.Errorf("foo bar 10"),
		Status: http.StatusConflict,
	}

	if !reflect.DeepEqual(err, expectedErr) {
		t.Fatalf("errors does not match: %v", err)
	}
}
