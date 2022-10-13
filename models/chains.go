package models

import (
	"fmt"
	"log"
	"sort"
	"strings"
	"time"
)

type ChainInfo struct {
	Symbol          string   `gorm:"primaryKey"`
	AvailableDates  []string `gorm:"type:integer[]"`
	LastRefreshedAt time.Time
}

type Chain struct {
	Symbol          string `gorm:"primaryKey"`
	DateString      string `gorm:"primaryKey"`
	IsCall          bool   `gorm:"primaryKey" gorm:"type:boolean; column:is_call"`
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
		option := NewOption(symbol, date, price_string, is_call, detail)
		options = append(options, option)
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

func (chain *Chain) ShortKey() string {
	ot := "P"
	if chain.IsCall {
		ot = "C"
	}
	return fmt.Sprintf("%s%s%s", ot, chain.Symbol, chain.DateString)
}
