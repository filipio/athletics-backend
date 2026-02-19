package middleware

import (
	"fmt"
	"net/http"

	"github.com/filipio/athletics-backend/pkg/httpio"
	"github.com/go-playground/validator/v10"
)

func ErrorsMiddleware(next httpio.HandlerWithError) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := next(w, r); err != nil {
			httpStatus, httpError := errorResponse(err)
			httpio.Encode(w, r, httpStatus, httpError)
		}
	})
}

// based on error type, builds map for response
func errorResponse(err error) (httpStatus int, errorResponse httpio.ErrorsResponse) {
	if err, ok := err.(httpio.AppValidationError); ok {
		return http.StatusBadRequest, httpio.ErrorsResponse{
			ErrorType: "validation_error",
			Details:   httpio.AppValidationErrorDetails(err),
		}
	}

	if err, ok := err.(validator.ValidationErrors); ok {
		return http.StatusBadRequest, httpio.ErrorsResponse{
			ErrorType: "validation_error",
			Details:   httpio.ValidationErrorDetails(err),
		}
	}

	if _, ok := err.(httpio.RecordNotFoundError); ok {
		return http.StatusNotFound, httpio.ErrorsResponse{
			ErrorType: "not_found",
			Details:   "record not found",
		}
	}

	if _, ok := err.(httpio.LoginError); ok {
		return http.StatusUnauthorized, httpio.ErrorsResponse{
			ErrorType: "login_error",
			Details:   "login or password is invalid",
		}
	}

	if _, ok := err.(httpio.AuthHeaderMissingError); ok {
		return http.StatusUnauthorized, httpio.ErrorsResponse{
			ErrorType: "auth_error",
			Details:   "header 'Authorization' is missing",
		}
	}

	if _, ok := err.(httpio.JwtTokenExpiredError); ok {
		return http.StatusUnauthorized, httpio.ErrorsResponse{
			ErrorType: "auth_error",
			Details:   "auth token is expired",
		}
	}

	if _, ok := err.(httpio.InvalidAuthHeaderError); ok {
		return http.StatusUnauthorized, httpio.ErrorsResponse{
			ErrorType: "auth_error",
			Details:   "auth header is not in format 'Bearer token'",
		}
	}

	if _, ok := err.(httpio.UserNotFoundError); ok {
		return http.StatusUnauthorized, httpio.ErrorsResponse{
			ErrorType: "auth_error",
			Details:   "user associated with token does not exist",
		}
	}

	if _, ok := err.(httpio.ActionForbiddenError); ok {
		return http.StatusUnauthorized, httpio.ErrorsResponse{
			ErrorType: "auth_error",
			Details:   "user is unauthorized to perform this action",
		}
	}

	if _, ok := err.(httpio.InvalidJwtClaimsError); ok {
		return http.StatusUnauthorized, httpio.ErrorsResponse{
			ErrorType: "auth_error",
			Details:   "jwt claims are invalid",
		}
	}

	if _, ok := err.(httpio.JwtTokenParsingError); ok {
		return http.StatusUnauthorized, httpio.ErrorsResponse{
			ErrorType: "auth_error",
			Details:   err.Error(),
		}
	}

	if _, ok := err.(httpio.InvalidUserError); ok {
		return http.StatusUnauthorized, httpio.ErrorsResponse{
			ErrorType: "auth_error",
			Details:   "user_id must be the same as in Bearer token",
		}
	}

	if _, ok := err.(httpio.EmailAlreadyExistsError); ok {
		return http.StatusBadRequest, httpio.ErrorsResponse{
			ErrorType: "validation_error",
			Details:   "email already registered",
		}
	}

	if err, ok := err.(httpio.EmailVerificationRateLimitError); ok {
		details := "too many registration attempts, please try again later"
		if err.BlockedUntil != nil {
			details = fmt.Sprintf("too many attempts, blocked until %s", *err.BlockedUntil)
		}
		return http.StatusTooManyRequests, httpio.ErrorsResponse{
			ErrorType: "rate_limit_error",
			Details:   details,
		}
	}

	if _, ok := err.(httpio.InvalidVerificationTokenError); ok {
		return http.StatusBadRequest, httpio.ErrorsResponse{
			ErrorType: "validation_error",
			Details:   "invalid or expired verification token",
		}
	}

	if err, ok := err.(httpio.EmailSendError); ok {
		return http.StatusInternalServerError, httpio.ErrorsResponse{
			ErrorType: "email_error",
			Details:   fmt.Sprintf("failed to send email: %v", err.OriginalError),
		}
	}

	if _, ok := err.(httpio.InvalidRefreshTokenError); ok {
		return http.StatusUnauthorized, httpio.ErrorsResponse{
			ErrorType: "auth_error",
			Details:   "invalid or expired refresh token",
		}
	}

	if _, ok := err.(httpio.RefreshTokenExpiredError); ok {
		return http.StatusUnauthorized, httpio.ErrorsResponse{
			ErrorType: "auth_error",
			Details:   "refresh token has expired",
		}
	}

	if _, ok := err.(httpio.SessionExpiredError); ok {
		return http.StatusUnauthorized, httpio.ErrorsResponse{
			ErrorType: "auth_error",
			Details:   "session has expired",
		}
	}

	return http.StatusInternalServerError, httpio.ErrorsResponse{
		ErrorType: "internal_server_error",
		Details:   err.Error(),
	}
}
