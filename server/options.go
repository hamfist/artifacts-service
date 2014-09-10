package server

import (
	"crypto/rsa"
	"crypto/x509"
	"os"
	"strings"

	"github.com/kelseyhightower/envconfig"
)

// Options contains the bits used to create a server
type Options struct {
	DatabaseURL     string
	FileStorePrefix string
	StorerType      string
	AutherType      string
	AuthToken       string

	TravisAPIServer        string
	TravisPrivateKeyString string
	TravisRequireRSA       bool

	S3Key    string
	S3Secret string
	S3Bucket string
	S3Region string

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
		DatabaseURL: os.Getenv("DATABASE_URL"),
		StorerType:  storerType,
		AutherType:  autherType,
		AuthToken:   os.Getenv("ARTIFACTS_TOKEN"),

		TravisAPIServer:        travisAPIServer,
		TravisPrivateKeyString: os.Getenv("TRAVIS_PRIVATE_KEY_STRING"),
		TravisRequireRSA:       os.Getenv("TRAVIS_REQUIRE_RSA") != "",

		S3Key:    os.Getenv("ARTIFACTS_KEY"),
		S3Secret: os.Getenv("ARTIFACTS_SECRET"),
		S3Bucket: os.Getenv("ARTIFACTS_BUCKET"),
		S3Region: s3Region,

		Debug: os.Getenv("DEBUG") != "",
	}

	envconfig.Process("artifacts", opts)
	return opts
}

func (o *Options) String() string {
	return "&server.Options{[secrets]}"
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
