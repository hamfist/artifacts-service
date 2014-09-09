package auth

import (
	"net/http"
)

// Auther deals with authentication and authorization
type Auther interface {
	Check(*http.Request, map[string]string) *AuthResult
}
