package currencyexternalapi

import (
	"context"
	"encoding/json"
	"github.com/kurtosis-tech/online-boutique-demo/src/currencyexternalapi/config"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
	"math"
	"net/http"
	"strings"
)

type CurrenciesResponse struct {
	Data map[string]Currency `json:"data"`
}

type Currency struct {
	Symbol        string `json:"symbol"`
	Name          string `json:"name"`
	SymbolNative  string `json:"symbol_native"`
	DecimalDigits int    `json:"decimal_digits"`
	Rounding      int    `json:"rounding"`
	Code          string `json:"code"`
	NamePlural    string `json:"name_plural"`
	Type          string `json:"type"`
}

type LatestRatesResponse struct {
	Data LatestRates `json:"data"`
}

type LatestRates map[string]float64

type CurrencyAPI struct {
	httpClient *http.Client
	cache      *Cache
	config     *config.CurrencyAPIConfig
}

func NewCurrencyAPI(config *config.CurrencyAPIConfig) *CurrencyAPI {
	return &CurrencyAPI{httpClient: http.DefaultClient, cache: NewCache(), config: config}
}

func (c *CurrencyAPI) GetSupportedCurrencies(ctx context.Context) ([]string, error) {

	currenciesURL, err := c.config.GetCurrenciesURLFunc()
	if err != nil {
		return nil, err
	}

	httpRequest := &http.Request{
		Method: http.MethodGet,
		URL:    currenciesURL,
	}
	httpRequestWithContext := httpRequest.WithContext(ctx)

	httpResponseBodyBytes, err := c.doHttpRequest(httpRequestWithContext)
	if err != nil {
		return nil, err
	}

	currenciesResp := &CurrenciesResponse{}
	if err = json.Unmarshal(httpResponseBodyBytes, currenciesResp); err != nil {
		return nil, err
	}

	currencyCodes := []string{}

	for code := range currenciesResp.Data {
		currencyCodes = append(currencyCodes, code)
	}

	return currencyCodes, nil
}

func (c *CurrencyAPI) Convert(ctx context.Context, fromCode string, fromUnits int64, fromNanos int32, to string) (string, int64, int32, error) {

	fromCode = strings.ToUpper(fromCode)
	toCode := strings.ToUpper(to)

	currencies, err := c.getLatestRatesFromAPI(ctx, fromCode, toCode)
	if err != nil {
		return "", 0, 0, err
	}
	fromCurrency, found := currencies[fromCode]
	if !found {
		return "", 0, 0, status.Errorf(codes.InvalidArgument, "unsupported currency: %s", fromCode)
	}
	toCurrency, found := currencies[toCode]
	if !found {
		return "", 0, 0, status.Errorf(codes.InvalidArgument, "unsupported currency: %s", toCode)
	}

	total := int64(math.Floor(float64(fromUnits*10^9+int64(fromNanos)) / fromCurrency * toCurrency))
	units := total / 1e9
	nanos := int32(total % 1e9)

	return toCode, units, nanos, nil
}

func (c *CurrencyAPI) getLatestRatesFromAPI(ctx context.Context, from string, to string) (map[string]float64, error) {

	latestRatesEndpointUrl, err := c.config.GetLatestRatesURLFunc(from, to)
	if err != nil {
		return nil, err
	}

	httpRequest := &http.Request{
		Method: http.MethodGet,
		URL:    latestRatesEndpointUrl,
	}
	httpRequestWithContext := httpRequest.WithContext(ctx)

	httpResponseBodyBytes, err := c.doHttpRequest(httpRequestWithContext)
	if err != nil {
		return nil, err
	}

	latestRatesResp := &LatestRatesResponse{}
	if err = json.Unmarshal(httpResponseBodyBytes, latestRatesResp); err != nil {
		return nil, err
	}

	return latestRatesResp.Data, nil
}

func (c *CurrencyAPI) doHttpRequest(
	request *http.Request,
) (
	resultResponseBodyBytes []byte,
	resultErr error,
) {

	var (
		httpResponseBodyBytes []byte
		err                   error
		ok                    bool
		urlStr                = request.URL.String()
	)

	if httpResponseBodyBytes, ok = c.cache.Get(urlStr); ok {
		logrus.Debugf("Cache hit for '%s'", urlStr)
		return httpResponseBodyBytes, nil
	}

	httpResponse, err := c.httpClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer httpResponse.Body.Close()

	if httpResponse.StatusCode == http.StatusOK {
		httpResponseBodyBytes, err = io.ReadAll(httpResponse.Body)
		if err != nil {
			return nil, err
		}
	}

	c.cache.Set(urlStr, httpResponseBodyBytes, c.config.CacheDuration)

	return httpResponseBodyBytes, nil
}
