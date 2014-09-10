package server

import (
	"crypto/rsa"
	"crypto/x509"
	"fmt"
	"os"
	"strings"

	"github.com/kelseyhightower/envconfig"
)

// Options contains the bits used to create a server
type Options struct {
	AuthToken       string
	AutherType      string
	DatabaseURL     string
	FileStorePrefix string
	StorerType      string

	TravisAPIServer        string
	TravisPrivateKeyString string
	TravisRequireRSA       bool

	S3Bucket string
	S3Key    string
	S3Region string
	S3Secret string

	Debug bool
}

// NewOptions makes new *Options wheeee
// defaulting to an enterprise-friendly configuration.
func NewOptions() *Options {
	storerType := os.Getenv("ARTIFACTS_STORER_TYPE")
	if storerType == "" {
		storerType = "file"
	}

	autherType := os.Getenv("ARTIFACTS_AUTHER_TYPE")
	if autherType == "" {
		autherType = "token"
	}

	travisAPIServer := os.Getenv("ARTIFACTS_TRAVIS_API_SERVER")
	if travisAPIServer == "" {
		travisAPIServer = "https://api.travis-ci.org"
	}

	s3Region := os.Getenv("ARTIFACTS_S3_REGION")
	if s3Region == "" {
		s3Region = "us-east-1"
	}

	opts := &Options{
		AuthToken:       os.Getenv("ARTIFACTS_TOKEN"),
		AutherType:      autherType,
		DatabaseURL:     os.Getenv("DATABASE_URL"),
		FileStorePrefix: os.Getenv("ARTIFACTS_FILE_STORE_PREFIX"),
		StorerType:      storerType,

		TravisAPIServer:        travisAPIServer,
		TravisPrivateKeyString: os.Getenv("TRAVIS_PRIVATE_KEY_STRING"),
		TravisRequireRSA:       os.Getenv("TRAVIS_REQUIRE_RSA") != "",

		S3Bucket: os.Getenv("ARTIFACTS_BUCKET"),
		S3Key:    os.Getenv("ARTIFACTS_KEY"),
		S3Region: s3Region,
		S3Secret: os.Getenv("ARTIFACTS_SECRET"),

		Debug: os.Getenv("DEBUG") != "",
	}

	envconfig.Process("artifacts", opts)
	return opts
}

func (o *Options) String() string {
	return fmt.Sprintf("&server.Options{[secrets]} %p", o)
}

// TravisPrivateKey parses and returns an *rsa.PrivateKey from the
// TravisPrivateKeyString if present and valid
func (o *Options) TravisPrivateKey() *rsa.PrivateKey {
	trimmed := strings.TrimSpace(o.TravisPrivateKeyString)
	if trimmed == "" {
		return nil
	}

	for _, s := range []string{"\n", " ", "\t"} {
		trimmed = strings.Replace(trimmed, s, "", -1)
	}

	privateKey, err := x509.ParsePKCS1PrivateKey([]byte(trimmed))
	if err != nil {
		return nil
	}

	return privateKey
}
