package artifact

import (
	"io"
	"io/ioutil"
)

// ReadSeekerCloser is both a ReadCloser and a ReadSeeker
type ReadSeekerCloser interface {
	io.ReadCloser
	io.Seeker
}

type readSeekerCloserWrapper struct {
	io.ReadCloser
}

// NewReadSeekerCloser wraps an io.ReadCloser and fakes out Seek
func NewReadSeekerCloser(in io.ReadCloser) *readSeekerCloserWrapper {
	return &readSeekerCloserWrapper{in}
}

func (rsc *readSeekerCloserWrapper) Seek(offset int64, whence int) (int64, error) {
	if offset == int64(0) && whence == 0 {
		return int64(0), nil
	}

	n, err := io.CopyN(ioutil.Discard, rsc, offset)
	if err != nil || n != offset {
		return int64(0), err
	}
	return n, nil
}
