package domain

// Service
type PingService interface {
	PingService() (string, error)
}
