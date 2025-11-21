package imap

import (
	"crypto/tls"
	"log"
	
	"github.com/emersion/go-imap/client"
)

// Connect establishes a connection to the IMAP server and logs in.
// addr: "hostname:port"
// useTLS: true to use TLS (usually port 993), false for plain TCP (usually 143).
func Connect(addr, username, password string, useTLS bool) (*client.Client, error) {
	var c *client.Client
	var err error

	if useTLS {
		// For production, we might want to allow custom TLS config (e.g. skip verify)
		// For now, use default.
		c, err = client.DialTLS(addr, nil)
	} else {
		c, err = client.Dial(addr)
	}

	if err != nil {
		return nil, err
	}

	        if err := c.Login(username, password); err != nil {
	                // If login fails, try to logout (best effort), then return the login error.
	                if logoutErr := c.Logout(); logoutErr != nil {
	                        log.Printf("Error logging out from IMAP after login failure: %v", logoutErr)
	                }
	                return nil, err
	        }
	return c, nil
}

// ConnectWithConfig allows passing a custom TLS config if needed (e.g. for self-signed certs)
func ConnectWithConfig(addr, username, password string, tlsConfig *tls.Config) (*client.Client, error) {
    // MVP extension: Not needed yet, but good to have in mind.
    // adhering to YAGNI, I won't implement it yet.
    return nil, nil 
}
