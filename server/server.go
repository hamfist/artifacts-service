package server

import (
	"fmt"
	"net/http"
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"github.com/meatballhat/artifacts-service/store"
)

// Server holds onto a router and a store
type Server struct {
	Router *mux.Router
	opts   *Options
	log    *logrus.Logger
	store  store.Storer
}

// Main is the top of the pile.  Start here.
func Main() {
	opts := NewOptions()
	if opts.FileStorePrefix == "" {
		opts.FileStorePrefix = "tmp"
	}

	log := logrus.New()

	server, err := NewServer(opts, log)
	if err != nil {
		log.Fatal(err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "9839"
	}

	addr := fmt.Sprintf(":%s", port)
	log.WithFields(logrus.Fields{
		"addr": addr,
	}).Info("artifacts-service listening")

	log.Fatal(http.ListenAndServe(addr, server.Router))
}

// NewServer creates a new *Server with a router and its routes registered
func NewServer(opts *Options, log *logrus.Logger) (*Server, error) {
	if opts.Debug {
		log.Level = logrus.Debug
	}

	server := &Server{
		opts: opts,
		log:  log,
	}

	server.setupRouter()
	err := server.setupStorer()
	if err != nil {
		return nil, err
	}

	return server, nil
}

func (srv *Server) setupRouter() {
	router := mux.NewRouter()

	router.HandleFunc(`/`,
		func(w http.ResponseWriter, r *http.Request) {
			srv.logRequest(r)
			srv.rootHandler(w, r, mux.Vars(r))
		}).Methods("GET")

	router.HandleFunc(`/{slug:[^/]+/[^/]+}/jobs/{job_id}/{filepath:[a-zA-Z0-9_\-\/]+}`,
		func(w http.ResponseWriter, r *http.Request) {
			srv.logRequest(r)
			srv.saveHandler(w, r, mux.Vars(r))
		}).Methods("PUT")

	router.HandleFunc(`/{slug:[^/]+/[^/]+}/jobs/{job_id}`,
		func(w http.ResponseWriter, r *http.Request) {
			srv.logRequest(r)
			srv.listHandler(w, r, mux.Vars(r))
		}).Methods("GET")

	router.HandleFunc(`/{slug:[^/]+/[^/]+}/jobs/{job_id}/{filepath:[a-zA-Z0-9_\-\/]+}`,
		func(w http.ResponseWriter, r *http.Request) {
			srv.logRequest(r)
			srv.getPathHandler(w, r, mux.Vars(r))
		}).Methods("GET")

	srv.Router = router
}

func (srv *Server) logRequest(r *http.Request) {
	srv.log.WithFields(logrus.Fields{
		"method": r.Method,
		"path":   r.URL.Path,
	}).Info("handling")
}

func (srv *Server) setupStorer() error {
	switch srv.opts.StorerType {
	case "s3":
		srv.store = store.NewS3Store(srv.opts.S3Key, srv.opts.S3Secret,
			srv.opts.S3Bucket, srv.log)
		return nil
	case "file":
		srv.store = store.NewFileStore(srv.opts.FileStorePrefix, srv.log)
		return nil
	default:
		srv.log.WithFields(logrus.Fields{
			"storer_type": srv.opts.StorerType,
		}).Error("unknown storer type")
		return fmt.Errorf("unknown storer type %q", srv.opts.StorerType)
	}

	return nil
}
