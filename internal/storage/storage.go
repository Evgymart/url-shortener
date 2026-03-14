package storage

import "errors"

var (
	ErrUrlNotFound       = errors.New("url not found")
	ErrAliasAlreadyTaken = errors.New("alias already taken")
)
