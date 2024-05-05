package utils

var RecordNotFoundResponse = ErrorsResponse{
	Errors: AnyMap{
		"invalid_request": NotFoundError,
	},
}
var HashError = ErrorsResponse{
	Errors: AnyMap{
		"internal_error": "error hashing password",
	},
}

var AuthError = ErrorsResponse{
	Errors: AnyMap{
		"invalid_request": "invalid email or password",
	},
}

const NotFoundError = "record not found"
const AdminRole = "admin"
const UserRole = "user"
const UserContextKey = 0
