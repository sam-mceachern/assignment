package util

func ToPtr[T any](input T) *T {
	return &input
}
