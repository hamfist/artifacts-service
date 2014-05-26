package store

import (
	"github.com/Sirupsen/logrus"
	"github.com/meatballhat/artifacts-service/artifact"
	"github.com/meatballhat/artifacts-service/metadata"
)

// S3Store is a Storer for S3!!!
type S3Store struct {
	key    string
	secret string
	bucket string
	log    *logrus.Logger
	db     *metadata.Database
}

// NewS3Store initializes an *S3Store.  Wow!
func NewS3Store(key, secret, bucket string,
	log *logrus.Logger, db *metadata.Database) *S3Store {

	return &S3Store{
		key:    key,
		secret: secret,
		bucket: bucket,
		log:    log,
		db:     db,
	}
}

// Store stores the stuff in the S3
func (s3s *S3Store) Store(a *artifact.Artifact) error {
	md := artifactToMetadata(a)
	err := s3s.db.Save(md)
	if err != nil {
		return err
	}
	return nil
}

// Fetch returns an artifact given a repo slug and path
func (s3s *S3Store) Fetch(slug, path, jobID string) (*artifact.Artifact, error) {
	// TODO: fetch crap from S3, ROFL!
	return artifact.New("", "", "", "", nil, uint64(0)), nil
}
