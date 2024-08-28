package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/caarlos0/httperr/v2"
)

func main() {
	mux := http.NewServeMux()

	mux.Handle("/s", httperr.NewF(func(w http.ResponseWriter, r *http.Request) error {
		fmt.Fprintln(w, "this is OK")
		return nil
	}))

	mux.Handle("/f", httperr.NewF(func(w http.ResponseWriter, r *http.Request) error {
		return fmt.Errorf("this will be a 500")
	}))

	mux.Handle("/e", httperr.NewF(func(w http.ResponseWriter, r *http.Request) error {
		return httperr.Wrap(fmt.Errorf("wrap another err into a bad request"), http.StatusBadRequest)
	}))

	mux.Handle("/ef", httperr.NewF(func(w http.ResponseWriter, r *http.Request) error {
		return httperr.Errorf(http.StatusConflict, "create new conflict error: %v", 123)
	}))

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	log.Fatal(server.ListenAndServe())
}
