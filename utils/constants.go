package utils

var RecordNotFoundResponse = ErrorsResponse{
	Errors: map[string]interface{}{
		"invalid_request": "record not found",
	},
}
