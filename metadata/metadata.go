package metadata

import (
	"fmt"
)

var (
	errNoMetadata = fmt.Errorf("no metadata available")
)

// Metadata is the stuff we care about in the metadata database
type Metadata struct {
	JobID       string
	Size        uint64
	Path        string
	ContentType string
}
