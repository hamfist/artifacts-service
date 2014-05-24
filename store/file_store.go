package store

import (
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/meatballhat/artifacts-service/artifact"
)

// FileStore stores stuff as files.  Wow!
type FileStore struct {
	Prefix string
}

// NewFileStore returns a *FileStore.  AMAZE.
func NewFileStore(prefix string) *FileStore {
	return &FileStore{
		Prefix: prefix,
	}
}

// Store does the storing
func (fs *FileStore) Store(a *artifact.Artifact) error {
	fullPath := filepath.Join(fs.Prefix, strings.TrimPrefix(a.Fullpath(), "/"))
	fullPathPrefix := path.Dir(fullPath)

	err := os.MkdirAll(fullPathPrefix, 0755)
	if err != nil {
		return err
	}

	fd, err := ioutil.TempFile("", "artifacts-tmp")
	if err != nil {
		return err
	}

	defer fd.Close()

	_, err = io.CopyN(fd, a.Instream, int64(a.Size))
	if err != nil {
		return err
	}

	defer os.Remove(fd.Name())
	return os.Rename(fd.Name(), fullPath)
}
