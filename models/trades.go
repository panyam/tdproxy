package models

import (
	"fmt"
	"github.com/panyam/goutils/utils"
	"log"
)

type OrderBase struct {
	Buy               bool
	Quantity          int32
	OverriddenPremium float64
}

type Single struct {
	OrderBase
	// OptionKey   string
	Symbol      string
	DateString  string
	PriceString string
	IsCall      bool
}

type Spread struct {
	OrderBase
	Name    string
	Singles []*Single
}

type Order struct {
	spread *Spread
	single *Single
}

func SpreadOrder(spread *Spread) *Order {
	return &Order{spread: spread}
}

func SingleOrder(single *Single) *Order {
	return &Order{single: single}
}

func (order *Order) IsSingle() bool {
	return order.single != nil
}

func (order *Order) IsSpread() bool {
	return order.spread != nil
}

func (order *Order) Single() *Single {
	return order.single
}

func (order *Order) Spread() *Spread {
	return order.spread
}

type Trade struct {
	TradeId            string `gorm:"primaryKey"`
	Symbol             string `gorm:"index:BySymbol,priority:1"`
	Date               string `gorm:"index:ByDate,priority:1" gorm:"index:ByDateAndGain,priority:1" gorm:"index:ByDateAndGainProb,priority:1"`
	LoadDate           string
	Strategy           string
	PayoffExpectedGain float64 `gorm:"index:ByDateAndGain,priority:2"`
	PayoffExpectedLoss float64
	PayoffGainProb     float64 `gorm:"index:ByDateAndGainProb,priority:2"`
	PayoffBPE          float64
	PayoffMaxProfit    float64
	Orders             []*Order        `gorm:"-"`
	Metadata           utils.StringMap `gorm:"-"`
}

func (o *OrderBase) ToJson() utils.StringMap {
	out := utils.StringMap{
		"buy": o.Buy,
	}
	if o.Quantity > 1 {
		out["quantity"] = o.Quantity
	}
	if o.OverriddenPremium > 0 {
		out["premium"] = o.OverriddenPremium
	}
	return out
}

func (o *OrderBase) FromJson(input utils.StringMap) error {
	if value, exists := input["buy"]; exists {
		o.Buy = value.(bool)
	} else {
		o.Buy = false
	}

	if value, exists := input["quantity"]; exists {
		if intval, ok := value.(int32); ok {
			o.Quantity = intval
		} else if floatval, ok := value.(int32); ok {
			o.Quantity = int32(floatval)
		}
	} else {
		o.Quantity = 1
	}

	if value, exists := input["premium"]; exists {
		o.OverriddenPremium = value.(float64)
	} else {
		o.OverriddenPremium = 1
	}
	return nil
}

func (o *Single) ToJson() utils.StringMap {
	out := o.OrderBase.ToJson()
	// out["option"] = o.OptionKey
	out["symbol"] = o.Symbol
	out["date_string"] = o.DateString
	out["price_string"] = o.PriceString
	out["is_call"] = o.IsCall
	return out
}

func (o *Single) FromJson(input utils.StringMap) error {
	if err := o.OrderBase.FromJson(input); err != nil {
		return err
	}

	if value, exists := input["symbol"]; exists {
		o.Symbol = value.(string)
	} else {
		return fmt.Errorf("Symbol missing in Single")
	}
	if value, exists := input["date_string"]; exists {
		o.DateString = value.(string)
	} else {
		return fmt.Errorf("DateString missing in Single")
	}

	if value, exists := input["price_string"]; exists {
		o.PriceString = value.(string)
	} else {
		return fmt.Errorf("PriceString missing in Single")
	}

	if value, exists := input["is_call"]; exists {
		o.IsCall = value.(bool)
	} else {
		return fmt.Errorf("IsCall missing in Single")
	}
	return nil
}

func (o *Spread) ToJson() utils.StringMap {
	out := o.OrderBase.ToJson()
	out["name"] = o.Name
	var singles []utils.StringMap
	for _, s := range o.Singles {
		singles = append(singles, s.ToJson())
	}
	out["singles"] = singles
	return out
}

