package server

import (
	"fmt"
	"net/http"
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/hamfist/artifacts-service/auth"
	"github.com/hamfist/artifacts-service/metadata"
	"github.com/hamfist/artifacts-service/store"
)

var (
	errMissingS3Key    = fmt.Errorf("missing s3 key")
	errMissingS3Secret = fmt.Errorf("missing s3 secret")
)

// Server holds onto a router and a store
type Server struct {
	Router *mux.Router

	n     *negroni.Negroni
	opts  *Options
	log   *logrus.Logger
	store store.Storer
	auth  auth.Auther
	md    *metadata.Database
}

// Main is the top of the pile.  Start here.
func Main(log *logrus.Logger) {
	opts := NewOptions()
	if opts.Debug {
		log.Level = logrus.DebugLevel
	}

	if opts.FileStorePrefix == "" {
		opts.FileStorePrefix = "tmp"
	}

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

	server.Run(addr)
}

// NewServer creates a new *Server with a router and its routes registered
func NewServer(opts *Options, log *logrus.Logger) (*Server, error) {
	var err error

	log.Debug("creating new server")
	srv := &Server{
		opts: opts,
		log:  log,
	}

	srv.setupRouter()
	srv.setupNegroni()

	err = srv.getDB()
	if err != nil {
		return nil, err
	}

	err = srv.setupStorer()
	if err != nil {
		return nil, err
	}

	err = srv.setupAuther()
	if err != nil {
		return nil, err
	}

	return srv, nil
}

// Run starts up the server on the given addr
func (srv *Server) Run(addr string) {
	srv.n.Run(addr)
}

func (srv *Server) setupRouter() {
	srv.log.Debug("setting up router")
	router := mux.NewRouter()

	router.HandleFunc(`/job/{job_id}/{filepath:.+}`,
		muxVarsWrapper(srv.saveHandler)).Methods("PUT").Name("save_job_artifact")
	router.HandleFunc(`/job/{job_id}`,
		muxVarsWrapper(srv.listHandler)).Methods("GET").Name("list_job_artifacts")
	router.HandleFunc(`/job/{job_id}/{filepath:.+}`,
		muxVarsWrapper(srv.getPathHandler)).Methods("GET").Name("get_job_artifact")
	router.HandleFunc(`/{owner}/{repo}/jobs/{job_id}/{filepath:.+}`,
		muxVarsWrapper(srv.saveHandler)).Methods("PUT").Name("legacy_save_job_artifact")
	router.HandleFunc(`/{owner}/{repo}/jobs/{job_id}`,
		muxVarsWrapper(srv.listHandler)).Methods("GET").Name("legacy_list_job_artifacts")
	router.HandleFunc(`/{owner}/{repo}/jobs/{job_id}/{filepath:.+}`,
		muxVarsWrapper(srv.getPathHandler)).Methods("GET").Name("legacy_get_job_artifact")

	srv.Router = router
}

func muxVarsWrapper(f func(http.ResponseWriter, *http.Request, map[string]string) int) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		f(w, r, mux.Vars(r))
	}
}

func (srv *Server) setupNegroni() {
	srv.log.Debug("setting up negroni")
	srv.n = negroni.New()
	srv.n.Use(negroni.NewRecovery())
	srv.n.Use(NewLoggerMiddleware())
	srv.n.Use(negroni.NewStatic(http.Dir("public")))
	srv.n.UseHandler(srv.Router)
}

func (srv *Server) setupStorer() error {
	srv.log.WithField("storer_type", srv.opts.StorerType).Debug("setting up storer")
	switch srv.opts.StorerType {
	case "s3":
		if srv.opts.S3Key == "" {
			return errMissingS3Key
		}
		if srv.opts.S3Secret == "" {
			return errMissingS3Secret
		}
		store, err := store.NewS3Store(srv.opts.S3Key,
			srv.opts.S3Secret, srv.opts.S3Bucket, srv.opts.S3Region, srv.log, srv.md)
		if err != nil {
			srv.log.WithField("err", err).Debug("error setting up s3 store")
			return err
		}

		srv.log.WithField("store", store).Debug("assigning s3 store")
		srv.store = store
		return nil
	case "file":
		srv.store = store.NewFileStore(srv.opts.FileStorePrefix,
			srv.log, srv.md)
		srv.log.WithField("store", srv.store).Debug("assigning file store")
		return nil
	default:
		srv.log.WithFields(logrus.Fields{
			"storer_type": srv.opts.StorerType,
		}).Error("unknown storer type")
		return fmt.Errorf("unknown storer type %q", srv.opts.StorerType)
	}

	panic("fell through to a bad place ¯\\_(ツ)_/¯")
}

func (srv *Server) setupAuther() error {
	srv.log.WithField("auther_type", srv.opts.AutherType).Debug("setting up auther")
	switch srv.opts.AutherType {
	case "token":
		srv.auth = auth.NewTokenAuther(srv.opts.AuthToken)
		return nil
	case "travis":
		srv.auth = auth.NewTravisAuther(srv.opts.TravisAPIServer,
			srv.opts.TravisPrivateKey(), srv.opts.TravisRequireRSA)
		return nil
	case "null":
		srv.auth = auth.NewNullAuther()
		return nil
	default:
		srv.log.WithFields(logrus.Fields{
			"auther_type": srv.opts.AutherType,
		}).Error("unknown auther type")
		return fmt.Errorf("unknown auther type %q", srv.opts.AutherType)
	}

	panic("fell through to a bad place ¯\\_(ツ)_/¯")
}

func (srv *Server) getDB() error {
	srv.log.Debug("getting database handle")
	db, err := metadata.NewDatabase(srv.opts.DatabaseURL, srv.log)
	if err != nil {
		return err
	}

	err = db.Init()
	if err != nil {
		return err
	}

	srv.md = db
	return nil
}

func (srv *Server) canWrite(r *http.Request, vars map[string]string) bool {
	return srv.auth.Check(r, vars).CanWrite
}

func (srv *Server) canRead(r *http.Request, vars map[string]string) bool {
	return srv.auth.Check(r, vars).CanRead
}
