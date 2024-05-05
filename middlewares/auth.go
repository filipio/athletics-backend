package middlewares

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/filipio/athletics-backend/models"
	"github.com/filipio/athletics-backend/utils"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

func AdminOnly(next http.Handler, db *gorm.DB) http.Handler {
	return authMiddleware(next, utils.AdminRole, db)
}

func UserOnly(next http.Handler, db *gorm.DB) http.Handler {
	return authMiddleware(next, utils.UserRole, db)
}

func authMiddleware(next http.Handler, allowedRole string, db *gorm.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString, extractionError := extractToken(r)

		if extractionError != nil {
			utils.BadRequest(w, r, extractionError.Error())
			return
		}

		token, parsingError := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			return []byte(os.Getenv("JWT_SIGNING_SECRET")), nil
		})

		if parsingError != nil {
			utils.BadRequest(w, r, parsingError.Error())
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			// Chec k the expiry date
			if float64(time.Now().Unix()) > claims["exp"].(float64) {
				utils.Unauthorized(w, r, "expired token")
				return
			}

			userRoles := claims["roles"].([]interface{})

			roleMatch := false
			for _, role := range userRoles {
				stringRole := role.(string)
				roleMatch = (stringRole == allowedRole)
			}

			if roleMatch {
				userID := claims["sub"]

				var user models.User
				db.First(&user, userID)
				if user.ID == 0 {
					utils.BadRequest(w, r, "user does not exist")
					return
				}

				ctx := context.WithValue(r.Context(), utils.UserContextKey, user)

				next.ServeHTTP(w, r.WithContext(ctx))
			} else {
				utils.Forbidden(w, r, "action forbidden")
				return
			}

		} else {
			utils.Unauthorized(w, r, "invalid jwt token claims")
			return
		}
	})
}

func extractToken(r *http.Request) (string, error) {
	// Extract the token from the Authorization header
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", fmt.Errorf("authorization header is missing")
	}

	// Bearer token format check
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return "", fmt.Errorf("authorization header format must be Bearer {token}")
	}

	return parts[1], nil
}
