package store

import (
	"fmt"

	"github.com/hamfist/artifacts-service/artifact"
	"github.com/hamfist/artifacts-service/metadata"
)

var (
	errNotImplemented = fmt.Errorf("brb under construction")
)

// Storer defines how stuff gets stored
type Storer interface {
	Store(*artifact.Artifact) error
	Fetch(string, string, string) (*artifact.Artifact, error)
}

func artifactToMetadata(a *artifact.Artifact) *metadata.Metadata {
	return &metadata.Metadata{
		JobID:       a.JobID,
		Size:        a.Size,
		Path:        a.FullDestination(),
		ContentType: a.ContentType,
	}
}
