package store

import "testing"

func TestNewS3StoreDefaults(t *testing.T) {
	log := getPanicLogger()
	db := getTestDB()
	s3s, err := NewS3Store("key", "secret", "bucket", "us-west-2", log, db)
	if err != nil {
		t.Fatalf("%v", err)
	}

	if s3s == nil {
		t.Fatalf("OH NO")
	}
}
