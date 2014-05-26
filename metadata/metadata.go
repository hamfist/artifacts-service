package metadata

import (
	"fmt"
)

var (
	errNoMetadata = fmt.Errorf("no metadata available")
)

// Metadata is the stuff we care about in the metadata database
type Metadata struct {
	Owner       string
	Repo        string
	BuildID     string
	BuildNumber string
	JobID       string
	JobNumber   string
	Path        string
}
