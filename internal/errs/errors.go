package errs

import "errors"

var (
	ErrDuplicateEmail = errors.New("user with this email already exists")
)
