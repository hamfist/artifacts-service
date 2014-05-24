package server

import (
	"os"
)

// Options contains the bits used to create a server
type Options struct {
	DatabaseURL     string
	FileStorePrefix string
}

// NewOptions makes new *Options wheeee
func NewOptions() *Options {
	dbURL := os.Getenv("ARTIFACTS_DATABASE_URL")
	if dbURL == "" {
		dbURL = os.Getenv("DATABASE_URL")
	}

	return &Options{
		DatabaseURL:     dbURL,
		FileStorePrefix: os.Getenv("ARTIFACTS_FILE_STORE_PREFIX"),
	}
}
