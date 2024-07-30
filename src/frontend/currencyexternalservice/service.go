package currencyexternalservice

import (
	"context"
	"github.com/kurtosis-tech/online-boutique-demo/src/currencyexternalapi"
)

type Service struct {
	primaryApi   *currencyexternalapi.CurrencyAPI
	secondaryApi *currencyexternalapi.CurrencyAPI
}

func NewService(primaryApi *currencyexternalapi.CurrencyAPI, secondaryApi *currencyexternalapi.CurrencyAPI) *Service {
	return &Service{primaryApi: primaryApi, secondaryApi: secondaryApi}
}

func (s *Service) GetSupportedCurrencies(ctx context.Context) ([]string, error) {

	var (
		currencyCodes []string
		err           error
	)

	currencyCodes, err = s.primaryApi.GetSupportedCurrencies(ctx)
	if err != nil {
		currencyCodes, err = s.secondaryApi.GetSupportedCurrencies(ctx)
		if err != nil {
			return nil, err
		}
	}

	return currencyCodes, nil
}

/*
func (s *Service) Convert(ctx context.Context, fromCode string, fromUnits int64, fromNanos int32, to string) (*fep.Money, error) {

	var (
		money = &fep.Money{}
		code  string
		units int64
		nanos int32
		err   error
	)

	code, units, nanos, err = s.secondaryApi.Convert(ctx, in.From.CurrencyCode, in.From.Units, in.From.Nanos, in.ToCode)
	if err != nil {
		code, units, nanos, err = s.secondaryApi.Convert(ctx, in.From.CurrencyCode, in.From.Units, in.From.Nanos, in.ToCode)
		if err != nil {
			return nil, err
		}
	}

	money.CurrencyCode = code
	money.Units = units
	money.Nanos = nanos

	return money, nil
}
*/
