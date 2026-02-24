package domain

import "errors"

var (
	ErrEmptyDepartmentName    = errors.New("empty department name")
	ErrToLongDepartmentName   = errors.New("too long department name")
	ErrEmptyFullName          = errors.New("empty full name")
	ErrToLongFullName         = errors.New("to long full name")
	ErrEmptyPositionName      = errors.New("empty position name")
	ErrToLongPositionName     = errors.New("to long position name")
	ErrNotFound               = errors.New("not found")
	ErrConflict               = errors.New("conflict")
	ErrIDMustNotReferToItself = errors.New("id must not refer to itself")
	ErrResignIDMustBeSet      = errors.New("resign id must be set")
	ErrGatewayTimeout         = errors.New("gateway timeout")
)
