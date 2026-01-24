package middlewares

import (
	"fmt"
	"net/http"

	"github.com/filipio/athletics-backend/utils"
	"github.com/go-playground/validator/v10"
)

func ErrorsMiddleware(next utils.HandlerWithError) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := next(w, r); err != nil {
			httpStatus, httpError := errorResponse(err)
			utils.Encode(w, r, httpStatus, httpError)
		}
	})
}

// based on error type, builds map for response
func errorResponse(err error) (httpStatus int, errorResponse utils.ErrorsResponse) {
	if err, ok := err.(utils.AppValidationError); ok {
		return http.StatusBadRequest, utils.ErrorsResponse{
			ErrorType: "validation_error",
			Details:   utils.AppValidationErrorDetails(err),
		}
	}

	if err, ok := err.(validator.ValidationErrors); ok {
		return http.StatusBadRequest, utils.ErrorsResponse{
			ErrorType: "validation_error",
			Details:   utils.ValidationErrorDetails(err),
		}
	}

	if _, ok := err.(utils.RecordNotFoundError); ok {
		return http.StatusNotFound, utils.ErrorsResponse{
			ErrorType: "not_found",
			Details:   "record not found",
		}
	}

	if _, ok := err.(utils.LoginError); ok {
		return http.StatusUnauthorized, utils.ErrorsResponse{
			ErrorType: "login_error",
			Details:   "login or password is invalid",
		}
	}

	if _, ok := err.(utils.AuthHeaderMissingError); ok {
		return http.StatusUnauthorized, utils.ErrorsResponse{
			ErrorType: "auth_error",
			Details:   "header 'Authorization' is missing",
		}
	}

	if _, ok := err.(utils.JwtTokenExpiredError); ok {
		return http.StatusUnauthorized, utils.ErrorsResponse{
			ErrorType: "auth_error",
			Details:   "auth token is expired",
		}
	}

	if _, ok := err.(utils.InvalidAuthHeaderError); ok {
		return http.StatusUnauthorized, utils.ErrorsResponse{
			ErrorType: "auth_error",
			Details:   "auth header is not in format 'Bearer token'",
		}
	}

	if _, ok := err.(utils.UserNotFoundError); ok {
		return http.StatusUnauthorized, utils.ErrorsResponse{
			ErrorType: "auth_error",
			Details:   "user associated with token does not exist",
		}
	}

	if _, ok := err.(utils.ActionForbiddenError); ok {
		return http.StatusUnauthorized, utils.ErrorsResponse{
			ErrorType: "auth_error",
			Details:   "user is unauthorized to perform this action",
		}
	}

	if _, ok := err.(utils.InvalidJwtClaimsError); ok {
		return http.StatusUnauthorized, utils.ErrorsResponse{
			ErrorType: "auth_error",
			Details:   "jwt claims are invalid",
		}
	}

	if _, ok := err.(utils.JwtTokenParsingError); ok {
		return http.StatusUnauthorized, utils.ErrorsResponse{
			ErrorType: "auth_error",
			Details:   err.Error(),
		}
	}

	if _, ok := err.(utils.InvalidUserError); ok {
		return http.StatusUnauthorized, utils.ErrorsResponse{
			ErrorType: "auth_error",
			Details:   "user_id must be the same as in Bearer token",
		}
	}

	if _, ok := err.(utils.EmailAlreadyExistsError); ok {
		return http.StatusBadRequest, utils.ErrorsResponse{
			ErrorType: "validation_error",
			Details:   "email already registered",
		}
	}

	if err, ok := err.(utils.EmailVerificationRateLimitError); ok {
		details := "too many registration attempts, please try again later"
		if err.BlockedUntil != nil {
			details = fmt.Sprintf("too many attempts, blocked until %s", *err.BlockedUntil)
		}
		return http.StatusTooManyRequests, utils.ErrorsResponse{
			ErrorType: "rate_limit_error",
			Details:   details,
		}
	}

	if _, ok := err.(utils.InvalidVerificationTokenError); ok {
		return http.StatusBadRequest, utils.ErrorsResponse{
			ErrorType: "validation_error",
			Details:   "invalid or expired verification token",
		}
	}

	if err, ok := err.(utils.EmailSendError); ok {
		return http.StatusInternalServerError, utils.ErrorsResponse{
			ErrorType: "email_error",
			Details:   fmt.Sprintf("failed to send email: %v", err.OriginalError),
		}
	}

	return http.StatusInternalServerError, utils.ErrorsResponse{
		ErrorType: "internal_server_error",
		Details:   err.Error(),
	}
}
