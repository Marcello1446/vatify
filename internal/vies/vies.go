package vies

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type VIESResult struct {
	IsValid     bool `json:"isValid"`
	CountryCode string
	TradeName   string `json:"name"`
	Address     string `json:"address"`
	CleanedVAT  string `json:"vatNumber"`
}

type VIESRequest struct {
	VAT         string
	CountryCode string
}

func Check(countryCode, VATNumber string) (*VIESResult, error) {
	url := fmt.Sprintf("https://ec.europa.eu/taxation_customs/vies/rest-api/ms/%s/vat/%s", countryCode, VATNumber)
	res, err := http.Get(url)

	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	v := &VIESResult{}
	v.CountryCode = countryCode

	errr := json.Unmarshal(data, &v)

	if errr != nil {
		return nil, errr
	}

	return v, nil

}
