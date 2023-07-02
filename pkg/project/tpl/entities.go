package tpl

func EntityHandlersTempalte() []byte {
	return []byte(`package {{.DirNames.EntityDir}}

import "github.com/adharshmk96/stk/gsk"

type PingHandler interface {
	PingHandler(ctx gsk.Context)
}
`)
}

func EntityServicesTemplate() []byte {
	return []byte(`package {{.DirNames.EntityDir}}

type PingService interface {
	PingService() string
}
`)
}

func EntityStorageTemplate() []byte {
	return []byte(`package {{.DirNames.EntityDir}}

type PingStorage interface {
	Ping() error
}
`)
}
