package server

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/julienschmidt/httprouter"
)

// Main is the top of the pile.  Start here.
func Main() {
	router := buildRouter()

	port := os.Getenv("PORT")
	if port == "" {
		port = "9839"
	}

	addr := fmt.Sprintf(":%s", port)
	log.Printf("artifacts-service listening on %v\n", addr)
	log.Fatal(http.ListenAndServe(addr, router))
}

func buildRouter() *httprouter.Router {
	router := httprouter.New()
	router.GET(`/`, rootHandler)
	router.POST(`/save`, saveHandler)
	router.GET(`/list/:owner/:repo`, listHandler)
	return router
}
