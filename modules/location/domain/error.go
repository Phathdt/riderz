package domain

import "errors"

var (
	ErrCannotUpdateLocation     = errors.New("cannot update location")
	ErrCannotSearchNearLocation = errors.New("cannot search near location")
)
