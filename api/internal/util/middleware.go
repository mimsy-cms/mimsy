package util

import "net/http"

type Middleware func(http.Handler) http.Handler

// ApplyMiddlewares returns a middleware that will execute the provided
// middlewares in the order in which they were given as parameters.
func ApplyMiddlewares(middlewares ...Middleware) Middleware {
	return func(next http.Handler) http.Handler {
		for i := len(middlewares) - 1; i >= 0; i-- {
			middleware := middlewares[i]
			next = middleware(next)
		}

		return next
	}
}
