package domain

import "github.com/adharshmk96/stk/gsk"

// Handler
type PingHandlers interface {
	PingHandler(gc *gsk.Context)
}
