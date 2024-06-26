package middlewares

import (
	"context"
	"errors"
	"net/http"
	"os"
	"slices"
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

func authMiddleware(next utils.HandlerWithError, requiredRole string, db *gorm.DB) utils.HandlerWithError {
	return utils.HandlerWithError(func(w http.ResponseWriter, r *http.Request) error {
		tokenString, extractionError := extractToken(r)

		if extractionError != nil {
			return extractionError
		}

		token, parsingError := parseToken(tokenString)

		if parsingError != nil {
			return app_errors.JwtTokenParsingError{AppError: app_errors.AppError{Message: parsingError.Error()}}
		}

		claims, ok := token.Claims.(jwt.MapClaims)

		if !ok || !token.Valid {
			return app_errors.InvalidJwtClaimsError{}
		}

		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			return app_errors.JwtTokenExpiredError{}
		}

		if !requiredRoleFound(claims, requiredRole) {
			return app_errors.ActionForbiddenError{}
		}

		clientContext, err := buildClientContext(r, claims, db)
		if err != nil {
			return err
		}

		return next.ServeHTTP(w, r.WithContext(clientContext))
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

func parseToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected jwt signing method")
		}

		return []byte(os.Getenv("JWT_SIGNING_SECRET")), nil
	})
}

func requiredRoleFound(claims jwt.MapClaims, requiredRole string) bool {
	userRoles := claims["roles"].([]interface{})

	if slices.Contains(userRoles, utils.AdminRole) {
		return true
	}

	roleFound := false

	for _, role := range userRoles {
		actualRole := role.(string)
		roleFound = (actualRole == requiredRole)
	}

	return roleFound
}

func buildClientContext(r *http.Request, claims jwt.MapClaims, db *gorm.DB) (context.Context, error) {
	userID := claims["sub"]

	var user models.User
	db.First(&user, userID)
	if user.ID == 0 {
		return nil, app_errors.UserNotFoundError{}
	}

	return context.WithValue(r.Context(), utils.UserContextKey, user), nil
}
