package store

import (
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/dustin/go-humanize"
	"github.com/hamfist/artifacts-service/artifact"
	"github.com/hamfist/artifacts-service/metadata"
)

// FileStore stores stuff as files.  Wow!
type FileStore struct {
	Prefix string
	log    *logrus.Logger
	db     *metadata.Database
}

// NewFileStore returns a *FileStore.  AMAZE.
func NewFileStore(prefix string, log *logrus.Logger, db *metadata.Database) *FileStore {
	return &FileStore{
		Prefix: prefix,
		log:    log,
		db:     db,
	}
}

func (fs *FileStore) artifactFullPath(a *artifact.Artifact) string {
	return filepath.Join(fs.Prefix, strings.TrimPrefix(a.FullDestination(), "/"))
}

// Store does the storing
func (fs *FileStore) Store(a *artifact.Artifact) error {
	fullPath := fs.artifactFullPath(a)
	fullPathPrefix := path.Dir(fullPath)

	err := os.MkdirAll(fullPathPrefix, 0755)
	if err != nil {
		fs.log.WithFields(logrus.Fields{
			"err":    err,
			"prefix": fullPathPrefix,
		}).Error("failed to make dest file path prefix")
		return err
	}

	fd, err := ioutil.TempFile("", "artifacts-tmp")
	if err != nil {
		fs.log.WithFields(logrus.Fields{
			"err": err,
		}).Error("failed to get tempfile")
		return err
	}

	defer fd.Close()

	_, err = io.CopyN(fd, a.Instream, int64(a.Size))
	if err != nil {
		fs.log.WithFields(logrus.Fields{
			"err":       err,
			"temp_file": fd.Name(),
		}).Error("failed to copy to temporary file")
		return err
	}

	defer os.Remove(fd.Name())
	err = os.Rename(fd.Name(), fullPath)
	if err != nil {
		fs.log.WithFields(logrus.Fields{
			"err":       err,
			"temp_file": fd.Name(),
			"dest":      fullPath,
		}).Error("failed to move temporary file to dest")
	}

	fs.log.WithFields(logrus.Fields{
		"source": a.Source,
		"prefix": fs.Prefix,
		"dest":   a.FullDestination(),
		"size":   humanize.Bytes(a.Size),
	}).Info("stored artifact to file")

	return nil
}

// Fetch returns an artifact given a path and job id
func (fs *FileStore) Fetch(path, jobID string) (*artifact.Artifact, error) {
	a := artifact.New(path, path, jobID, nil, uint64(0))
	fullPath := fs.artifactFullPath(a)

	fi, err := os.Stat(fullPath)
	if err != nil {
		return nil, err
	}

	fd, err := os.Open(fullPath)
	if err != nil {
		return nil, err
	}

	a.Size = uint64(fi.Size())
	a.DateModified = fi.ModTime()
	a.OutReadSeeker = fd

	return a, nil
}
