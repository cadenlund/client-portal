// Author: Caden Lund
// Created: 4/11/2026
// Last updated: 4/11/2026
// Notes:

package util

// Function for putting any type in a generic variable and return its address
func Ptr[T any](v T) *T { return &v }
