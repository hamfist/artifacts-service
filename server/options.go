package server

import (
	"os"
)

// Options contains the bits used to create a server
type Options struct {
	DatabaseURL     string
	FileStorePrefix string
	StorerType      string

	S3Key    string
	S3Secret string

	Debug bool
}

// NewOptions makes new *Options wheeee
func NewOptions() *Options {
	dbURL := os.Getenv("ARTIFACTS_DATABASE_URL")
	if dbURL == "" {
		dbURL = os.Getenv("DATABASE_URL")
	}

	storerType := os.Getenv("ARTIFACTS_STORER_TYPE")
	if storerType == "" {
		storerType = "file"
	}

	return &Options{
		DatabaseURL:     dbURL,
		FileStorePrefix: os.Getenv("ARTIFACTS_FILE_STORE_PREFIX"),
		StorerType:      storerType,

		S3Key:    os.Getenv("ARTIFACTS_KEY"),
		S3Secret: os.Getenv("ARTIFACTS_SECRET"),

		Debug: os.Getenv("ARTIFACTS_DEBUG") != "",
	}
}
