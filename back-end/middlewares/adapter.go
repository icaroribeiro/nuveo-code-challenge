package middlewares

import (
    "net/http"
)

type AdaptedHandlerFunc struct {
    handlerFunc http.HandlerFunc
}

type Adapter func(http.HandlerFunc) http.HandlerFunc

func AdaptFunc(f func(w http.ResponseWriter, r *http.Request)) *AdaptedHandlerFunc {
    return &AdaptedHandlerFunc{handlerFunc: f}
}

// It applies a chain of middlewares to the http.HandlerFunc.
func (a *AdaptedHandlerFunc) With(adapters ...Adapter) http.HandlerFunc {
    var h http.HandlerFunc
    var last int
    var i int

    h = a.handlerFunc

    last = len(adapters) - 1

    for i = range adapters {
        h = adapters[last-i](h)
    }

    return h
}
