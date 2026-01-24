package controllers

import (
	"os"
	"time"

	"github.com/filipio/athletics-backend/models"
	"github.com/filipio/athletics-backend/utils"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
}

// generateAccessToken creates a short-lived JWT access token with session_id in claims
func generateAccessToken(user models.User, roleNames []string, sessionID string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":        user.ID,
		"exp":        time.Now().Add(utils.AccessTokenExpiration).Unix(),
		"roles":      roleNames,
		"username":   user.Username,
		"session_id": sessionID,
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SIGNING_SECRET")))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// createRefreshToken creates a long-lived refresh token and stores it in the database
func createRefreshToken(user models.User, db *gorm.DB, sessionID string) (string, error) {
	refreshTokenModel := models.RefreshToken{
		UserID:    user.ID,
		SessionID: sessionID,
		ExpiresAt: time.Now().Add(utils.RefreshTokenExpiration),
	}

	plainToken := refreshTokenModel.GenerateToken()

	hashedToken, err := refreshTokenModel.HashToken(plainToken)
	if err != nil {
		return "", err
	}

	refreshTokenModel.TokenHash = hashedToken

	if err := db.Create(&refreshTokenModel).Error; err != nil {
		return "", err
	}

	return plainToken, nil
}

// generateTokenPair creates both access and refresh tokens with the same session_id (refresh token is stored in DB)
// If sessionID is nil, a new one will be generated
func generateTokenPair(user models.User, db *gorm.DB, sessionID *string) (TokenPair, error) {
	var id string
	if sessionID == nil {
		id = (&models.RefreshToken{}).GenerateSessionID()
	} else {
		id = *sessionID
	}

	var roles []models.Role
	if err := db.Model(&user).Association("Roles").Find(&roles); err != nil {
		return TokenPair{}, err
	}

	roleNames := make([]string, len(roles))
	for i, role := range roles {
		roleNames[i] = role.Name
	}

	accessToken, err := generateAccessToken(user, roleNames, id)
	if err != nil {
		return TokenPair{}, err
	}

	refreshToken, err := createRefreshToken(user, db, id)
	if err != nil {
		return TokenPair{}, err
	}

	return TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(utils.AccessTokenExpiration.Seconds()),
	}, nil
}
