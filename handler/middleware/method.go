package middleware

import (
	"errors"
	"net/http"
	"server/util"
)

type ReqMethod string

const (
	GET    ReqMethod = http.MethodGet
	POST   ReqMethod = http.MethodPost
	DELETE ReqMethod = http.MethodDelete
)

func (r ReqMethod) String() string {
	return string(r)
}

func NewReqMethodMW(m ReqMethod) *reqMethodMW {
	return &reqMethodMW{method: m}
}

type reqMethodMW struct {
	method ReqMethod
}

func (m *reqMethodMW) Bind(next http.Handler) http.Handler {
	util.Log.FunctionInfo()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != m.method.String() {
			util.Log.Error(errors.New(http.StatusText(http.StatusMethodNotAllowed)))
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			return
		}

		next.ServeHTTP(w, r.WithContext(r.Context()))

	})
}
