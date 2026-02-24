package repository

import "errors"

var (
	ErrNotFound          = errors.New("not found")
	ErrConflict          = errors.New("conflict")
	ErrResignIDMustBeSet = errors.New("resign id must be set")
	ErrIDMustNotReferToItself = errors.New("id must not refer to itself")
)
