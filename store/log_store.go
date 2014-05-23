package store

import (
	"github.com/Sirupsen/logrus"
	"github.com/meatballhat/artifacts-service/artifact"
)

// LogStore implements Storer and logs stuff
type LogStore struct {
	log *logrus.Logger
}

// NewLogStore makes a *LogStore. wow!
func NewLogStore() *LogStore {
	return &LogStore{
		log: logrus.New(),
	}
}

// Store reports stuff about what is being stored
func (ls *LogStore) Store(a *artifact.Artifact) error {
	return nil
}
