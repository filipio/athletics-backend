package responses

type DefaultResponse struct {
	Data any `json:"data"`
}

type BuildResponseFunc[T any, V any] func(T) V

func BuildDefaultResponse[T any](model T) T {
	return model
}
