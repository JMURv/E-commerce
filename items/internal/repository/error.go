package repository

import "errors"

var ErrNotFound = errors.New("not found")

var ErrUserIDRequired = errors.New("userID is required")
var ErrCategoryIDRequired = errors.New("categoryID is required")
var ErrNoSuchCategory = errors.New("no such category")
var ErrNameRequired = errors.New("name is required")
var ErrDescriptionRequired = errors.New("description is required")
var ErrPriceRequired = errors.New("price is required")

var ErrCategoryNameIsRequired = errors.New("category name is required")
var ErrCategoryDescriptionIsRequired = errors.New("category description is required")

var ErrTagNameIsRequired = errors.New("tag name is required")
