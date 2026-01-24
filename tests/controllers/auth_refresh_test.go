package controllers

import (
	"net/http"
	"testing"

	"github.com/filipio/athletics-backend/utils"
)

// TestRefreshToken_Success tests that a valid refresh token returns new tokens
func TestRefreshToken_Success(t *testing.T) {
	t.Run("refresh with valid token returns new tokens", testCase(func(t *testing.T) {
		// First, login to get tokens
		email, password := getAdminCredentials()
		loginPayload := utils.AnyMap{
			"email":    email,
			"password": password,
		}
		_, loginResult, err := Post[map[string]any]("/api/v1/login", loginPayload)
		if err != nil {
			t.Fatalf("Failed to login: %v", err)
		}

		refreshToken, ok := (*loginResult)["refresh_token"].(string)
		if !ok || refreshToken == "" {
			t.Fatalf("No refresh_token in login response")
		}

		// Now use the refresh token to get new tokens
		refreshPayload := utils.AnyMap{
			"refresh_token": refreshToken,
		}
		response, refreshResult, err := Post[map[string]any]("/api/v1/auth/refresh", refreshPayload)
		if err != nil {
			t.Fatalf("Failed to refresh: %v", err)
		}

		if response.StatusCode != http.StatusOK {
			t.Fatalf("Expected status 200, got %d", response.StatusCode)
		}

		// Verify the response has the required fields
		if _, ok := (*refreshResult)["access_token"].(string); !ok {
			t.Fatal("No access_token in refresh response")
		}
		if _, ok := (*refreshResult)["refresh_token"].(string); !ok {
			t.Fatal("No refresh_token in refresh response")
		}
		if _, ok := (*refreshResult)["expires_in"].(float64); !ok {
			t.Fatal("No expires_in in refresh response")
		}
	}))
}

// TestRefreshToken_InvalidToken tests that an invalid refresh token returns 401
func TestRefreshToken_InvalidToken(t *testing.T) {
	t.Run("refresh with invalid token returns 401", testCase(func(t *testing.T) {
		refreshPayload := utils.AnyMap{
			"refresh_token": "invalid-token-12345",
		}
		response, errorResult, err := Post[map[string]any]("/api/v1/auth/refresh", refreshPayload)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}

		if response.StatusCode != http.StatusUnauthorized {
			t.Fatalf("Expected status 401, got %d", response.StatusCode)
		}

		if (*errorResult)["error_type"] != "auth_error" {
			t.Fatalf("Expected error_type to be 'auth_error', got %v", (*errorResult)["error_type"])
		}
	}))
}

// TestRefreshToken_TokenRotation tests that old refresh token becomes invalid after refresh
func TestRefreshToken_TokenRotation(t *testing.T) {
	t.Run("token rotation - old token becomes invalid", testCase(func(t *testing.T) {
		// First, login to get tokens
		email, password := getAdminCredentials()
		loginPayload := utils.AnyMap{
			"email":    email,
			"password": password,
		}
		_, loginResult, err := Post[map[string]any]("/api/v1/login", loginPayload)
		if err != nil {
			t.Fatalf("Failed to login: %v", err)
		}

		refreshToken := (*loginResult)["refresh_token"].(string)

		// Use the refresh token to get new tokens
		refreshPayload := utils.AnyMap{
			"refresh_token": refreshToken,
		}
		_, _, err = Post[map[string]any]("/api/v1/auth/refresh", refreshPayload)
		if err != nil {
			t.Fatalf("Failed to refresh: %v", err)
		}

		// Try to use the old refresh token again - should fail
		response, _, err := Post[map[string]any]("/api/v1/auth/refresh", refreshPayload)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}

		if response.StatusCode != http.StatusUnauthorized {
			t.Fatalf("Expected status 401 when using old token, got %d", response.StatusCode)
		}
	}))
}

