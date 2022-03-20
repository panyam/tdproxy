package models

import (
	"fmt"
	"github.com/panyam/goutils/utils"
	"log"
	"sort"
	"strconv"
	"strings"
	"time"
)

type Auth struct {
	ClientId       string `gorm:"primaryKey"`
	CallbackUrl    string
	CreatedAt      time.Time `gorm:"autoCreateTime"`
	UpdatedAt      time.Time `gorm:"autoUpdateTime"`
	ExpiresAt      time.Time
	AuthToken      utils.StringMap
	UserPrincipals utils.StringMap
}

func (a *Auth) ToJson() utils.StringMap {
	out := make(utils.StringMap)
	out["client_id"] = a.ClientId
	out["callback_url"] = a.CallbackUrl
	out["auth_token"] = a.AuthToken
	out["user_principals"] = a.UserPrincipals
	out["expires_at"] = utils.FormatTime(a.ExpiresAt)
	return out
}

func (auth *Auth) FromJson(json utils.StringMap) {
	if json != nil {
		auth.ClientId = json["client_id"].(string)
		auth.CallbackUrl = json["callback_url"].(string)
		if val, ok := json["auth_token"]; ok && val != nil {
			auth.AuthToken = val.(utils.StringMap)
		}
		if val, ok := json["user_principals"]; ok && val != nil {
			auth.UserPrincipals = val.(utils.StringMap)
		}
		if val, ok := json["expires_at"]; ok && val != nil {
			auth.ExpiresAt = utils.ParseTime(val.(string))
		}
	}
}

func (auth *Auth) IsAuthenticated() bool {
	if auth.AuthToken == nil {
		return false
	}
	if auth.ExpiresAt.Sub(time.Now().UTC()) <= 0 {
		return false
	}
	return true
}

func (auth *Auth) AccessToken() string {
	access_token := auth.AuthToken["access_token"]
	if access_token == nil {
		return ""
	}
	return access_token.(string)
}

type Ticker struct {
	Symbol          string `gorm:"primaryKey"`
	LastRefreshedAt time.Time
	Info            map[string]interface{} `gorm:"-"`
}

func NewTicker(symbol string, refreshed_at time.Time, info map[string]interface{}) *Ticker {
	return &Ticker{
		Symbol:          symbol,
		LastRefreshedAt: refreshed_at,
		Info:            info,
	}
}

type Option struct {
	Symbol       string `gorm:"primaryKey" gorm:"index:ByCallSymbolDate,priority:1"`
	DateString   string `gorm:"primaryKey" gorm:"index:ByCallSymbolDate,priority:2"`
	IsCall       bool   `gorm:"primaryKey" gorm:"index:ByCallSymbolDate,priority:3"`
	PriceString  string `gorm:"primaryKey"`
	StrikePrice  float64
	AskPrice     float64
	BidPrice     float64
	MarkPrice    float64
	OpenInterest int32
	Delta        float64
	Multiplier   float64
	Info         map[string]interface{} `gorm:"-"`
}

func NewOption(symbol string, date_string string, price_string string, is_call bool, info map[string]interface{}) *Option {
	out := &Option{
		Symbol:      symbol,
		DateString:  date_string,
		PriceString: price_string,
		IsCall:      is_call,
	}
	out.Refresh()
	return out
}

func (opt *Option) Refresh() bool {
	if val, ok := opt.Info["ask"]; ok {
		opt.AskPrice = val.(float64)
	}
	if val, ok := opt.Info["bid"]; ok {
		opt.BidPrice = val.(float64)
	}
	if val, ok := opt.Info["mark"]; ok {
		opt.MarkPrice = val.(float64)
	} else {
		opt.MarkPrice = (opt.AskPrice + opt.BidPrice) / 2
	}

	if val, ok := opt.Info["openInterest"]; ok {
		opt.OpenInterest = int32(val.(float64))
	}
	if val, ok := opt.Info["delta"]; ok {
		defer func() {
			if err := recover(); err != nil {
				log.Println("Panic Occurred: ", err, val)
			}
		}()
		opt.Delta = val.(float64)
	}
	if val, ok := opt.Info["multiplier"]; ok {
		opt.Multiplier = val.(float64)
	}
	if opt.StrikePrice <= 0 {
		result, err := strconv.ParseFloat(opt.PriceString, 64)
		if err != nil {
			fmt.Printf("Invalid price string: %s\n", opt.PriceString)
		}
		opt.StrikePrice = result
	}
	return true
}

type ChainInfo struct {
	Symbol          string   `gorm:"primaryKey"`
	AvailableDates  []string `gorm:"type:integer[]"`
	LastRefreshedAt time.Time
}

type Chain struct {
	Symbol          string `gorm:"primaryKey"`
	DateString      string `gorm:"primaryKey"`
	IsCall          bool   `gorm:"primaryKey"`
	LastRefreshedAt time.Time
	Options         []*Option `gorm:"-"` // dont read/write this
}

func ChainFromDict(symbol string, date string, is_call bool,
	options_by_price map[string]interface{},
	refreshed_at time.Time) *Chain {
	var options []*Option

	var detail map[string]interface{}
	for price_string, option_details := range options_by_price {
		detail_list, ok := option_details.([]interface{})
		if ok {
			detail = detail_list[0].(map[string]interface{})
		} else {
			detail = option_details.(map[string]interface{})
		}
		option := Option{
			Symbol:      symbol,
			DateString:  date,
			IsCall:      is_call,
			PriceString: price_string,
			Info:        detail,
		}
		if !option.Refresh() {
			log.Println("Refresh failed: ", option)
		} else {
			options = append(options, &option)
		}
	}
	chain := NewChain(symbol, date, is_call, options)
	chain.LastRefreshedAt = refreshed_at
	return chain
}

func NewChain(symbol string, date string, is_call bool, options []*Option) *Chain {
	chain := Chain{
		Symbol:     strings.ToUpper(symbol),
		DateString: date,
		IsCall:     is_call,
		Options:    options,
	}
	chain.SortOptions()
	return &chain
}

func (chain *Chain) SortOptions() *Chain {
	options := chain.Options
	sort.Slice(options, func(i, j int) bool {
		a := options[i]
		b := options[j]
		if a == nil || b == nil {
			log.Println("How can this be? ", a, b)
		}
		return a.StrikePrice < b.StrikePrice
	})
	return chain
}
