package currencyexternalservice

import (
	"github.com/kurtosis-tech/online-boutique-demo/src/currencyexternalapi"
	"github.com/kurtosis-tech/online-boutique-demo/src/currencyexternalapi/config/freecurrency"
	"github.com/kurtosis-tech/online-boutique-demo/src/currencyexternalapi/config/ghgist"
)

func CreateService(apiKey string) *Service {
	primaryApi := currencyexternalapi.NewCurrencyAPI(freecurrency.GetFreeCurrencyAPIConfig(apiKey))
	secondaryApi := currencyexternalapi.NewCurrencyAPI(ghgist.GHGistCurrencyAPIConfig)

	service := NewService(primaryApi, secondaryApi)

	return service
}
