package freecurrency

import (
	"fmt"
	"github.com/kurtosis-tech/online-boutique-demo/src/currencyexternalapi/config"
	"net/url"
	"strings"
	"time"
)

const (
	apiBaseURL              = "https://api.freecurrencyapi.com/v1/"
	apiKeyQueryParamKey     = "apikey"
	currenciesQueryParamKey = "currencies"
	currenciesEndpointPath  = "currencies"
	latestRatesEndpointPath = "latest"
)

func GetFreeCurrencyAPIConfig(apiKey string) *config.CurrencyAPIConfig {
	var FreeCurrencyAPIConfig = config.NewCurrencyAPIConfig(
		// saving the response for a week because app.freecurrencyapi.com has a low limit
		// and this is a demo project, it's not important to have the latest data
		168*time.Hour,
		getGetCurrenciesURLFunc(apiKey),
		getGetLatestRatesURLFunc(apiKey),
	)
	return FreeCurrencyAPIConfig
}

func getGetCurrenciesURLFunc(apiKey string) func() (*url.URL, error) {

	getCurrenciesURLFunc := func() (*url.URL, error) {
		currenciesEndpointUrlStr := fmt.Sprintf("%s%s", apiBaseURL, currenciesEndpointPath)

		currenciesEndpointUrl, err := url.Parse(currenciesEndpointUrlStr)
		if err != nil {
			return nil, err
		}

		currenciesEndpointQuery := currenciesEndpointUrl.Query()

		currenciesEndpointQuery.Set(apiKeyQueryParamKey, apiKey)

		currenciesEndpointUrl.RawQuery = currenciesEndpointQuery.Encode()

		return currenciesEndpointUrl, nil
	}

	return getCurrenciesURLFunc
}

func getGetLatestRatesURLFunc(apiKey string) func(string, string) (*url.URL, error) {

	getLatestRatesURLFunc := func(from string, to string) (*url.URL, error) {
		latestRatesEndpointUrlStr := fmt.Sprintf("%s%s", apiBaseURL, latestRatesEndpointPath)

		latestRatesEndpointUrl, err := url.Parse(latestRatesEndpointUrlStr)
		if err != nil {
			return nil, err
		}

		latestRatesEndpointQuery := latestRatesEndpointUrl.Query()

		currenciesQueryParamValue := strings.Join([]string{strings.ToUpper(from), strings.ToUpper(to)}, ",")

		latestRatesEndpointQuery.Set(apiKeyQueryParamKey, apiKey)
		latestRatesEndpointQuery.Set(currenciesQueryParamKey, currenciesQueryParamValue)

		latestRatesEndpointUrl.RawQuery = latestRatesEndpointQuery.Encode()

		return latestRatesEndpointUrl, nil
	}

	return getLatestRatesURLFunc
}
