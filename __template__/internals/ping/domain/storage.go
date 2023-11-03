package domain

// Storage
type PingStorage interface {
	Ping() error
}
