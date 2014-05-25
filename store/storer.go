package store

import (
	"github.com/meatballhat/artifacts-service/artifact"
)

// Storer defines how stuff gets stored
type Storer interface {
	Store(*artifact.Artifact) error
	Fetch(slug, path string) (*artifact.Artifact, error)
}
