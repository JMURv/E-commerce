package repository

import "errors"

var ErrNotFound = errors.New("user not found")
var ErrUsernameIsRequired = errors.New("username is required")
var ErrEmailIsRequired = errors.New("email is required")
