package api

import (
	"crypto/subtle"
	"net/http"
)

// adminUser is the username for the admin login.
const adminUser = "admin"

// AdminOnly is middleware: it wraps another handler and only lets the
// request through if it carries the admin username and password
// (HTTP Basic Auth). Otherwise it replies 401 Unauthorized, and the
// WWW-Authenticate header makes the browser show its login box.
//
// next is an http.Handler: an interface with one method, ServeHTTP.
// Anything that has that method can be wrapped, whether it's one of our
// handlers or Go's file server.
func AdminOnly(password string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, pass, ok := r.BasicAuth()
		if !ok || !sameText(user, adminUser) || !sameText(pass, password) {
			w.Header().Set("WWW-Authenticate", `Basic realm="MacBook Marketplace admin", charset="UTF-8"`)
			writeError(w, http.StatusUnauthorized, "admin login required")
			return
		}
		next.ServeHTTP(w, r) // logged in: run the real handler
	})
}

// sameText compares two strings in constant time. A normal == stops at the
// first different letter, and an attacker could time that to guess the
// password one letter at a time.
func sameText(a, b string) bool {
	return subtle.ConstantTimeCompare([]byte(a), []byte(b)) == 1
}
