package xkom

import (
	"crypto/tls"
	"errors"
	"net/http"
	"net/url"
)

type Account struct {
	Email    string
	Password string

	AccessToken  string
	RefreshToken string

	HttpClient *http.Client
}

// NewAccount creates new Account struct containing account information
//
// Parameters:
//   - email: reqired
//   - password: required
//   - proxy: optional proxy address, format: http://host:port:username:password
//
// Returns:
//   - *Account: pointer to Account
//   - error: nil if error didn't occur
func NewAccount(email, password, proxy string) (*Account, error) {
	if email == "" || password == "" {
		return &Account{}, errors.New("no email or password provided")
	}

	//Force HTTP/1.1 protocol
	transport := &http.Transport{
		TLSNextProto: make(map[string]func(authority string, c *tls.Conn) http.RoundTripper),
	}

	//Set proxy
	if proxy != "" {
		proxyURL, err := url.Parse(proxy)
		if err != nil {
			return &Account{}, err
		}

		transport.Proxy = http.ProxyURL(proxyURL)
	}

	client := &http.Client{
		Transport: transport,
	}
	return &Account{
		Email:      email,
		Password:   password,
		HttpClient: client,
	}, nil
}
