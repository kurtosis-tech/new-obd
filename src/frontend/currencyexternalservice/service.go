package currencyexternalservice

import (
	"context"
	fep "github.com/kurtosis-tech/online-boutique-demo/frontend/proto"
	"github.com/kurtosis-tech/online-boutique-demo/src/currencyexternalapi"
	"go-micro.dev/v4/client"
)

type Service struct {
	primaryApi   *currencyexternalapi.CurrencyAPI
	secondaryApi *currencyexternalapi.CurrencyAPI
}

func NewService(primaryApi *currencyexternalapi.CurrencyAPI, secondaryApi *currencyexternalapi.CurrencyAPI) *Service {
	return &Service{primaryApi: primaryApi, secondaryApi: secondaryApi}
}

func (s *Service) GetSupportedCurrencies(ctx context.Context, _ *fep.Empty, _ ...client.CallOption) (*fep.GetSupportedCurrenciesResponse, error) {

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

	response := &fep.GetSupportedCurrenciesResponse{
		CurrencyCodes: currencyCodes,
	}

	return response, nil
}

func (s *Service) Convert(ctx context.Context, in *fep.CurrencyConversionRequest, _ ...client.CallOption) (*fep.Money, error) {

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
