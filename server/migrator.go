package server

import (
	"github.com/Sirupsen/logrus"
	"github.com/hamfist/artifacts-service/metadata"
)

func MigratorMain(log *logrus.Logger) {
	opts := NewOptions()
	if opts.Debug {
		log.Level = logrus.DebugLevel
	}

	log.Debug("spinning up database")

	db, err := metadata.NewDatabase(opts.DatabaseURL, log)
	if err != nil {
		log.Fatal(err)
	}

	log.Debug("migrating")

	err = db.Migrate(log)
	if err != nil {
		log.Fatal(err)
	}

	log.Info("database migration complete")
}
