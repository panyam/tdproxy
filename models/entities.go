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
	Info            map[string]interface{}
}

type Option struct {
	Symbol      string `gorm:"primaryKey" gorm:"index:ByCallSymbolDate,priority:1"`
	DateString  string `gorm:"primaryKey" gorm:"index:ByCallSymbolDate,priority:2"`
	IsCall      bool   `gorm:"primaryKey" gorm:"index:ByCallSymbolDate,priority:3"`
	PriceString string `gorm:"primaryKey"`
	Info        map[string]interface{}
	strikePrice float64
}

func (opt *Option) StrikePrice() float64 {
	if opt.strikePrice <= 0 {
		result, err := strconv.ParseFloat(opt.PriceString, 64)
		if err != nil {
			fmt.Printf("Invalid price string: %s\n", opt.PriceString)
		}
		opt.strikePrice = result
	}
	return opt.strikePrice
}

func (opt *Option) AskPrice() float64 {
	return opt.Info["ask"].(float64)
}

func (opt *Option) Mark() float64 {
	val, ok := opt.Info["mark"]
	if ok {
		return val.(float64)
	}
	return (opt.AskPrice() + opt.BidPrice()) / 2
}

func (opt *Option) BidPrice() float64 {
	return opt.Info["bid"].(float64)
}

func (opt *Option) OpenInterest() int {
	return int(opt.Info["openInterest"].(float64))
}

func (opt *Option) Delta() float64 {
	return opt.Info["delta"].(float64)
}

func (opt *Option) Multiplier() float64 {
	return opt.Info["multiplier"].(float64)
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
		options = append(options, &option)
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
		return options[i].StrikePrice() < options[j].StrikePrice()
	})
	return chain
}
