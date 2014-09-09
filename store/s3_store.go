package store

import (
	"fmt"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/hamfist/artifacts-service/artifact"
	"github.com/hamfist/artifacts-service/metadata"
	"github.com/mitchellh/goamz/aws"
	"github.com/mitchellh/goamz/s3"
)

var (
	errNoBucket = fmt.Errorf("no bucket found")
)

// S3Store is a Storer for S3!!!
type S3Store struct {
	key    string
	secret string
	bucket string
	log    *logrus.Logger
	db     *metadata.Database
	b      *s3.Bucket
}

// NewS3Store initializes an *S3Store.  Wow!
func NewS3Store(key, secret, bucket, regionName string,
	log *logrus.Logger, db *metadata.Database) (*S3Store, error) {

	log.Debug("getting aws auth")
	auth, err := aws.GetAuth(key, secret)
	if err != nil {
		log.WithField("err", err).Error("failed to get auth")
		return nil, err
	}

	region, ok := aws.Regions[regionName]
	if !ok {
		log.WithFields(logrus.Fields{
			"region": regionName,
		}).Warn(fmt.Sprintf("nonexistent region, falling back to %s", aws.USEast.Name))
		region = aws.USEast
	}

	log.Debug("getting new s3 connection")
	s3Conn := s3.New(auth, region)
	b := s3Conn.Bucket(bucket)

	if b == nil || b.Name == "" {
		return nil, errNoBucket
	}

	log.WithFields(logrus.Fields{
		"bucket": b.Name,
	}).Debug("got back this bucket")

	return &S3Store{
		key:    key,
		secret: secret,
		bucket: bucket,

		log: log,
		db:  db,
		b:   b,
	}, nil
}

// Store stores the stuff in the S3
func (s3s *S3Store) Store(a *artifact.Artifact) error {
	destination := a.FullDestination()
	ctype := a.ContentType
	size := a.Size

	s3s.log.WithFields(logrus.Fields{
		"source":       a.Source,
		"dest":         destination,
		"bucket":       s3s.b.Name,
		"content_type": ctype,
	}).Debug("more artifact details")

	err := s3s.b.PutReaderHeader(destination, a.Instream, int64(size),
		map[string][]string{
			"Content-Type": []string{ctype},
		}, s3.Private)
	if err != nil {
		return err
	}

	md := artifactToMetadata(a)
	_, err = s3s.db.Save(md)
	if err != nil {
		return err
	}
	return nil
}

// Fetch returns an artifact given a path and job id
func (s3s *S3Store) Fetch(path, jobID string) (*artifact.Artifact, error) {
	a := artifact.New("", path, jobID, nil, uint64(0))

	s3Key, err := s3s.b.GetKey(a.FullDestination())
	if err != nil {
		return nil, err
	}

	s3Stream, err := s3s.b.GetReader(a.FullDestination())
	if err != nil {
		return nil, err
	}

	a.OutReadCloser = s3Stream
	a.Size = uint64(s3Key.Size)
	dateMod, err := time.Parse(time.RFC1123, s3Key.LastModified)
	if err != nil {
		return nil, err
	}
	a.DateModified = dateMod

	s3s.log.WithFields(logrus.Fields{
		"path":          a.FullDestination(),
		"size":          a.Size,
		"date_modified": a.DateModified,
	}).Debug("returning artifact from s3")

	return a, nil
}

func (s3s *S3Store) String() string {
	return "&store.S3Store{[secrets]}"
}
