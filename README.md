# httperr

I've been doing this in several different projects, I finally decided to
convert it to a proper lib.

The idea is to add an `error` return to HTTP handler functions, so you can
avoid writing `if err != nil { http.Error(w, err); return }` everywhere.

The basic usage looks like:

```go
mux.Handle("/", httperr.NewF(func(w http.ResponseWriter, r *http.Request) error {
  err := doSomething()
  return err
}))
```

This will yield a `500` and return a JSON like `{"error":"doSomething error"}`.

The lib also provide a `Wrap` function, so you can decide which status code
you want:

```go
mux.Handle("/e", httperr.NewF(func(w http.ResponseWriter, r *http.Request) error {
  err := doSomething()
  return httperr.Wrap(err, http.StatusBadRequest)
}))
```

So, this is it! You can also check the `examples` folder for a "real" usage.
