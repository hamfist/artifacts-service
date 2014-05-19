package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "9839"
	}

	addr := fmt.Sprintf(":%s", port)

	log.Printf("artifacts-service listening on %v\n", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
