package xkom

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	fhttp "github.com/bogdanfinn/fhttp"

	"github.com/xMajkel/x-kom-unboxer/pkg/utility"
	"github.com/xMajkel/x-kom-unboxer/pkg/utility/config"
	"github.com/xMajkel/x-kom-unboxer/pkg/utility/shared"
)

type Boxes struct {
	BoxId        int `json:"BoxId"`
	Requirements []struct {
		IsMatched bool `json:"IsMatched"`
	} `json:"Requirements"`
	NextBoxOpeningPossibleDate string `json:"NextBoxOpeningPossibleDate"`
}

type ErrorResponse struct {
	Message string `json:"Message"`
}

// GetBoxes returns list of all boxes
func (acc *Account) GetBoxes() ([]Boxes, error) {
	var err error
	var req *fhttp.Request
	var resp *fhttp.Response
	var respJson []Boxes

	url := "https://mobileapi.x-kom.pl/api/v1/xkom/Box/Boxes"

	req, err = fhttp.NewRequest(fhttp.MethodGet, url, nil)
	if err != nil {
		return []Boxes{}, err
	}

	req.Header = map[string][]string{
		"x-api-key":       {config.GlobalConfig.ApiKey},
		"clientversion":   {"1.103.0"},
		"time-zone":       {"UTC"},
		"User-Agent":      {"x-kom_prod/20240916.1 CFNetwork/1496.0.7 Darwin/23.5.0"},
		"authorization":   {"Bearer " + acc.AccessToken},
		"accept-encoding": {"gzip"},
	}

	resp, err = acc.HttpClient.Do(req)
	if err != nil {
		return []Boxes{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		var errJson ErrorResponse

		body, err := utility.ReadHttpResponseBody(resp.Header.Get("Content-Encoding"), resp.Body)
		if err != nil {
			return []Boxes{}, errors.New(resp.Status)
		}
		defer body.Close()

		err = json.NewDecoder(body).Decode(&errJson)
		if err != nil {
			return []Boxes{}, errors.New(resp.Status)
		}
		if errJson.Message != "" {
			return []Boxes{}, errors.New(errJson.Message)
		}
		return []Boxes{}, errors.New(resp.Status)
	}

	body, err := utility.ReadHttpResponseBody(resp.Header.Get("Content-Encoding"), resp.Body)
	if err != nil {
		return []Boxes{}, err
	}
	defer body.Close()

	err = json.NewDecoder(body).Decode(&respJson)
	if err != nil {
		return []Boxes{}, err
	}

	return respJson, nil
}

type BoxItem struct {
	Item struct {
		Name  string `json:"Name"`
		Photo struct {
			Url          string `json:"Url"`
			ThumbnailUrl string `json:"ThumbnailUrl"`
		} `json:"Photo"`
		CatalogPrice         float64 `json:"CatalogPrice"`
		CategoryNameSingular string  `json:"CategoryNameSingular"`
	} `json:"Item"`
	BoxRarity struct {
		Id   string `json:"Id"`
		Name string `json:"Name"`
	} `json:"BoxRarity"`
	BoxPrice      float64 `json:"BoxPrice"`
	WebUrl        string  `json:"WebUrl"`
	ExpireDate    string  `json:"ExpireDate"`
	PromotionGain struct {
		Value     float64 `json:"Value"`
		GainValue string  `json:"GainValue"`
		GainType  string  `json:"GainType"`
	} `json:"PromotionGain"`
	ProductCommentsStatistics struct {
		TotalCount    int     `json:"TotalCount"`
		AverageRating float64 `json:"AverageRating"`
	} `json:"ProductCommentsStatistics"`
	NextBoxOpeningPossibleDate string `json:"NextBoxOpeningPossibleDate"`
}

// RollBox rolls box with given id
//
// Parameters:
//   - id: currently boxes id's are 1-3
//
// Returns:
//   - BoxItem: rolled item
//   - error: nil if error didn't occur
func (acc *Account) RollBox(id string) (BoxItem, error) {
	var err error
	var req *fhttp.Request
	var resp *fhttp.Response
	var respJson BoxItem

	if acc.AccessToken == "" {
		return respJson, shared.ErrNoAccessToken
	}

	url := "https://mobileapi.x-kom.pl/api/v1/xkom/Box/" + id + "/Roll"

	req, err = fhttp.NewRequest(fhttp.MethodPost, url, strings.NewReader(""))
	if err != nil {
		return respJson, err
	}

	req.Header = map[string][]string{
		"x-api-key":       {config.GlobalConfig.ApiKey},
		"clientversion":   {"1.103.0"},
		"time-zone":       {"UTC"},
		"User-Agent":      {"x-kom_prod/20240916.1 CFNetwork/1496.0.7 Darwin/23.5.0"},
		"authorization":   {"Bearer " + acc.AccessToken},
		"accept-encoding": {"gzip"},
	}

	resp, err = acc.HttpClient.Do(req)
	if err != nil {
		return respJson, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		if resp.StatusCode == 403 {
			return respJson, shared.ErrBoxNotYetAvailable
		}
		return respJson, fmt.Errorf(resp.Status)
	}

	body, err := utility.ReadHttpResponseBody(resp.Header.Get("Content-Encoding"), resp.Body)
	if err != nil {
		return respJson, err
	}
	defer body.Close()

	err = json.NewDecoder(body).Decode(&respJson)
	if err != nil {
		return respJson, err
	}

	return respJson, nil

}
