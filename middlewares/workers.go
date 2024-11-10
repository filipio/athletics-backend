package middlewares

import (
	"context"
	"net/http"

	"github.com/filipio/athletics-backend/config"
	"github.com/filipio/athletics-backend/utils"
)

func WorkersMiddleware(next http.Handler, insertWorkersClient *config.InsertWorkerClient) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		workersCtx := context.WithValue(ctx, utils.WorkersContextKey, insertWorkersClient)
		r = r.WithContext(workersCtx)
		next.ServeHTTP(w, r)
	})
}
