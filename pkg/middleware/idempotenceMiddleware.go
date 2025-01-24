package middleware

import (
	"im/pkg/interceptor"
	"net/http"
)

type IdempotenceMiddleware struct {
}

func (m *IdempotenceMiddleware) Handler(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r = r.WithContext(interceptor.ContextWithVal(r.Context()))
		next(w, r)
	}
}
