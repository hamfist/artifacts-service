package auth

import (
	"net/http"
)

// Auther deals with authentication and authorization
type Auther interface {
	Check(*http.Request, map[string]string) *AuthResult
}

// AuthResult is what an Auther returns mkay
type AuthResult struct {
	UserID   string
	Resource string
	CanRead  bool
	CanWrite bool
	Errors   []error
}
