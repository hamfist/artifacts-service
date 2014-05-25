package artifact

import (
	"fmt"
	"io"
	"path/filepath"
	"time"
)

var (
	errNoReader = fmt.Errorf("no reader available")
)

// Artifact contains the bits!
type Artifact struct {
	Branch       string
	BuildID      string
	BuildNumber  string
	Source       string
	Destination  string
	Instream     io.Reader
	Outstream    io.ReadSeeker
	JobID        string
	JobNumber    string
	RepoSlug     string
	Size         uint64
	DateModified time.Time
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

// Reader provides an io.Reader for the raw bytes
func (a *Artifact) Reader() (io.ReadSeeker, error) {
	if a.Outstream == nil {
		return nil, errNoReader
	}
	return a.Outstream, nil
}

// Fullpath returns the full destination path
func (a *Artifact) Fullpath() string {
	return filepath.Join(a.RepoSlug, "jobs", a.JobID, a.Destination)
}
