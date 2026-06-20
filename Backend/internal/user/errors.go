package user

import "errors"

var (
	ErrEmailAlreadyExists   = errors.New("email already exists")
	ErrPANAlreadyExists     = errors.New("pan already exists")
	ErrAadhaarAlreadyExists = errors.New("aadhaar already exists")
	ErrUserNotFound         = errors.New("user not found")
)