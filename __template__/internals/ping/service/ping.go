package service

import (
	"github.com/adharshmk96/stktemplate/internals/ping/domain"
)

type pingService struct {
	storage domain.PingStorage
}

func NewPingService(storage domain.PingStorage) domain.PingService {
	return &pingService{
		storage: storage,
	}
}

func (s *pingService) PingService() (string, error) {
	err := s.storage.Ping()
	if err != nil {
		return "", err
	}
	return "pong", nil
}
