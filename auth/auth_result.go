package auth

import (
	"net/http"
)

// Result is what an Auther returns mkay
type Result struct {
	UserID   string
	Resource string
	CanRead  bool
	CanWrite bool
	Errors   []error
}

// NewResult initializes a *Result from an *http.Request and vars map
func NewResult(r *http.Request, vars map[string]string) *Result {
	return &Result{
		Errors:   []error{},
		UserID:   r.Header.Get("Artifacts-User"),
		Resource: vars["job_id"],
	}
}