// TestLogout_Success tests that logout successfully revokes the session
func TestLogout_Success(t *testing.T) {
	t.Run("logout successfully revokes session", testCase(func(t *testing.T) {
		// Create a new token for this test so we don't invalidate the global adminToken
		email, password := getAdminCredentials()
		loginPayload := utils.AnyMap{
			"email":    email,
			"password": password,
		}
		_, loginResult, err := Post[map[string]any]("/api/v1/login", loginPayload)
		if err != nil {
			t.Fatalf("Failed to login: %v", err)
		}

		testToken := (*loginResult)["access_token"].(string)

		// Now logout with this token
		response, logoutResult, err := executeLogout("/api/v1/auth/logout", testToken)
		if err != nil {
			t.Fatalf("Failed to logout: %v", err)
		}

		if response.StatusCode != http.StatusOK {
			t.Fatalf("Expected status 200, got %d", response.StatusCode)
		}

		if (*logoutResult)["message"] != "logged out successfully" {
			t.Fatalf("Expected logout message, got %v", (*logoutResult)["message"])
		}
	}))
}

// TestLogout_InvalidatesAccessToken tests that access token fails after logout
func TestLogout_InvalidatesAccessToken(t *testing.T) {
	t.Run("access token fails after logout", testCase(func(t *testing.T) {
		// Login to get a new access token
		email, password := getAdminCredentials()
		loginPayload := utils.AnyMap{
			"email":    email,
			"password": password,
		}
		_, loginResult, err := Post[map[string]any]("/api/v1/login", loginPayload)
		if err != nil {
			t.Fatalf("Failed to login: %v", err)
		}

		accessToken := (*loginResult)["access_token"].(string)

		// Verify the access token works first
		response, _, err := executeHttpWithToken[map[string]any]("GET", "/api/v1/users/me", nil, accessToken)
		if err != nil {
			t.Fatalf("Failed to get user: %v", err)
		}
		if response.StatusCode != http.StatusOK {
			t.Fatalf("Expected status 200 for authenticated request, got %d", response.StatusCode)
		}

		// Now logout using the new token
		_, _, err = executeLogout("/api/v1/auth/logout", accessToken)
		if err != nil {
			t.Fatalf("Failed to logout: %v", err)
		}

		// Try to use the access token again - should fail
		response, _, err = executeHttpWithToken[map[string]any]("GET", "/api/v1/users/me", nil, accessToken)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}

		if response.StatusCode != http.StatusUnauthorized {
			t.Fatalf("Expected status 401 after logout, got %d", response.StatusCode)
		}
	}))
}

// TestLogin_ReturnsNewTokenPair tests successful login
func TestLogin_ReturnsNewTokenPair(t *testing.T) {
	t.Run("login returns new format (access_token, refresh_token, expires_in)", testCase(func(t *testing.T) {
		email, password := getAdminCredentials()
		loginPayload := utils.AnyMap{
			"email":    email,
			"password": password,
		}
		response, loginResult, err := Post[map[string]any]("/api/v1/login", loginPayload)
		if err != nil {
			t.Fatalf("Failed to login: %v", err)
		}

		if response.StatusCode != http.StatusOK {
			t.Fatalf("Expected status 200, got %d", response.StatusCode)
		}

		// Check for new fields
		if _, ok := (*loginResult)["access_token"].(string); !ok {
			t.Fatal("No access_token in login response")
		}
		if _, ok := (*loginResult)["refresh_token"].(string); !ok {
			t.Fatal("No refresh_token in login response")
		}
		expiresIn, ok := (*loginResult)["expires_in"].(float64)
		if !ok {
			t.Fatal("No expires_in in login response")
		}

		if expiresIn != 3600 {
			t.Fatalf("Expected expires_in to be 3600, got %v", expiresIn)
		}
	}))
}

