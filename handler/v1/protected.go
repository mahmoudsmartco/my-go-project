package v1

import (
	"app2_http_api_database/auth"
	"fmt"
	"net/http"
)

// ---------------- Protected Handler ----------------
func ProtectedHandler(w http.ResponseWriter, r *http.Request) {
	entry, ok := auth.UserEntryFromContext(r.Context())
	if !ok {
		http.Error(w, "No user info in context", http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "Welcome by LDAP, %s!", entry.DN)
}
