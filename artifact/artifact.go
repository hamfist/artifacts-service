package artifact

import (
	"io"
	"path/filepath"
)

// Artifact contains the bits!
type Artifact struct {
	Branch      string
	BuildID     string
	BuildNumber string
	Source      string
	Destination string
	Instream    io.Reader
	JobID       string
	JobNumber   string
	RepoSlug    string
	Size        uint64
}

// New creates a new *Artifact
func New(repoSlug, src, dest, jobID string, in io.Reader, size uint64) *Artifact {
	return &Artifact{
		Source:      src,
		Destination: dest,
		Instream:    in,
		JobID:       jobID,
		RepoSlug:    repoSlug,
		Size:        size,
	}
}

// Fullpath returns the full destination path
func (a *Artifact) Fullpath() string {
	return filepath.Join(a.RepoSlug,
		a.BuildNumber, a.JobNumber, a.Destination)
}
