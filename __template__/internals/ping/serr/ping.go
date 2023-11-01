package serr

import "errors"

var (
	ErrPingFailed = errors.New("ping failed")
)
