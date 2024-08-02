package currencyexternalservice

import (
	"github.com/kurtosis-tech/new-obd/src/currencyexternalapi"
	"github.com/kurtosis-tech/new-obd/src/currencyexternalapi/config/freecurrency"
	"github.com/kurtosis-tech/new-obd/src/currencyexternalapi/config/jsdelivr"
)

func CreateService(apiKey string) *CurrencyExternalService {
	primaryApi := currencyexternalapi.NewCurrencyAPI(jsdelivr.JsdelivrAPIConfig)
	secondaryApi := currencyexternalapi.NewCurrencyAPI(freecurrency.GetFreeCurrencyAPIConfig(apiKey))

	service := NewService(primaryApi, secondaryApi)

	return service
}
