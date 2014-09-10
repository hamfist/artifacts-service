package metadata

// LookupSaver is the interface needed for reads and writes of metadata
type LookupSaver interface {
	Save(*Metadata) error
	Lookup(string, string) (*Metadata, error)
	LookupAll(string) ([]*Metadata, error)
}
