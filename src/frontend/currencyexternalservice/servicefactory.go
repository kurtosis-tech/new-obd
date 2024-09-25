package currencyexternalservice

import (
	"github.com/kurtosis-tech/new-obd/src/libraries/currencyexternalapi"
	"github.com/kurtosis-tech/new-obd/src/libraries/currencyexternalapi/config/jsdelivr"
)

func CreateService(apiKey string) *CurrencyExternalService {
	primaryApi := currencyexternalapi.NewCurrencyAPI(jsdelivr.GetJsdelivrAPIConfig(apiKey))
	service := NewService(primaryApi)
	return service
}
