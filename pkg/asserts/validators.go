package asserts

import (
	"golang.org/x/exp/constraints"
)

// AssertInBetween checks if a value is in between two other values
func AssertInBetween[T constraints.Ordered](left T, right T, value T) bool {
	if left > value || right < value {
		return false
	}

	return true
}

// AssertContains checks if a slice contains a value
func AssertContains[T comparable](slice []T, value T) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}

	return false
}

// AssertSize checks if the size of a byte slice is smaller, equal, or bigger than the given maximum size
func AssertSize(slice []byte, maxSize int) int {
	length := len(slice)
	switch {
	case length < maxSize:
		return -1
	case length == maxSize:
		return 0
	default:
		return 1
	}
}
