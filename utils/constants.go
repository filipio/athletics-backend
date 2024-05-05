package utils

var RecordNotFoundResponse = ErrorsResponse{
	Errors: AnyMap{
		"invalid_request": NotFoundError,
	},
}

const NotFoundError = "record not found"
