package server

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/julienschmidt/httprouter"
)

// Server holds onto a router and a store
type Server struct {
	Router *httprouter.Router
}

// Main is the top of the pile.  Start here.
func Main() {
	server := buildServer()

	port := os.Getenv("PORT")
	if port == "" {
		port = "9839"
	}

	addr := fmt.Sprintf(":%s", port)
	log.Printf("artifacts-service listening on %v\n", addr)
	log.Fatal(http.ListenAndServe(addr, server.Router))
}

func buildServer() *Server {
	server := &Server{}
	router := httprouter.New()

	router.GET(`/`, server.rootHandler)
	router.POST(`/save`, server.saveHandler)
	router.GET(`/list/:owner/:repo`, server.listHandler)

	server.Router = router

	return server
}
