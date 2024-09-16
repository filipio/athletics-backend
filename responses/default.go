package responses

type BuildResponseFunc[T any, V any] func(T) V

func DefaultResponse[T any](model T) T {
	return model
}
