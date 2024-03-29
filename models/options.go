package models

import (
	"fmt"
	"strconv"
)

/*
type OptionJsonField struct {
	*Json
	OptionSymbol      string `gorm:"primaryKey"`
	OptionDateString  string `gorm:"primaryKey"`
	OptionIsCall      bool   `gorm:"primaryKey"`
	OptionPriceString string `gorm:"primaryKey"`
}
*/

type Option struct {
	Symbol       string `gorm:"primaryKey" gorm:"index:ByCallSymbolDate,priority:1"`
	DateString   string `gorm:"primaryKey" gorm:"index:ByCallSymbolDate,priority:2"`
	IsCall       bool   `gorm:"primaryKey" gorm:"index:ByCallSymbolDate,priority:3" gorm:"type:boolean; column:is_call"`
	PriceString  string `gorm:"primaryKey"`
	StrikePrice  float64
	AskPrice     float64
	BidPrice     float64
	MarkPrice    float64
	OpenInterest int32
	Delta        float64
	Multiplier   float64
	Info         JsonField `gorm:"type:text"` // OptionJsonJson
}

func (opt *Option) HasValidDelta() bool {
	// For calls - for puts it should be -ve or 0
	return opt.Delta >= 0
}

func NewOption(symbol string, date_string string, price_string string, is_call bool, info map[string]interface{}) *Option {
	// log.Println("Args: ", symbol, date_string, price_string, is_call)
	out := &Option{
		Symbol:      symbol,
		DateString:  date_string,
		PriceString: price_string,
		IsCall:      is_call,
		Info:        info,
		/*
			Info: OptionJsonField{
				Json:              NewJson(info),
				OptionSymbol:      symbol,
				OptionDateString:  date_string,
				OptionIsCall:      is_call,
				OptionPriceString: price_string,
			},
		*/
	}
	out.Refresh()
	return out
}

func (opt *Option) ShortKey() string {
	ot := "P"
	if opt.IsCall {
		ot = "C"
	}
	return fmt.Sprintf("%s%s%s%s", ot, opt.Symbol, opt.DateString, opt.PriceString)
}

func (opt *Option) Refresh() bool {
	/*
		res, err := opt.Info.Value()
		if err != nil {
			log.Println("Err refreshing option info: ", err)
			return false
		}
	*/
	info := opt.Info // .(utils.StringMap)
	if val, ok := info["ask"]; ok {
		opt.AskPrice = val.(float64)
	}
	if val, ok := info["bid"]; ok {
		opt.BidPrice = val.(float64)
	}
	if val, ok := info["mark"]; ok {
		opt.MarkPrice = val.(float64)
	} else {
		opt.MarkPrice = (opt.AskPrice + opt.BidPrice) / 2
	}

	if val, ok := info["openInterest"]; ok {
		opt.OpenInterest = int32(val.(float64))
	}
	if val, ok := info["delta"]; ok {
		if opt.Delta, ok = val.(float64); !ok {
			opt.Delta = -1
		}
	}
	if val, ok := info["multiplier"]; ok {
		if opt.Multiplier, ok = val.(float64); !ok {
			opt.Multiplier = -1
		}
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
