package middleware

import (
	"app2_http_api_database/auth"
	"context"
	"fmt"
	"net/http"
)

// --------------- Context Helpers ----------------
type ctxKey string

const userEntryKey ctxKey = "ldapUserEntry"

// Store LDAP entry in context
func contextWithUserEntry(ctx context.Context, entry interface{}) context.Context {
	return context.WithValue(ctx, userEntryKey, entry)
}

// Exported function to retrieve LDAP entry from context
func UserEntryFromContext(ctx context.Context) (interface{}, bool) {
	entry := ctx.Value(userEntryKey)
	if entry == nil {
		return nil, false
	}
	return entry, true
}

// ----------------- Middleware -----------------
func LDAPMiddleware(next http.HandlerFunc, cfg auth.LDAPConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()
		if !ok {
			// Missing Authorization header → trigger browser login popup
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
			return
		}

		entry, err := auth.AuthenticateUser(cfg, username, password)
		if err != nil {
			// Wrong credentials or cannot connect
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			http.Error(w, fmt.Sprintf("Authentication failed: %s", err.Error()), http.StatusUnauthorized)
			return
		}

		// Add LDAP entry to request context
		ctx := contextWithUserEntry(r.Context(), entry)
		r = r.WithContext(ctx)

		// Call the next handler
		next.ServeHTTP(w, r)
	}
}