func (o *Spread) FromJson(input utils.StringMap) error {
	if err := o.OrderBase.FromJson(input); err != nil {
		return err
	}

	if value, exists := input["name"]; exists {
		o.Name = value.(string)
	} else {
		return fmt.Errorf("Name missing in Spread")
	}
	o.Singles = nil
	singles := input["singles"].([]interface{})
	for _, s := range singles {
		single := &Single{}
		if err := single.FromJson(s.(utils.StringMap)); err != nil {
			return err
		}
		o.Singles = append(o.Singles, single)
	}
	return nil
}

func (o *Order) ToJson() utils.StringMap {
	out := make(utils.StringMap)
	if o.IsSingle() {
		out["single"] = o.Single().ToJson()
	} else {
		out["spread"] = o.Spread().ToJson()
	}
	return out
}

func (o *Order) FromJson(input utils.StringMap) error {
	o.single = nil
	o.spread = nil
	if singlejson, exists := input["single"]; exists {
		o.single = &Single{}
		if err := o.single.FromJson(singlejson.(utils.StringMap)); err != nil {
			return err
		}
	}
	if spreadjson, exists := input["spread"]; exists {
		o.spread = &Spread{}
		if err := o.spread.FromJson(spreadjson.(utils.StringMap)); err != nil {
			return err
		}
	}
	return nil
}

func (t *Trade) ToJson() utils.StringMap {
	out := make(utils.StringMap)
	out["id"] = t.TradeId
	out["sym"] = t.Symbol
	out["date"] = t.Date
	out["strategy"] = t.Strategy
	out["metadata"] = t.Metadata
	out["payoff"] = utils.StringMap{
		"expected_gain": t.PayoffExpectedGain,
		"expected_loss": t.PayoffExpectedLoss,
		"gain_prob":     t.PayoffGainProb,
		"bpe":           t.PayoffBPE,
		"max_profit":    t.PayoffMaxProfit,
	}
	var orders []utils.StringMap
	for _, s := range t.Orders {
		orders = append(orders, s.ToJson())
	}
	out["orders"] = orders
	return out
}

func (o *Trade) FromJson(input utils.StringMap) error {
	if value, exists := input["id"]; exists {
		o.TradeId = value.(string)
	} else {
		return fmt.Errorf("TradeId does not exist in input")
	}

	if value, exists := input["sym"]; exists {
		o.Symbol = value.(string)
	} else {
		return fmt.Errorf("Symbol does not exist in input")
	}

	if value, exists := input["date"]; exists {
		o.Date = value.(string)
	} else {
		return fmt.Errorf("Date does not exist in input")
	}

	if value, exists := input["strategy"]; exists {
		o.Strategy = value.(string)
	} else {
		o.Strategy = ""
	}

	if value, exists := input["metadata"]; exists && value != nil {
		o.Metadata = value.(utils.StringMap)
	} else {
		o.Metadata = nil
	}

	if value, exists := input["payoff"]; exists {
		payoff := value.(utils.StringMap)
		if value, exists := payoff["expected_gain"]; exists {
			o.PayoffExpectedGain = value.(float64)
		} else {
			return fmt.Errorf("Missing payoff expected gain")
		}

		if value, exists := payoff["expected_loss"]; exists {
			o.PayoffExpectedLoss = value.(float64)
		} else {
			return fmt.Errorf("Missing payoff expected loss")
		}

		if value, exists := payoff["gain_prob"]; exists {
			o.PayoffGainProb = value.(float64)
		} else {
			return fmt.Errorf("Missing payoff expected gain_prob")
		}

		if value, exists := payoff["bpe"]; exists {
			o.PayoffBPE = value.(float64)
		} else {
			return fmt.Errorf("Missing payoff expected bpe")
		}

		if value, exists := payoff["max_profit"]; exists {
			log.Println("Max Profit: ", value)
			o.PayoffMaxProfit = value.(float64)
		} else {
			return fmt.Errorf("Missing payoff expected bpe")
		}
	} else {
		return fmt.Errorf("Payoff not exist in input")
	}

	o.Orders = nil
	orders := input["orders"].([]interface{})
	for _, s := range orders {
		order := &Order{}
		if err := order.FromJson(s.(utils.StringMap)); err != nil {
			return err
		}
		o.Orders = append(o.Orders, order)
	}

	return nil
}
