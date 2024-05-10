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
			httpStatus, errorMap := errorMap(err)
			utils.EncodeError(w, r, httpStatus, errorMap)
		}
	})
}

// based on error type, builds map for response
func errorMap(err error) (httpStatus int, errorMap utils.AnyMap) {
	if err, ok := err.(validator.ValidationErrors); ok {
		return http.StatusBadRequest, utils.ValidationErrorMap(err)
	}

	if _, ok := err.(app_errors.RecordNotFoundError); ok {
		return http.StatusBadRequest, utils.AnyMap{"bad_request": "record not found"}
	}

	if _, ok := err.(app_errors.LoginError); ok {
		return http.StatusUnauthorized, utils.AnyMap{"login_error": "login or password is invalid"}
	}

	if _, ok := err.(app_errors.AuthHeaderMissingError); ok {
		return http.StatusUnauthorized, utils.AnyMap{"auth_error": "header 'Authorization' is missing"}
	}

	if _, ok := err.(app_errors.JwtTokenExpiredError); ok {
		return http.StatusUnauthorized, utils.AnyMap{"auth_error": "auth token is expired"}
	}

	if _, ok := err.(app_errors.InvalidAuthHeaderError); ok {
		return http.StatusUnauthorized, utils.AnyMap{"auth_error": "auth header is not in format 'Bearer token'"}
	}

	if _, ok := err.(app_errors.UserNotFoundError); ok {
		return http.StatusUnauthorized, utils.AnyMap{"auth_error": "user associated with token does not exist"}
	}

	if _, ok := err.(app_errors.ActionForbiddenError); ok {
		return http.StatusUnauthorized, utils.AnyMap{"auth_error": "user is unauthorized to perform this action"}
	}

	if _, ok := err.(app_errors.InvalidJwtClaimsError); ok {
		return http.StatusUnauthorized, utils.AnyMap{"auth_error": "jwt claims are invalid"}
	}

	return http.StatusInternalServerError, utils.AnyMap{"internal_server_error": err.Error()}
}
