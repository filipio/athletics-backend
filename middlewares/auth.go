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
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

func AdminOnly(next utils.HandlerWithError) utils.HandlerWithError {
	return authMiddleware(next, utils.AdminRole)
}

func UserOnly(next utils.HandlerWithError) utils.HandlerWithError {
	return authMiddleware(next, utils.UserRole)
}

func OrganizerOnly(next utils.HandlerWithError) utils.HandlerWithError {
	return authMiddleware(next, utils.OrganizerRole)
}

func authMiddleware(next utils.HandlerWithError, requiredRole string) utils.HandlerWithError {
	return utils.HandlerWithError(func(w http.ResponseWriter, r *http.Request) error {
		tokenString, extractionError := extractToken(r)

		if extractionError != nil {
			return extractionError
		}

		token, parsingError := parseToken(tokenString)

		if parsingError != nil {
			return utils.JwtTokenParsingError{AppError: utils.AppError{Message: parsingError.Error()}}
		}

		claims, ok := token.Claims.(jwt.MapClaims)

		if !ok || !token.Valid {
			return utils.InvalidJwtClaimsError{}
		}

		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			return utils.JwtTokenExpiredError{}
		}

		if !requiredRoleFound(claims, requiredRole) {
			return utils.ActionForbiddenError{}
		}

		clientContext, err := buildClientContext(r, claims, models.Db(r))
		if err != nil {
			return err
		}

		return next.ServeHTTP(w, r.WithContext(clientContext))
	})
}

func extractToken(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", utils.AuthHeaderMissingError{}
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return "", utils.InvalidAuthHeaderError{}
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
		return nil, utils.UserNotFoundError{}
	}

	return context.WithValue(r.Context(), utils.UserContextKey, user), nil
}
