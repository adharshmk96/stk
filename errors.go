package stk

import "errors"

var ErrInvalidJSON = errors.New("invalid_json")
var ErrInternalServer = errors.New("internal_server_error")
var ErrBodyTooLarge = errors.New("request_body_too_large")
