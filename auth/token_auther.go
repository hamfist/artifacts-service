package auth

import (
	"net/http"
)

// TokenAuther implements Auther with tokens!
type TokenAuther struct {
	Token string
}

// Check checks the token mkay
func (ta *TokenAuther) Check(r *http.Request, vars map[string]string) *AuthResult {
	ar := &AuthResult{
		Errors:   []error{},
		UserID:   r.Header.Get("Artifacts-User"),
		Resource: vars["slug"],
	}

	if r.Header.Get("Authorization") == ("token " + ta.Token) {
		ar.CanRead = true
		ar.CanWrite = true
	}

	return ar
}
