package gsk

import "errors"

var (
	ErrInvalidJSON    = errors.New("invalid_json")
	ErrInternalServer = errors.New("internal_server_error")
	ErrBodyTooLarge   = errors.New("request_body_too_large")
)
