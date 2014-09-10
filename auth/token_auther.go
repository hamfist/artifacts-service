package auth

import (
	"net/http"
)

// TokenAuther implements Auther with tokens!
type TokenAuther struct {
	AuthToken string
}

// NewTokenAuther makes a new *TokenAuther wow!
func NewTokenAuther(authToken string) *TokenAuther {
	return &TokenAuther{AuthToken: authToken}
}

// Check checks the token mkay
func (ta *TokenAuther) Check(r *http.Request, vars map[string]string) *Result {
	ar := NewResult(r, vars)

	if r.Header.Get("Authorization") == ("token "+ta.AuthToken) ||
		r.Header.Get("Authorization") == ("token="+ta.AuthToken) {
		ar.CanRead = true
		ar.CanWrite = true
	}

	return ar
}
