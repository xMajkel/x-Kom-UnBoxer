package xkom

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/xMajkel/x-kom-unboxer/pkg/utility"
	"github.com/xMajkel/x-kom-unboxer/pkg/utility/shared"
)

// Login logs in into account using refresh token or email and password and saves access_token for operations on account and refresh_token for further Login uses
func (acc *Account) Login() error {
	var err error
	var payload url.Values

	if acc.RefreshToken != "" {
		payload = url.Values{
			"grant_type":    {"refresh_token"},
			"refresh_token": {acc.RefreshToken},
			"client_id":     {"android"},
		}

		err = acc.postLogin(payload)
		if err == nil {
			return nil
		}
	}

	payload = url.Values{
		"grant_type": {"password"},
		"username":   {acc.Email},
		"password":   {acc.Password},
		"client_id":  {"android"},
		"scope":      {"api_v1 offline_access"},
	}

	return acc.postLogin(payload)
}

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func (acc *Account) postLogin(payload url.Values) error {
	var err error
	var req *http.Request
	var resp *http.Response
	var respJson LoginResponse

	url := "https://auth.x-kom.pl/xkom/Token"

	req, err = http.NewRequest(http.MethodPost, url, strings.NewReader(payload.Encode()))
	if err != nil {
		return err
	}

	req.Header = map[string][]string{
		"Content-Type": {"application/x-www-form-urlencoded"},
		"User-Agent":   {"xkom_prod/1.98.3"},
		"Host":         {"auth.x-kom.pl"},
	}

	resp, err = acc.HttpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf(resp.Status)
	}

	body, err := utility.ReadHttpResponseBody(resp.Header.Get("Content-Encoding"), resp.Body)
	if err != nil {
		return err
	}
	defer body.Close()

	err = json.NewDecoder(body).Decode(&respJson)
	if err != nil {
		return err
	}

	if respJson.AccessToken == "" || respJson.RefreshToken == "" {
		return shared.ErrNoAccessToken
	}

	acc.AccessToken = respJson.AccessToken
	acc.RefreshToken = respJson.RefreshToken

	return nil
}
