package auth

import (
	"net/http"
)

// TokenAuther implements Auther with tokens!
type TokenAuther struct {
	AuthToken string
}

// Check checks the token mkay
func (ta *TokenAuther) Check(r *http.Request, vars map[string]string) *AuthResult {
	ar := NewAuthResult(r, vars)

	if r.Header.Get("Authorization") == ("token "+ta.AuthToken) ||
		r.Header.Get("Authorization") == ("token="+ta.AuthToken) {
		ar.CanRead = true
		ar.CanWrite = true
	}

	return ar
}
