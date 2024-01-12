package repository

import "errors"

var ErrNotFound = errors.New("not found")

var ErrUserIDRequired = errors.New("userID is required")
var ErrItemIDRequired = errors.New("itemID or ReviewedUserID is required")
var ErrRatingRequired = errors.New("rating is required")
