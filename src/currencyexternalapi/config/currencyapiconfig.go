package config

import (
	"net/url"
	"time"
)

type CurrencyAPIConfig struct {
	CacheDuration         time.Duration
	GetCurrenciesURLFunc  func() (*url.URL, error)
	GetLatestRatesURLFunc func(from string, to string) (*url.URL, error)
}

func NewCurrencyAPIConfig(
	cacheDuration time.Duration,
	getCurrenciesURLFunc func() (*url.URL, error),
	getLatestRatesURLFunc func(from string, to string) (*url.URL, error),
) *CurrencyAPIConfig {
	return &CurrencyAPIConfig{
		CacheDuration:         cacheDuration,
		GetCurrenciesURLFunc:  getCurrenciesURLFunc,
		GetLatestRatesURLFunc: getLatestRatesURLFunc,
	}
}
