package store

import (
	"os"
	"testing"

	"github.com/Sirupsen/logrus"
	"github.com/hamfist/artifacts-service/metadata"
)

var (
	dbURL = os.Getenv("DATABASE_URL")
)

func init() {
	os.Clearenv()
	os.Setenv("DATABASE_URL", dbURL)
}

func getPanicLogger() *logrus.Logger {
	log := logrus.New()
	if os.Getenv("DEBUG") != "" && os.Getenv("ARTIFACTS_DEBUG") != "" {
		log.Level = logrus.DebugLevel
	}
	return log
}

func getTestDB() *metadata.Database {
	db, err := metadata.NewDatabase(dbURL, getPanicLogger())
	if err != nil {
		panic(err)
	}

	return db
}

func TestStuff(t *testing.T) {
	if false {
		t.Fail()
	}
}
