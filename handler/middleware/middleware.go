package middleware

import (
	"net/http"
	i "server/interfaces"
	"server/lib"
)

type Middleware interface {
	Bind(next http.Handler) http.Handler
}

// Middlewares wrapped right to left and are executed in left to right
// (slice reversed within to achieve this)
func AddMiddleware(h http.Handler, mws ...Middleware) http.Handler {
	var handler http.Handler = h
	lib.ReverseInPlace(mws)
	for _, mw := range mws {
		handler = mw.Bind(handler)
	}
	return handler
}

func WithAuth(a i.AuthService) *authMW {
	return NewAuthMW(a)
}

func WithMethod(method ReqMethod) *reqMethodMW {
	return NewReqMethodMW(method)
}

func WithNoAuth() *noAuthMW {
	return NewNoAuthMW()
}
