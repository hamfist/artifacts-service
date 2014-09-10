package server

import (
	"os"
	"testing"
)

func TestNewOptionsDefaults(t *testing.T) {
	opts := NewOptions()

	for actual, expected := range map[string]string{
		opts.AuthToken:       "",
		opts.AutherType:      "token",
		opts.DatabaseURL:     os.Getenv("DATABASE_URL"),
		opts.FileStorePrefix: "",
		opts.StorerType:      "file",

		opts.TravisAPIServer:        "https://api.travis-ci.org",
		opts.TravisPrivateKeyString: "",

		opts.S3Bucket: "",
		opts.S3Key:    "",
		opts.S3Region: "us-east-1",
		opts.S3Secret: "",
	} {
		if actual != expected {
			t.Fatalf("%v != %v", actual, expected)
		}
	}

	if opts.TravisRequireRSA {
		t.Fatal()
	}

	if opts.Debug {
		t.Fatal()
	}

}
