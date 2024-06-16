package handler

import "errors"

var (
	errInvalidJSONBodyRequest = errors.New("invalid JSON body request")
	errInvalidInput           = errors.New("invalid input")
	errSomethingWentWrong     = errors.New("something went wrong")
)
