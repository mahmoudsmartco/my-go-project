// config/ldap_config.go
package config

import (
	"os"
	"time"
)

type LDAPConfig struct {
	URL            string        // example: "ldap://localhost:389"
	BaseDN         string        // example: "dc=mycompany,dc=com"
	UserDNPattern  string        // example: "uid=%s,ou=users,dc=mycompany,dc=com"
	AdminDN        string        // example: "cn=admin,dc=mycompany,dc=com"
	AdminPassword  string        // admin password
	ConnectTimeout time.Duration // e.g. 5 * time.Second
}

func GetLDAPConfig() LDAPConfig {
	// قيم افتراضية ممكن تغيرها عبر env
	url := getEnv("LDAP_URL", "ldap://localhost:389")
	baseDN := getEnv("LDAP_BASE_DN", "dc=mycompany,dc=com")
	userPattern := getEnv("LDAP_USER_DN_PATTERN", "uid=%s,ou=users,dc=mycompany,dc=com")
	adminDN := getEnv("LDAP_ADMIN_DN", "cn=admin,dc=mycompany,dc=com")
	adminPass := getEnv("LDAP_ADMIN_PASSWORD", "admin")
	timeoutStr := getEnv("LDAP_TIMEOUT", "5s")

	timeout, err := time.ParseDuration(timeoutStr)
	if err != nil {
		timeout = 5 * time.Second
	}

	return LDAPConfig{
		URL:            url,
		BaseDN:         baseDN,
		UserDNPattern:  userPattern,
		AdminDN:        adminDN,
		AdminPassword:  adminPass,
		ConnectTimeout: timeout,
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
