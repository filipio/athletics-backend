package middlewares

import (
	"context"
	"net/http"

	"github.com/filipio/athletics-backend/utils"
	"gorm.io/gorm"
)

func DbMiddleware(next http.Handler, db *gorm.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		dbCtx := context.WithValue(ctx, utils.DbContextKey, db)
		r = r.WithContext(dbCtx)
		next.ServeHTTP(w, r)
	})
}
