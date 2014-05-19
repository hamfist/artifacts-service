package server

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

// Main is the top of the pile.  Start here.
func Main() {
	http.Handle(`/`, newRootHandler())

	uh, err := newUploadHandler()
	if err != nil {
		log.Fatalf("failed to build upload handler: %v\n", err)
	}
	http.Handle(`/upload`, uh)

	port := os.Getenv("PORT")
	if port == "" {
		port = "9839"
	}

	addr := fmt.Sprintf(":%s", port)
	log.Printf("artifacts-service listening on %v\n", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
