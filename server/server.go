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

// Server holds onto a router and a store
type Server struct {
	Router    *mux.Router
	AuthToken string

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

	if opts.Debug {
		log.Level = logrus.DebugLevel
	}

	server := &Server{
		opts: opts,
		log:  log,
	}

	server.setupRouter()
	server.setupNegroni()

	err = server.getDB()
	if err != nil {
		return nil, err
	}

	err = server.setupStorer()
	if err != nil {
		return nil, err
	}

	err = server.setupAuther()
	if err != nil {
		return nil, err
	}

	return server, nil
}

// Run starts up the server on the given addr
func (srv *Server) Run(addr string) {
	srv.n.Run(addr)
}

func (srv *Server) setupRouter() {
	router := mux.NewRouter()

	router.HandleFunc(`/job/{job_id}/{filepath:.+}`,
		func(w http.ResponseWriter, r *http.Request) {
			srv.saveHandler(w, r, mux.Vars(r))
		}).Methods("PUT").Name("save_job_artifact")

	router.HandleFunc(`/job/{job_id}`,
		func(w http.ResponseWriter, r *http.Request) {
			srv.listHandler(w, r, mux.Vars(r))
		}).Methods("GET").Name("list_job_artifacts")

	router.HandleFunc(`/job/{job_id}/{filepath:.+}`,
		func(w http.ResponseWriter, r *http.Request) {
			srv.getPathHandler(w, r, mux.Vars(r))
		}).Methods("GET").Name("get_job_artifact")

	srv.Router = router
}

func (srv *Server) setupNegroni() {
	srv.n = negroni.New()
	srv.n.Use(negroni.NewRecovery())
	srv.n.Use(NewLoggerMiddleware())
	srv.n.Use(negroni.NewStatic(http.Dir("public")))
	srv.n.UseHandler(srv.Router)
}

func (srv *Server) setupStorer() error {
	switch srv.opts.StorerType {
	case "s3":
		store, err := store.NewS3Store(srv.opts.S3Key,
			srv.opts.S3Secret, srv.opts.S3Bucket, srv.log, srv.md)
		if err != nil {
			return err
		}

		srv.store = store
		return nil
	case "file":
		srv.store = store.NewFileStore(srv.opts.FileStorePrefix,
			srv.log, srv.md)
		return nil
	default:
		srv.log.WithFields(logrus.Fields{
			"storer_type": srv.opts.StorerType,
		}).Error("unknown storer type")
		return fmt.Errorf("unknown storer type %q", srv.opts.StorerType)
	}

	return nil
}

func (srv *Server) setupAuther() error {
	srv.auth = &auth.TokenAuther{Token: srv.AuthToken}
	return nil
}

func (srv *Server) getDB() error {
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
	return true
}

func (srv *Server) canRead(r *http.Request, vars map[string]string) bool {
	return true
}
