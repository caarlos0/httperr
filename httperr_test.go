package httperr

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestServeNoError(t *testing.T) {
	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	NewF(func(w http.ResponseWriter, r *http.Request) error {
		return nil
	}).ServeHTTP(w, req)
	resp := w.Result()
	require.Equal(t, http.StatusOK, resp.StatusCode)
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
	require.Equal(t, status, resp.StatusCode)
	bts, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)
	require.Contains(t, string(bts), msg)
}

func TestServeError(t *testing.T) {
	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	msg := "server is doing funky stuff"
	NewF(func(w http.ResponseWriter, r *http.Request) error {
		return fmt.Errorf(msg)
	}).ServeHTTP(w, req)
	resp := w.Result()
	require.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	bts, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)
	require.Contains(t, string(bts), msg)
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
			var err = Wrap(tt.args.err, tt.args.status)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
