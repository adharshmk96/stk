package tpl

func ServiceServiceTemplate() []byte {
	return []byte(`package {{.DirNames.ServiceDir}}

import (
	"{{.PkgName}}/{{.DirTree.EntityPath}}"
)

type pingService struct {
	pingStorage {{.DirNames.EntityDir}}.PingStorage
}

func NewPingService(storage {{.DirNames.EntityDir}}.PingStorage) {{.DirNames.EntityDir}}.PingService {
	return &pingService{
		pingStorage: storage,
	}
}
`)
}

func ServicePingTemplate() []byte {
	return []byte(`package {{.DirNames.ServiceDir}}

func (s *pingService) PingService() string {
	err := s.pingStorage.Ping()
	if err != nil {
		return "error"
	}
	return "pong"
}
`)
}
