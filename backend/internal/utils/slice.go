package utils

import "fmt"

func Remove[T any](s []T, i int) ([]T, error) {
	if i < 0 || i >= len(s) {
		return nil, fmt.Errorf("index %d out of range", i)
	}

	s[i] = s[len(s)-1]
	return s[:len(s)-1], nil
}

func RemoveOrdered[T any](s []T, i int) ([]T, error) {
	if i < 0 || i >= len(s) {
		return nil, fmt.Errorf("index %d out of range", i)
	}

	copy(s[i:], s[i+1:])
	return s[:len(s)-1], nil
}
