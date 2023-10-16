package service

import "github.com/adharshmk96/stk-template/singlemod/internals/core/entity"

type pingService struct {
	storage entity.PingStorage
}

func NewPingService(storage entity.PingStorage) entity.PingService {
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
