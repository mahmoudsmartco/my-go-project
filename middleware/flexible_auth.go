package middleware

import (
	"app2_http_api_database/auth"
	"fmt"
	"net/http"
	"strings"
)

func FlexibleAuthMiddleware(next http.Handler, ldapCfg auth.LDAPConfig) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// ---------------- JWT check ----------------
		tokenStr := r.Header.Get("Authorization")
		tokenStr = strings.TrimPrefix(tokenStr, "Bearer ")
		if tokenStr != "" {
			if _, err := auth.ValidateJWT(tokenStr); err == nil {
				next.ServeHTTP(w, r)
				return
			}
		}

		// ---------------- LDAP check ----------------
		username, password, ok := r.BasicAuth()
		if ok && username != "" && password != "" {
			entry, err := auth.AuthenticateUser(ldapCfg, username, password)
			if err == nil {
				// ممكن نخزن الـ entry في الـ context لو محتاجين بعد كده
				ctx := r.Context()
				ctx = auth.ContextWithUserEntry(ctx, entry)
				r = r.WithContext(ctx)

				next.ServeHTTP(w, r)
				return
			} else {
				fmt.Println("LDAP error:", err)
			}
		}

		// لو مفيش JWT صحيح أو LDAP صحيح
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	})
}
