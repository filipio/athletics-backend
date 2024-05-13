package app_errors

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
