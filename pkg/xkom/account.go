package xkom

import (
	"errors"

	tls_client "github.com/bogdanfinn/tls-client"
	"github.com/bogdanfinn/tls-client/profiles"
)

type Account struct {
	Email    string
	Password string

	AccessToken  string
	RefreshToken string

	HttpClient tls_client.HttpClient
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
		return nil, errors.New("no email or password provided")
	}

	client, err := tls_client.NewHttpClient(tls_client.NewNoopLogger(),
		tls_client.WithClientProfile(profiles.Safari_IOS_18_0),
		tls_client.WithForceHttp1(),
	)

	if err != nil {
		return nil, err
	}

	//Set proxy
	if proxy != "" {
		err := client.SetProxy(proxy)
		if err != nil {
			return nil, err
		}
	}

	return &Account{
		Email:      email,
		Password:   password,
		HttpClient: client,
	}, nil
}
