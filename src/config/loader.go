package config

import "time"

// Loader is a object that can import, initialize, and Get config values
//go:generate go run -mod=mod github.com/golang/mock/mockgen -package loadermock -destination=./loadermock/mock_loader.go -source=../config/loader.go -build_flags=-mod=mod
type Loader interface {
	Import(data []byte) error
	Initialize() error
	Get(key string) ([]byte, error)
	Put(key string, value []byte) error

	// Must functions will panic if they can't do what is requested.
	// They are maingly meant for use with configs that are required for an app to start up
	MustGetString(key string) string
	MustGetBool(key string) bool
	MustGetInt(key string) int
	MustGetDuration(key string) time.Duration

	//TODO add array support?
}
