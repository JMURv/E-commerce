package repository

import "errors"

var ErrNotFound = errors.New("not found")

var ErrUserIDRequired = errors.New("userID is required")
var ErrCategoryIDRequired = errors.New("categoryID is required")
var ErrNameRequired = errors.New("name is required")
var ErrDescriptionRequired = errors.New("description is required")
var ErrPriceRequired = errors.New("price is required")
