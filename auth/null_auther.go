package auth

import (
	"net/http"
)

// NullAuther implements Auther but doesn't even care!
type NullAuther struct{}

// Check always allows reads and writes
func (na *NullAuther) Check(r *http.Request, vars map[string]string) *AuthResult {
	ar := NewAuthResult(r, vars)
	ar.CanRead = true
	ar.CanWrite = true
	return ar
}
