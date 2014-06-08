package server

import (
	"os"

	"github.com/kelseyhightower/envconfig"
)

// Options contains the bits used to create a server
type Options struct {
	DatabaseURL     string
	FileStorePrefix string
	StorerType      string

	S3Key    string
	S3Secret string
	S3Bucket string

	Debug bool
}

// NewOptions makes new *Options wheeee
func NewOptions() *Options {
	dbURL := os.Getenv("ARTIFACTS_DATABASEURL")
	if dbURL == "" {
		dbURL = os.Getenv("DATABASE_URL")
	}

	storerType := os.Getenv("ARTIFACTS_STORERTYPE")
	if storerType == "" {
		storerType = "file"
	}

	opts := &Options{
		DatabaseURL: dbURL,
		StorerType:  storerType,

		S3Key:    os.Getenv("ARTIFACTS_KEY"),
		S3Secret: os.Getenv("ARTIFACTS_SECRET"),
		S3Bucket: os.Getenv("ARTIFACTS_BUCKET"),

		Debug: os.Getenv("DEBUG") != "",
	}

	envconfig.Process("artifacts", opts)
	return opts
}
