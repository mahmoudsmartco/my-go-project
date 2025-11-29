package routes

import (
	"app2_http_api_database/auth"
	v1 "app2_http_api_database/handler/v1"
	"app2_http_api_database/middleware"
	"net/http"
	"time"
)

func RegisterProtectedRoutes(mux *http.ServeMux) {

	// Protected route (LDAP)
	ldapCfg := auth.LDAPConfig{
		URL:            "ldap://localhost:389",
		BindDNTemplate: "uid=%s,ou=users,dc=mycompany,dc=com",
		ConnectTimeout: 5 * time.Second,
	}
	mux.Handle("/protected", middleware.FlexibleAuthMiddleware(http.HandlerFunc(v1.ProtectedHandler), ldapCfg))

}
