package middleware

import (
	"context"
	"errors"
	"net/http"
	"os"
	"slices"
	"strings"
	"time"

	"github.com/filipio/athletics-backend/pkg/config"
	"github.com/filipio/athletics-backend/internal/models"
	"github.com/filipio/athletics-backend/pkg/httpio"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

type AuthMiddleware struct {
	deps *config.Dependencies
}

func NewAuthMiddleware(deps *config.Dependencies) *AuthMiddleware {
	return &AuthMiddleware{deps: deps}
}

func (a *AuthMiddleware) AdminOnly(next httpio.HandlerWithError) httpio.HandlerWithError {
	return a.authMiddleware(next, httpio.AdminRole)
}

func (a *AuthMiddleware) UserOnly(next httpio.HandlerWithError) httpio.HandlerWithError {
	return a.authMiddleware(next, httpio.UserRole)
}

func (a *AuthMiddleware) OrganizerOnly(next httpio.HandlerWithError) httpio.HandlerWithError {
	return a.authMiddleware(next, httpio.OrganizerRole)
}

func (a *AuthMiddleware) authMiddleware(next httpio.HandlerWithError, requiredRole string) httpio.HandlerWithError {
	return httpio.HandlerWithError(func(w http.ResponseWriter, r *http.Request) error {
		tokenString, extractionError := extractToken(r)
		if extractionError != nil {
			return extractionError
		}

		token, parsingError := parseToken(tokenString)
		if parsingError != nil {
			return httpio.JwtTokenParsingError{AppError: httpio.AppError{Message: parsingError.Error()}}
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			return httpio.InvalidJwtClaimsError{}
		}

		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			return httpio.JwtTokenExpiredError{}
		}

		if !requiredRoleFound(claims, requiredRole) {
			return httpio.ActionForbiddenError{}
		}

		// Validate session if session_id is present in claims
		sessionID, hasSessionID := claims["session_id"].(string)
		if hasSessionID && sessionID != "" {
			db := a.deps.DB
			var refreshToken models.RefreshToken
			err := db.Where("session_id = ? AND revoked_at IS NULL AND expires_at > ?",
				sessionID, time.Now()).First(&refreshToken).Error
			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return httpio.SessionExpiredError{}
				}
				return err
			}
		}

		clientContext, err := a.buildClientContext(r, claims)
		if err != nil {
			return err
		}

		// Add session_id to context for logout functionality
		if sessionID != "" {
			clientContext = context.WithValue(clientContext, httpio.SessionIDContextKey, sessionID)
		}

		return next.ServeHTTP(w, r.WithContext(clientContext))
	})
}

func (a *AuthMiddleware) buildClientContext(r *http.Request, claims jwt.MapClaims) (context.Context, error) {
	userID := claims["sub"]

	var user models.User
	a.deps.DB.Preload("Roles").First(&user, userID)
	if user.ID == 0 {
		return nil, httpio.UserNotFoundError{}
	}

	return context.WithValue(r.Context(), httpio.UserContextKey, user), nil
}

func extractToken(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", httpio.AuthHeaderMissingError{}
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return "", httpio.InvalidAuthHeaderError{}
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
	userRolesInterface, ok := claims["roles"].([]interface{})
	if !ok || userRolesInterface == nil {
		return false
	}

	userRoles := make([]string, len(userRolesInterface))
	for i, v := range userRolesInterface {
		userRoles[i] = v.(string)
	}

	if slices.Contains(userRoles, httpio.AdminRole) {
		return true
	}

	return slices.Contains(userRoles, requiredRole)
}
