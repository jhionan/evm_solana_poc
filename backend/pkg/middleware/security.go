// Package middleware provides reusable HTTP middleware.
package middleware

import "net/http"

// SecurityHeaders is an HTTP middleware that adds security-related response
// headers to every outgoing response.
func SecurityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h := w.Header()

		// Prevent MIME-type sniffing.
		h.Set("X-Content-Type-Options", "nosniff")

		// Disallow embedding in frames/iframes.
		h.Set("X-Frame-Options", "DENY")

		// Enable browser's XSS filter (legacy, still useful).
		h.Set("X-XSS-Protection", "1; mode=block")

		// Force HTTPS for 2 years, include subdomains.
		h.Set("Strict-Transport-Security", "max-age=63072000; includeSubDomains; preload")

		// Minimal CSP — restrict everything to same origin by default.
		h.Set("Content-Security-Policy", "default-src 'self'")

		// Never send the Referer header to other origins.
		h.Set("Referrer-Policy", "no-referrer")

		next.ServeHTTP(w, r)
	})
}
