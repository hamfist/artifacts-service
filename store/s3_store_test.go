package store

import (
	"testing"

	"github.com/hamfist/artifacts-service/metadata"
)

func TestNewS3StoreDefaults(t *testing.T) {
	log := getPanicLogger()
	md := metadata.NewNullLookupSaver()
	s3s, err := NewS3Store("key", "secret", "bucket", "us-west-2", log, md)
	if err != nil {
		t.Fatalf("%v", err)
	}

	if s3s == nil {
		t.Fatalf("OH NO")
	}
}
