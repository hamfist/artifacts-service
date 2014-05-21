package store

import (
	"io"

	// "github.com/mitchellh/goamz/s3"
)

// S3Store stores stuff in S3
type S3Store struct{}

// NewS3Store makes a new *S3Store
func NewS3Store() *S3Store {
	return &S3Store{}
}

// Store does it up real nice!
func (s3s *S3Store) Store(in io.Reader) error {
	return nil
}
