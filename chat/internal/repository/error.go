package repository

import "errors"

var ErrNotFound = errors.New("not found")

var ErrUserIDRequired = errors.New("userID is required")
var ErrRoomIDRequired = errors.New("roomID is required")
var ErrItemIDRequired = errors.New("itemID is required")
var ErrTextRequired = errors.New("text is required")
