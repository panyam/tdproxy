package models

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"
)

type AuthToken struct {
	AccessToken   string
	Scope         string
	ExpiresIn     int32
	TokenType     string
	CreatedAt     string
	LastFetchedAt string
}

type Ticker struct {
	Symbol          string
	LastRefreshedAt time.Time
	Info            map[string]interface{}
}

type Option struct {
	Symbol      string
	DateString  string
	PriceString string
	IsCall      bool
	Info        map[string]interface{}
}

type TickerChainInfo struct {
	Symbol          string
	AvailableDates  []string
	LastRefreshedAt time.Time
}

type Chain struct {
	Symbol          string
	DateString      string
	IsCall          bool
	LastRefreshedAt time.Time
	Options         []*Option
}

func (opt *Option) StrikePrice() float64 {
	result, err := strconv.ParseFloat(opt.PriceString, 64)
	if err != nil {
		fmt.Printf("Invalid price string: %s\n", opt.PriceString)
	}
	return result
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
