package middlewares

import (
	"net/http"

	"github.com/filipio/athletics-backend/utils"
	"github.com/filipio/athletics-backend/utils/app_errors"
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
	if err, ok := err.(validator.ValidationErrors); ok {
		return http.StatusBadRequest, utils.ErrorsResponse{
			ErrorType: "validation_error",
			Details:   utils.ValidationErrorDetails(err),
		}
	}

	if _, ok := err.(app_errors.RecordNotFoundError); ok {
		return http.StatusBadRequest, utils.ErrorsResponse{
			ErrorType: "bad_request",
			Details:   "record not found",
		}
	}

	if _, ok := err.(app_errors.LoginError); ok {
		return http.StatusUnauthorized, utils.ErrorsResponse{
			ErrorType: "login_error",
			Details:   "login or password is invalid",
		}
	}

	if _, ok := err.(app_errors.AuthHeaderMissingError); ok {
		return http.StatusUnauthorized, utils.ErrorsResponse{
			ErrorType: "auth_error",
			Details:   "header 'Authorization' is missing",
		}
	}

	if _, ok := err.(app_errors.JwtTokenExpiredError); ok {
		return http.StatusUnauthorized, utils.ErrorsResponse{
			ErrorType: "auth_error",
			Details:   "auth token is expired",
		}
	}

	if _, ok := err.(app_errors.InvalidAuthHeaderError); ok {
		return http.StatusUnauthorized, utils.ErrorsResponse{
			ErrorType: "auth_error",
			Details:   "auth header is not in format 'Bearer token'",
		}
	}

	if _, ok := err.(app_errors.UserNotFoundError); ok {
		return http.StatusUnauthorized, utils.ErrorsResponse{
			ErrorType: "auth_error",
			Details:   "user associated with token does not exist",
		}
	}

	if _, ok := err.(app_errors.ActionForbiddenError); ok {
		return http.StatusUnauthorized, utils.ErrorsResponse{
			ErrorType: "auth_error",
			Details:   "user is unauthorized to perform this action",
		}
	}

	if _, ok := err.(app_errors.InvalidJwtClaimsError); ok {
		return http.StatusUnauthorized, utils.ErrorsResponse{
			ErrorType: "auth_error",
			Details:   "jwt claims are invalid",
		}
	}

	if _, ok := err.(app_errors.JwtTokenParsingError); ok {
		return http.StatusUnauthorized, utils.ErrorsResponse{
			ErrorType: "auth_error",
			Details:   err.Error(),
		}
	}

	return http.StatusInternalServerError, utils.ErrorsResponse{
		ErrorType: "internal_server_error",
		Details:   err.Error(),
	}
}
