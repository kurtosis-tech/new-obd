// Copyright 2018 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/runtime/protoimpl"
	"html/template"

	"net"
	"net/http"
	"os"

	"strings"
)

type platformDetails struct {
	css      string
	provider string
}

// TODO fix
type Money struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The 3-letter currency code defined in ISO 4217.
	CurrencyCode string `protobuf:"bytes,1,opt,name=currency_code,json=currencyCode,proto3" json:"currency_code,omitempty"`
	// The whole units of the amount.
	// For example if `currencyCode` is `"USD"`, then 1 unit is one US dollar.
	Units int64 `protobuf:"varint,2,opt,name=units,proto3" json:"units,omitempty"`
	// Number of nano (10^-9) units of the amount.
	// The value must be between -999,999,999 and +999,999,999 inclusive.
	// If `units` is positive, `nanos` must be positive or zero.
	// If `units` is zero, `nanos` can be positive, zero, or negative.
	// If `units` is negative, `nanos` must be negative or zero.
	// For example $-1.75 is represented as `units`=-1 and `nanos`=-750,000,000.
	Nanos int32 `protobuf:"varint,3,opt,name=nanos,proto3" json:"nanos,omitempty"`
}

// TODO fix
type Product struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id          string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Name        string `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Description string `protobuf:"bytes,3,opt,name=description,proto3" json:"description,omitempty"`
	Picture     string `protobuf:"bytes,4,opt,name=picture,proto3" json:"picture,omitempty"`
	PriceUsd    *Money `protobuf:"bytes,5,opt,name=price_usd,json=priceUsd,proto3" json:"price_usd,omitempty"`
	// Categories such as "clothing" or "kitchen" that can be used to look up
	// other related products.
	Categories []string `protobuf:"bytes,6,rep,name=categories,proto3" json:"categories,omitempty"`
}

var (
	isCymbalBrand = strings.ToLower(os.Getenv("CYMBAL_BRANDING")) == "true"
	templates     = template.Must(template.New("").
			Funcs(template.FuncMap{
			"renderMoney":        renderMoney,
			"renderCurrencyLogo": renderCurrencyLogo,
		}).ParseGlob("templates/*.html"))
	plat platformDetails
)

var validEnvs = []string{"local", "gcp", "azure", "aws", "onprem", "alibaba"}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	//log := r.Context().Value(ctxKeyLog{}).(logrus.FieldLogger)//TODO fix
	//log.WithField("currency", currentCurrency(r)).Info("home")
	/*currencies, err := fe.getCurrencies(r.Context())
	if err != nil {
		renderHTTPError(log, r, w, errors.Wrap(err, "could not retrieve currencies"), http.StatusInternalServerError)
		return
	}*/
	//TODO fix
	currencies := []string{"USD"}
	//TODO fix
	/*products, err := fe.getProducts(r.Context())
	if err != nil {
		renderHTTPError(log, r, w, errors.Wrap(err, "could not retrieve products"), http.StatusInternalServerError)
		return
	}
	//TODO fix
	cart, err := fe.getCart(r.Context(), sessionID(r))
	if err != nil {
		renderHTTPError(log, r, w, errors.Wrap(err, "could not retrieve cart"), http.StatusInternalServerError)
		return
	}*/

	type productView struct {
		Item  *Product
		Price *Money
	}
	//TODO fix
	ps := []productView{
		{
			Item: &Product{
				Picture: "",
				Id:      "77",
				Name:    "fake",
				PriceUsd: &Money{
					CurrencyCode: "USD",
					Units:        1,
					Nanos:        0,
				},
			},
		},
	}

	// Set ENV_PLATFORM (default to local if not set; use env var if set; otherwise detect GCP, which overrides env)_
	var env = os.Getenv("ENV_PLATFORM")
	// Only override from env variable if set + valid env
	if env == "" || !stringinSlice(validEnvs, env) {
		fmt.Println("env platform is either empty or invalid")
		env = "local"
	}
	// Autodetect GCP
	addrs, err := net.LookupHost("metadata.google.internal.")
	if err == nil && len(addrs) >= 0 {
		logrus.Debugf("Detected Google metadata server: %v, setting ENV_PLATFORM to GCP.", addrs)
		env = "gcp"
	}

	logrus.Debugf("ENV_PLATFORM is: %s", env)
	plat = platformDetails{}
	plat.setPlatformDetails(strings.ToLower(env))

	if err := templates.ExecuteTemplate(w, "home", map[string]interface{}{
		"session_id":    sessionID(r),
		"request_id":    r.Context().Value(ctxKeyRequestID{}),
		"user_currency": currentCurrency(r),
		"show_currency": true,
		"currencies":    currencies,
		"products":      ps,
		"cart_size":     1,                         //TODO fix
		"banner_color":  os.Getenv("BANNER_COLOR"), // illustrates canary deployments
		//"ad":                fe.chooseAd(r.Context(), []string{}, log), //TODO fix
		"platform_css":    plat.css,
		"platform_name":   plat.provider,
		"is_cymbal_brand": isCymbalBrand,
		//"deploymentDetails": deploymentDetailsMap, //TODO fix
	}); err != nil {
		logrus.Error(err)
	}
}

func (plat *platformDetails) setPlatformDetails(env string) {
	if env == "aws" {
		plat.provider = "AWS"
		plat.css = "aws-platform"
	} else if env == "onprem" {
		plat.provider = "On-Premises"
		plat.css = "onprem-platform"
	} else if env == "azure" {
		plat.provider = "Azure"
		plat.css = "azure-platform"
	} else if env == "gcp" {
		plat.provider = "Google Cloud"
		plat.css = "gcp-platform"
	} else if env == "alibaba" {
		plat.provider = "Alibaba Cloud"
		plat.css = "alibaba-platform"
	} else {
		plat.provider = "local"
		plat.css = "local"
	}
}

func renderMoney(money *Money) string {
	currencyLogo := renderCurrencyLogo(money.CurrencyCode)
	return fmt.Sprintf("%s%d.%02d", currencyLogo, money.Units, money.Nanos/10000000)
}

func sessionID(r *http.Request) string {
	v := r.Context().Value(ctxKeySessionID{})
	if v != nil {
		return v.(string)
	}
	return ""
}

func currentCurrency(r *http.Request) string {
	c, _ := r.Cookie(cookieCurrency)
	if c != nil {
		return c.Value
	}
	return defaultCurrency
}

func renderCurrencyLogo(currencyCode string) string {
	logos := map[string]string{
		"USD": "$",
		"CAD": "$",
		"JPY": "¥",
		"EUR": "€",
		"TRY": "₺",
		"GBP": "£",
	}

	logo := "$" //default
	if val, ok := logos[currencyCode]; ok {
		logo = val
	}
	return logo
}

func stringinSlice(slice []string, val string) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}
