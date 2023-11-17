package dbErrors

import "errors"

var (
	ErrKeyNotFound  = errors.New("key not found")
	ErrFileNotFound = errors.New("file not found")
)
