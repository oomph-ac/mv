package pointer

// Make returns a pointer of passed value.
func Make[T any](v T) *T {
	return &v
}
