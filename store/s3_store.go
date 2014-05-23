package store

import (
	// "github.com/mitchellh/goamz/s3"
	"github.com/meatballhat/artifacts-service/artifact"
)

// S3Store stores stuff in S3
type S3Store struct {
	Key    string
	Secret string
	Bucket string
}

// NewS3Store makes a new *S3Store
func NewS3Store(key, secret, bucket string) *S3Store {
	return &S3Store{}
}

// Store does it up real nice!
func (s3s *S3Store) Store(a *artifact.Artifact) error {
	return nil
}
