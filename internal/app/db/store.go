package db

// Store ...
type Store interface {
	User() UserRepository
}
