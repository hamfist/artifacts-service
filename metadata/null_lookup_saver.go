package metadata

import (
	"fmt"
	"strings"
	"sync"
)

type NullLookupSaver struct {
	mdMap map[string]*Metadata

	l *sync.Mutex
}

func NewNullLookupSaver() *NullLookupSaver {
	return &NullLookupSaver{
		mdMap: map[string]*Metadata{},

		l: &sync.Mutex{},
	}
}

func (nls *NullLookupSaver) Save(m *Metadata) error {
	nls.l.Lock()
	defer nls.l.Unlock()
	nls.mdMap[fmt.Sprintf("%s-%s", m.JobID, m.Path)] = m
	return nil
}

func (nls *NullLookupSaver) Lookup(jobID, path string) (*Metadata, error) {
	nls.l.Lock()
	defer nls.l.Unlock()
	m, ok := nls.mdMap[fmt.Sprintf("%s-%s", jobID, path)]
	if ok {
		return m, nil
	}
	return nil, errNoMetadata
}

func (nls *NullLookupSaver) LookupAll(jobID string) ([]*Metadata, error) {
	nls.l.Lock()
	defer nls.l.Unlock()
	mds := []*Metadata{}

	for key, m := range nls.mdMap {
		keyParts := strings.SplitN(key, "-", 1)
		if len(keyParts) < 2 {
			return nil, errNoMetadata
		}

		if keyParts[0] == jobID {
			mds = append(mds, m)
		}
	}

	return mds, nil
}
