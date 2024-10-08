# httperr

[![Build Status](https://img.shields.io/github/actions/workflow/status/caarlos0/httperr/build.yml?branch=main&style=for-the-badge)](https://github.com/caarlos0/httperr/actions?workflow=build)
[![Coverage Status](https://img.shields.io/codecov/c/gh/caarlos0/httperr.svg?logo=codecov&style=for-the-badge)](https://codecov.io/gh/caarlos0/httperr)
[![](http://img.shields.io/badge/godoc-reference-5272B4.svg?style=for-the-badge)](http://godoc.org/github.com/caarlos0/httperr/v2)

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

Or, you can throw errors with a status, e.g.:

```go
mux.Handle("/e", httperr.NewF(func(w http.ResponseWriter, r *http.Request) error {
  if something {
  	return httperr.Errorf(http.StatusBadRequest, "something: %v", something)
  }
  return nil
}))
```

So, this is it! You can also check the `examples` folder for a "real" usage.
