package metadata

import (
	"encoding/json"
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

// MarshalJSON is all about the JSON
func (md *Metadata) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]string{
		"size":         fmt.Sprintf("%d", md.Size),
		"path":         md.Path,
		"content_type": md.ContentType,
	})
}
