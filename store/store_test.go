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

type testSaver struct {
	mds map[string]*metadata.Metadata
}

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

func TestStuff(t *testing.T) {
	if false {
		t.Fail()
	}
}
