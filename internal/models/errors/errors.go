package errorsModels

import "errors"

var (
	ErrUserExists         = errors.New("user exists")
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrServer             = errors.New("server error")
	ErrNoPermission       = errors.New("no permission")
	ErrChangeOwnRole      = errors.New("cannot change own role")
	ErrUserDoesNotExist   = errors.New("user does not exist")
	ErrGroupExists        = errors.New("group exists")
	ErrGroupDoesNotExist  = errors.New("group does not exist")
	ErrUserNotInGroup     = errors.New("user is not in the group")
)
