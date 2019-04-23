package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/caarlos0/httperr"
)

func main() {
	var mux = http.NewServeMux()

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

	var server = &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	log.Fatal(server.ListenAndServe())
}
