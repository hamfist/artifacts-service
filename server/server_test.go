package server

import (
	"os"
	"reflect"
	"testing"

	"github.com/Sirupsen/logrus"
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

func TestServerDefaults(t *testing.T) {
	opts := NewOptions()
	log := getPanicLogger()

	srv, err := NewServer(opts, log)
	if err != nil {
		t.Fatalf("%v", err)
	}

	if !reflect.DeepEqual(srv.opts, opts) {
		t.Fatalf("opts %v != %v", srv.opts, opts)
	}

}
