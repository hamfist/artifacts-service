package server

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/julienschmidt/httprouter"
	"github.com/meatballhat/artifacts-service/store"
)

// Server holds onto a router and a store
type Server struct {
	Router *httprouter.Router
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

	server, err := NewServer(opts)
	if err != nil {
		log.Fatal(err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "9839"
	}

	addr := fmt.Sprintf(":%s", port)
	log.Printf("artifacts-service listening on %v\n", addr)
	log.Fatal(http.ListenAndServe(addr, server.Router))
}

// NewServer creates a new *Server with a router and its routes registered
func NewServer(opts *Options) (*Server, error) {
	log := logrus.New()
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
	router := httprouter.New()

	router.GET(`/`, func(w http.ResponseWriter, r *http.Request, vars map[string]string) {
		srv.logRequest(r)
		srv.rootHandler(w, r, vars)
	})
	router.PUT(`/save`, func(w http.ResponseWriter, r *http.Request, vars map[string]string) {
		srv.logRequest(r)
		srv.saveHandler(w, r, vars)
	})
	router.GET(`/list/:owner/:repo`, func(w http.ResponseWriter, r *http.Request, vars map[string]string) {
		srv.logRequest(r)
		srv.listHandler(w, r, vars)
	})

	srv.Router = router
}

func (srv *Server) logRequest(r *http.Request) {
	srv.log.WithFields(logrus.Fields{
		"method": r.Method,
		"path":   r.URL.Path,
	}).Info("handling")
}

func (srv *Server) setupStorer() error {
	// TODO: make this work mkay?
	//pgstore, err := store.NewPostgreSQLStore(srv.opts.DatabaseURL)
	//if err != nil {
	//return err
	//}

	// s3store := store.NewS3Store()

	//multiStore := store.NewMultiStore()
	//multiStore.AddStorer("postgresql", pgstore)
	//multiStore.AddStorer("s3", s3store)

	//srv.store = multiStore
	srv.store = store.NewFileStore(srv.opts.FileStorePrefix)
	return nil
}
