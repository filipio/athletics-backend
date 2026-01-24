package utils

type AppError struct {
	Message string
}

func (e AppError) Error() string {
	return e.Message
}

type InternalError struct {
	AppError
}

type LoginError struct {
	AppError
}

type RecordNotFoundError struct {
	AppError
}

type AuthHeaderMissingError struct {
	AppError
}

type JwtTokenExpiredError struct {
	AppError
}

type InvalidAuthHeaderError struct {
	AppError
}

type UserNotFoundError struct {
	AppError
}

type ActionForbiddenError struct {
	AppError
}

type InvalidJwtClaimsError struct {
	AppError
}

type JwtTokenParsingError struct {
	AppError
}

type InvalidUserError struct {
	AppError
}

type AppValidationError struct {
	AppError
	FieldPath string
}

type EmailAlreadyExistsError struct {
	AppError
}

type EmailVerificationRateLimitError struct {
	AppError
	BlockedUntil *string
}

type InvalidVerificationTokenError struct {
	AppError
}

type EmailSendError struct {
	AppError
	OriginalError error
}

type InvalidRefreshTokenError struct {
	AppError
}

type RefreshTokenExpiredError struct {
	AppError
}

type SessionExpiredError struct {
	AppError
}
