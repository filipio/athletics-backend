package middlewares

import (
	"context"
	"errors"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/filipio/athletics-backend/models"
	"github.com/filipio/athletics-backend/utils"
	"github.com/filipio/athletics-backend/utils/app_errors"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

func AdminOnly(next utils.HandlerWithError, db *gorm.DB) utils.HandlerWithError {
	return authMiddleware(next, utils.AdminRole, db)
}

func UserOnly(next utils.HandlerWithError, db *gorm.DB) utils.HandlerWithError {
	return authMiddleware(next, utils.UserRole, db)
}

func authMiddleware(next utils.HandlerWithError, allowedRole string, db *gorm.DB) utils.HandlerWithError {
	return utils.HandlerWithError(func(w http.ResponseWriter, r *http.Request) error {
		tokenString, extractionError := extractToken(r)

		if extractionError != nil {
			return extractionError
		}

		token, parsingError := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected jwt signing method")
			}

			return []byte(os.Getenv("JWT_SIGNING_SECRET")), nil
		})

		if parsingError != nil {
			return app_errors.JwtTokenParsingError{AppError: app_errors.AppError{parsingError.Error()}}
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			if float64(time.Now().Unix()) > claims["exp"].(float64) {
				return app_errors.JwtTokenExpiredError{}
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
					return app_errors.UserNotFoundError{}
				}

				ctx := context.WithValue(r.Context(), utils.UserContextKey, user)
				return next.ServeHTTP(w, r.WithContext(ctx))
			} else {
				return app_errors.ActionForbiddenError{}
			}

		} else {
			return app_errors.InvalidJwtClaimsError{}
		}
	})
}

func extractToken(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", app_errors.AuthHeaderMissingError{}
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return "", app_errors.InvalidAuthHeaderError{}
	}

	return parts[1], nil
}
