package utils

import "errors"

// handler must understand which type of errors
// came from service or even from repository
// so we use custom errors in whole app

var (
	ErrIDEmpty       = errors.New("id is empty")
	ErrIDNotNumber   = errors.New("id must be a valid integer")
	ErrIDNotPositive = errors.New("id must be a positive integer")
	ErrTaskNotFound  = errors.New("task not found")
)
