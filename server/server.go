package server

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/julienschmidt/httprouter"
	"github.com/meatballhat/artifacts-service/store"
)

// Server holds onto a router and a store
type Server struct {
	Router *httprouter.Router
	opts   *Options
	store  store.Storer
}

// Main is the top of the pile.  Start here.
func Main() {
	server, err := NewServer(NewOptions())
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
	server := &Server{opts: opts}
	server.setupRouter()
	err := server.setupStorer()
	if err != nil {
		return nil, err
	}

	return server, nil
}

func (srv *Server) setupRouter() {
	router := httprouter.New()

	router.GET(`/`, srv.rootHandler)
	router.POST(`/save`, srv.saveHandler)
	router.GET(`/list/:owner/:repo`, srv.listHandler)

	srv.Router = router
}

func (srv *Server) setupStorer() error {
	pgstore, err := store.NewPostgreSQLStore(srv.opts.DatabaseURL)
	if err != nil {
		return err
	}

	s3store := store.NewS3Store()

	multiStore := store.NewMultiStore()
	multiStore.AddStorer("postgresql", pgstore)
	multiStore.AddStorer("s3", s3store)

	srv.store = multiStore
	return nil
}
