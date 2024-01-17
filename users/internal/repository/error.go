package repository

import "errors"

var ErrUsernameIsRequired = errors.New("username is required")
var ErrEmailIsRequired = errors.New("email is required")
