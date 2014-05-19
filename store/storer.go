package store

import (
	"io"
)

// Storer defines how stuff gets stored
type Storer interface {
	Store(io.Reader) error
}
