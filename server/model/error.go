package model

import "errors"

var (
	ERROR_USER_NOTEXISTS = errors.New("user not exists...")
	ERROR_USER_EXISTS    = errors.New("user already exists...")
	ERROR_USER_PWD       = errors.New("password or userId not correct..")
)
