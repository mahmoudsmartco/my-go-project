package auth

import (
	"context"
	"errors"
	"fmt"
	"net"
	"time"

	"github.com/go-ldap/ldap/v3"
)

type contextKey string

const userEntryKey contextKey = "userEntry"

// تخزين LDAP entry في context
func ContextWithUserEntry(ctx context.Context, entry *ldap.Entry) context.Context {
	return context.WithValue(ctx, userEntryKey, entry)
}

// استرجاع LDAP entry من context
func UserEntryFromContext(ctx context.Context) (*ldap.Entry, bool) {
	entry, ok := ctx.Value(userEntryKey).(*ldap.Entry)
	return entry, ok
}

// ----------------- Errors -----------------
var (
	ErrLDAPUnavailable    = errors.New("cannot connect to LDAP server")
	ErrInvalidCredentials = errors.New("invalid username or password")
)

// ----------------- Config -----------------
type LDAPConfig struct {
	URL            string        // مثال: "ldap://localhost:389"
	BindDNTemplate string        // مثال: "uid=%s,ou=users,dc=mycompany,dc=com"
	ConnectTimeout time.Duration // مثال: 5 * time.Second
}

// ----------------- AuthenticateUser -----------------
func AuthenticateUser(cfg LDAPConfig, username, password string) (*ldap.Entry, error) {
	// Dial to LDAP server with timeout
	dialer := &net.Dialer{Timeout: cfg.ConnectTimeout}
	conn, err := ldap.DialURL(cfg.URL, ldap.DialWithDialer(dialer))
	if err != nil {
		return nil, ErrLDAPUnavailable
	}
	defer conn.Close()

	// DN of the user
	userDN := fmt.Sprintf(cfg.BindDNTemplate, username)

	// Try to bind (login) with user credentials
	err = conn.Bind(userDN, password)
	if err != nil {
		if ldapErr, ok := err.(*ldap.Error); ok && ldapErr.ResultCode == ldap.LDAPResultInvalidCredentials {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}

	// Fetch user info
	searchRequest := ldap.NewSearchRequest(
		userDN,
		ldap.ScopeBaseObject,
		ldap.NeverDerefAliases,
		0, 0, false,
		"(objectClass=*)",
		[]string{"dn", "cn", "mail"},
		nil,
	)

	result, err := conn.Search(searchRequest)
	if err != nil {
		return nil, err
	}
	if len(result.Entries) == 0 {
		return nil, ErrInvalidCredentials
	}

	return result.Entries[0], nil
}
