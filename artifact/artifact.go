package artifact

import (
	"fmt"
	"io"
	"path/filepath"
	"time"

	"github.com/mitchellh/goamz/s3"
)

var (
	errNoReader = fmt.Errorf("no reader available")
)

// Artifact contains the bits!
type Artifact struct {
	Branch      string
	BuildID     string
	BuildNumber string
	JobID       string
	JobNumber   string
	RepoSlug    string

	Source       string
	Destination  string
	Size         uint64
	DateModified time.Time
	ContentType  string
	Perm         s3.ACL

	Instream      io.Reader
	OutReadSeeker io.ReadSeeker
	OutReadCloser io.ReadCloser
}

// New creates a new *Artifact
func New(repoSlug, src, dest, jobID string, in io.Reader, size uint64) *Artifact {
	return &Artifact{
		Source:       src,
		Destination:  dest,
		Instream:     in,
		JobID:        jobID,
		RepoSlug:     repoSlug,
		Size:         size,
		DateModified: time.Now().UTC(),
	}
}

// ReadCloser provides an io.ReadCloser for the raw bytes
func (a *Artifact) ReadCloser() (io.ReadCloser, error) {
	if a.OutReadCloser == nil {
		return nil, errNoReader
	}
	return a.OutReadCloser, nil
}

// ReadSeeker provides an io.ReadSeeker for the raw bytes
func (a *Artifact) ReadSeeker() (io.ReadSeeker, error) {
	if a.OutReadSeeker == nil {
		return nil, errNoReader
	}
	return a.OutReadSeeker, nil
}

// FullDestination returns the full destination path
func (a *Artifact) FullDestination() string {
	return filepath.Join(a.RepoSlug, "jobs", a.JobID, a.Destination)
}
