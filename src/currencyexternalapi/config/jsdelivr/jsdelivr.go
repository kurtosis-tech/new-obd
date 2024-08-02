package jsdelivr

import (
	"encoding/json"
	"fmt"
	"github.com/kurtosis-tech/new-obd/src/currencyexternalapi/config"
	"net/url"
	"strings"
	"time"
)

const (
	apiBaseURL              = "https://cdn.jsdelivr.net/npm/@fawazahmed0/currency-api@latest/v1/"
	currenciesEndpointPath  = "currencies.json"
	latestRatesEndpointPath = "currencies/usd.json"
)

type LatestRatesResponse struct {
	Date string             `json:"date"`
	Usd  map[string]float64 `json:"usd"`
}

var JsdelivrAPIConfig = config.NewCurrencyAPIConfig(
	// saving the response for a week because app.freecurrencyapi.com has a low limit
	// and this is a demo project, it's not important to have the latest data
	5*time.Second,
	getGetCurrenciesURLFunc,
	getGetLatestRatesURLFunc,
	getCurrencyListFromResponseFunc,
	getLatestRatesFromResponse,
)

func getGetCurrenciesURLFunc() (*url.URL, error) {

	currenciesEndpointUrlStr := fmt.Sprintf("%s%s", apiBaseURL, currenciesEndpointPath)

	currenciesEndpointUrl, err := url.Parse(currenciesEndpointUrlStr)
	if err != nil {
		return nil, err
	}

	return currenciesEndpointUrl, nil
}

func getGetLatestRatesURLFunc(from string, to string) (*url.URL, error) {

	latestRatesEndpointUrlStr := fmt.Sprintf("%s%s", apiBaseURL, latestRatesEndpointPath)

	latestRatesEndpointUrl, err := url.Parse(latestRatesEndpointUrlStr)
	if err != nil {
		return nil, err
	}

	return latestRatesEndpointUrl, nil
}

func getCurrencyListFromResponseFunc(httpResponseBodyBytes []byte) ([]string, error) {
	currencyCodes := []string{}
	currenciesResp := &map[string]string{}
	if err := json.Unmarshal(httpResponseBodyBytes, currenciesResp); err != nil {
		return currencyCodes, err
	}

	for code := range *currenciesResp {
		upperCode := strings.ToUpper(code)
		currencyCodes = append(currencyCodes, upperCode)
	}
	return currencyCodes, nil
}

func getLatestRatesFromResponse(httpResponseBodyBytes []byte) (map[string]float64, error) {

	data := map[string]float64{}
	latestRatesResp := &LatestRatesResponse{}
	if err := json.Unmarshal(httpResponseBodyBytes, latestRatesResp); err != nil {
		return data, err
	}
	data = latestRatesResp.Usd
	dataUpperCode := map[string]float64{}
	for code, rate := range data {
		upperCode := strings.ToUpper(code)
		dataUpperCode[upperCode] = rate
	}

	//add USD
	dataUpperCode["USD"] = 1
	return dataUpperCode, nil
}
