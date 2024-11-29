package middlewares

import (
	"context"
	"net/http"
	"strings"

	"github.com/filipio/athletics-backend/utils"
)

func OnlyCurrentUserMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		if strings.Contains(r.URL.Path, utils.OnlyCurrentUserPath) {
			ctx = context.WithValue(ctx, utils.OnlyCurrentUserContextKey, true)
		} else {
			ctx = context.WithValue(ctx, utils.OnlyCurrentUserContextKey, false)
		}

		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
