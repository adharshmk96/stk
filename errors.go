package stk

import "errors"

var ErrInvalidJSON = errors.New("invalid json")
var ErrInternalServer = errors.New("internal server error")
var ErrBodyTooLarge = errors.New("request body too large")
