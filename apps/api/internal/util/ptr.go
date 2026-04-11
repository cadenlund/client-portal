package util

// Function for putting any type in a generic variable and return its address
func Ptr[T any](v T) *T { return &v }
