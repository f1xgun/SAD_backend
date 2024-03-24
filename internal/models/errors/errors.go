package errorsModels

import "errors"

var ErrUserExists = errors.New("user exists")
var ErrUserNotFound = errors.New("user not found")
var ErrInvalidCredentials = errors.New("invalid credentials")
var ErrServer = errors.New("server error")
var ErrNoPermission = errors.New("no permission")
var ErrChangeOwnRole = errors.New("cannot change own role")
