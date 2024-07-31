package currencyexternalservice

import (
	"github.com/kurtosis-tech/new-obd/src/currencyexternalapi"
	"github.com/kurtosis-tech/new-obd/src/currencyexternalapi/config/freecurrency"
	"github.com/kurtosis-tech/new-obd/src/currencyexternalapi/config/ghgist"
)

func CreateService(apiKey string) *CurrencyExternalService {
	primaryApi := currencyexternalapi.NewCurrencyAPI(freecurrency.GetFreeCurrencyAPIConfig(apiKey))
	secondaryApi := currencyexternalapi.NewCurrencyAPI(ghgist.GHGistCurrencyAPIConfig)

	service := NewService(primaryApi, secondaryApi)

	return service
}
