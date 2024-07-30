package ghgist

import (
	"fmt"
	"github.com/kurtosis-tech/online-boutique-demo/src/currencyexternalapi/config"
	"net/url"
	"time"
)

const (
	apiBaseURL              = "https://gist.githubusercontent.com/leoporoli/"
	currenciesEndpointPath  = "4801500594b953e33fb87d2a34d31281/raw/dbb5537bf7f4cbe90cfca3f45fc7712d28a63944/currencies.json"
	latestRatesEndpointPath = "b84dc6e408cfeb4319840c1daf8bbc1f/raw/4257228f98aeb5ae1f0d4f7e258e9295c9c8cad8/latest.json"
)

var GHGistCurrencyAPIConfig = config.NewCurrencyAPIConfig(
	5*time.Second,
	getCurrenciesURL,
	getLatestRatesURL,
)

func getCurrenciesURL() (*url.URL, error) {
	currenciesEndpointUrlStr := fmt.Sprintf("%s%s", apiBaseURL, currenciesEndpointPath)

	currenciesEndpointUrl, err := url.Parse(currenciesEndpointUrlStr)
	if err != nil {
		return nil, err
	}

	return currenciesEndpointUrl, nil
}

func getLatestRatesURL(from string, to string) (*url.URL, error) {
	latestRatesEndpointUrlStr := fmt.Sprintf("%s%s", apiBaseURL, latestRatesEndpointPath)

	latestRatesEndpointUrl, err := url.Parse(latestRatesEndpointUrlStr)
	if err != nil {
		return nil, err
	}

	return latestRatesEndpointUrl, nil
}