// TestTokenLifecycle_FullFlow tests the complete lifecycle: Login -> Use -> Refresh -> Logout
func TestTokenLifecycle_FullFlow(t *testing.T) {
	t.Run("full token lifecycle", testCase(func(t *testing.T) {
		// Step 1: Login
		email, password := getAdminCredentials()
		loginPayload := utils.AnyMap{
			"email":    email,
			"password": password,
		}
		_, loginResult, err := Post[map[string]any]("/api/v1/login", loginPayload)
		if err != nil {
			t.Fatalf("Failed to login: %v", err)
		}

		accessToken1 := (*loginResult)["access_token"].(string)
		refreshToken1 := (*loginResult)["refresh_token"].(string)

		// Step 2: Use access token to make an API call
		response, getUserResult, err := executeHttpWithToken[map[string]any]("GET", "/api/v1/users/me", nil, accessToken1)
		if err != nil {
			t.Fatalf("Failed to get user: %v", err)
		}
		if response.StatusCode != http.StatusOK || (*getUserResult)["id"] == nil {
			t.Fatal("Failed to use access token for API call")
		}

		// Step 3: Refresh the token
		refreshPayload := utils.AnyMap{
			"refresh_token": refreshToken1,
		}
		_, refreshResult, err := Post[map[string]any]("/api/v1/auth/refresh", refreshPayload)
		if err != nil {
			t.Fatalf("Failed to refresh: %v", err)
		}

		accessToken2 := (*refreshResult)["access_token"].(string)
		refreshToken2 := (*refreshResult)["refresh_token"].(string)

		// Verify new tokens work
		response, getUserResult, err = executeHttpWithToken[map[string]any]("GET", "/api/v1/users/me", nil, accessToken2)
		if err != nil || response.StatusCode != http.StatusOK || (*getUserResult)["id"] == nil {
			t.Fatal("Failed to use new access token")
		}

		// Verify old refresh token no longer works
		response, _, _ = Post[map[string]any]("/api/v1/auth/refresh", utils.AnyMap{
			"refresh_token": refreshToken1,
		})
		if response.StatusCode != http.StatusUnauthorized {
			t.Fatal("Old refresh token should not work")
		}

		// Step 4: Logout
		_, _, err = executeLogout("/api/v1/auth/logout", accessToken2)
		if err != nil {
			t.Fatalf("Failed to logout: %v", err)
		}

		// Verify new access token no longer works
		response, _, _ = executeHttpWithToken[map[string]any]("GET", "/api/v1/users/me", nil, accessToken2)
		if response.StatusCode != http.StatusUnauthorized {
			t.Fatal("Access token should not work after logout")
		}

		// Verify new refresh token no longer works
		response, _, _ = Post[map[string]any]("/api/v1/auth/refresh", utils.AnyMap{
			"refresh_token": refreshToken2,
		})
		if response.StatusCode != http.StatusUnauthorized {
			t.Fatal("Refresh token should not work after logout")
		}
	}))
}

// TestAccessToken_SessionValidation tests that access tokens require valid sessions
func TestAccessToken_SessionValidation(t *testing.T) {
	t.Run("access token validates session on each request", testCase(func(t *testing.T) {
		// Login to get tokens
		email, password := getAdminCredentials()
		loginPayload := utils.AnyMap{
			"email":    email,
			"password": password,
		}
		_, loginResult, err := Post[map[string]any]("/api/v1/login", loginPayload)
		if err != nil {
			t.Fatalf("Failed to login: %v", err)
		}

		accessToken := (*loginResult)["access_token"].(string)

		// Verify the access token works
		response, _, err := executeHttpWithToken[map[string]any]("GET", "/api/v1/users/me", nil, accessToken)
		if err != nil || response.StatusCode != http.StatusOK {
			t.Fatal("Access token should work initially")
		}

		// Session should remain valid after multiple requests
		for i := 0; i < 3; i++ {
			response, _, err := executeHttpWithToken[map[string]any]("GET", "/api/v1/users/me", nil, accessToken)
			if err != nil || response.StatusCode != http.StatusOK {
				t.Fatalf("Access token should work on request %d", i+1)
			}
		}
	}))
}
