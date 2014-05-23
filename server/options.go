package server

// Options contains the bits used to create a server
type Options struct {
	DatabaseURL     string
	FileStorePrefix string
}

// NewOptions makes new *Options wheeee
func NewOptions() *Options {
	return &Options{}
}
