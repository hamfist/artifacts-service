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
	ID          uint64
	JobID       string
	Size        uint64
	Path        string
	ContentType string
}

// MarshalJSON is all about the JSON, everything as a string per jsonapi
func (md *Metadata) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]string{
		"id":           fmt.Sprintf("%d", md.ID),
		"size":         fmt.Sprintf("%d", md.Size),
		"path":         md.Path,
		"content_type": md.ContentType,
	})
}
