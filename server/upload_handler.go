package server

import (
	"fmt"
	"net/http"
	"os"

	"github.com/meatballhat/artifacts-service/store"
)

type uploadHandler struct {
	Store store.Storer
}

func newUploadHandler() (*uploadHandler, error) {
	store, err := store.NewPostgreSQLStore(os.Getenv("HEROKU_POSTGRESQL_SILVER_URL"))
	if err != nil {
		return nil, err
	}

	return &uploadHandler{
		Store: store,
	}, nil
}

func (uh *uploadHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		uh.handlePostUpload(w, r)
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "whatever, meatbag\n")
	}
}

func (uh *uploadHandler) handlePostUpload(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "why not!?\n")
}
