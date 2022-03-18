package models

import (
	"fmt"
	"github.com/panyam/goutils/utils"
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
	authToken      utils.StringMap
	userPrincipals utils.StringMap
}

type AuthToken struct {
	AccessToken string
	Scope       string
	TokenType   string
}

type Ticker struct {
	Symbol          string `gorm:"primaryKey"`
	LastRefreshedAt time.Time
	info            map[string]interface{} `gorm:"-"`
}

func NewTicker(symbol string, refreshed_at time.Time, info map[string]interface{}) *Ticker {
	return &Ticker{
		Symbol:          symbol,
		LastRefreshedAt: refreshed_at,
		info:            info,
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
	OpenInterest int
	Delta        float64
	Multiplier   float64
	info         map[string]interface{} `gorm:"-"`
}

func NewOption(symbol string, date_string string, price_string string, is_call bool, info map[string]interface{}) *Option {
	out := &Option{
		Symbol:      symbol,
		DateString:  date_string,
		PriceString: price_string,
		IsCall:      is_call,
	}
	return out.Refresh()
}

func (opt *Option) Refresh() *Option {
	if val, ok := opt.info["ask"]; ok {
		opt.AskPrice = val.(float64)
	}
	if val, ok := opt.info["bid"]; ok {
		opt.BidPrice = val.(float64)
	}
	if val, ok := opt.info["mark"]; ok {
		opt.MarkPrice = val.(float64)
	} else {
		opt.MarkPrice = (opt.AskPrice + opt.BidPrice) / 2
	}

	if val, ok := opt.info["openInterest"]; ok {
		opt.OpenInterest = int(val.(float64))
	}
	if val, ok := opt.info["delta"]; ok {
		opt.Delta = val.(float64)
	}
	if val, ok := opt.info["multiplier"]; ok {
		opt.Multiplier = val.(float64)
	}
	if opt.StrikePrice <= 0 {
		result, err := strconv.ParseFloat(opt.PriceString, 64)
		if err != nil {
			fmt.Printf("Invalid price string: %s\n", opt.PriceString)
		}
		opt.StrikePrice = result
	}
	return opt
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
			info:        detail,
		}
		options = append(options, option.Refresh())
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
		return options[i].StrikePrice < options[j].StrikePrice
	})
	return chain
}
