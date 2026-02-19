package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/filipio/athletics-backend/pkg/httpio"
)

func OnlyCurrentUserMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		if strings.Contains(r.URL.Path, httpio.OnlyCurrentUserPath) {
			ctx = context.WithValue(ctx, httpio.OnlyCurrentUserContextKey, true)
		} else {
			ctx = context.WithValue(ctx, httpio.OnlyCurrentUserContextKey, false)
		}

		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
