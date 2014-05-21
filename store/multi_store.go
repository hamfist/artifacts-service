package store

import (
	"io"
	"sync"
)

// MultiStore allows for multiple stores in one!
type MultiStore struct {
	storers     map[string]Storer
	storerMutex *sync.Mutex
}

// NewMultiStore creates a new *MultiStore with initialized bits
func NewMultiStore() *MultiStore {
	return &MultiStore{
		storers:     make(map[string]Storer),
		storerMutex: &sync.Mutex{},
	}
}

// Store tees the in io.Reader across all Storers
func (ms *MultiStore) Store(in io.Reader) error {
	ms.storerMutex.Lock()
	defer ms.storerMutex.Unlock()

	var (
		nextReader, r io.Reader
		w             io.Writer
	)

	nextReader = in

	for _, storer := range ms.storers {
		r, w = io.Pipe()
		nextReader = io.TeeReader(nextReader, w)
		err := storer.Store(r)
		if err != nil {
			return err
		}
	}

	return nil
}

// AddStorer allows one to add a storer by name
func (ms *MultiStore) AddStorer(name string, storer Storer) {
	ms.storerMutex.Lock()
	defer ms.storerMutex.Unlock()

	ms.storers[name] = storer
}

// RemoveStorer allows one to remove a storer by name
func (ms *MultiStore) RemoveStorer(name string) {
	ms.storerMutex.Lock()
	defer ms.storerMutex.Unlock()

	if _, ok := ms.storers[name]; ok {
		delete(ms.storers, name)
	}
}
