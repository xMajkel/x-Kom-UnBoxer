package xkom

import (
	"flag"
	"testing"

	"github.com/xMajkel/x-kom-unboxer/pkg/utility/config"
)

var (
	email    = flag.String("email", "", "Email for account")
	password = flag.String("password", "", "Password for account")
	proxy    = flag.String("proxy", "", "Proxy URL")
	xapikey  = flag.String("xapikey", config.DEFAULT_API_KEY, "API Key")
)

func TestAccount(t *testing.T) {
	if *email == "" || *password == "" {
		t.Fatal("Flags required missing: -email, -password")
	}

	config.GlobalConfig.ApiKey = *xapikey

	acc, err := NewAccount(*email, *password, *proxy)
	if err != nil {
		t.Fatalf("could not create account: %v", err)
	}

	err = acc.Login()
	if err != nil {
		t.Fatalf("could not login: %v", err)
	}

	_, err = acc.GetBoxes()
	if err != nil {
		t.Fatalf("could not get boxes: %v", err)
	}
}
