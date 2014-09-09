package auth

import (
	"net/http"
)

// AuthResult is what an Auther returns mkay
type AuthResult struct {
	UserID   string
	Resource string
	CanRead  bool
	CanWrite bool
	Errors   []error
}

func NewAuthResult(r *http.Request, vars map[string]string) *AuthResult {
	return &AuthResult{
		Errors:   []error{},
		UserID:   r.Header.Get("Artifacts-User"),
		Resource: vars["job_id"],
	}
}
