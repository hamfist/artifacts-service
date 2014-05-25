package store

import (
	"fmt"

	"github.com/meatballhat/artifacts-service/artifact"
)

var (
	errNotImplemented = fmt.Errorf("brb under construction")
)

// Storer defines how stuff gets stored
type Storer interface {
	Store(*artifact.Artifact) error
	Fetch(string, string, string) (*artifact.Artifact, error)
}
