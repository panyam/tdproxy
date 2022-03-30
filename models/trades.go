package models

import (
	"github.com/panyam/goutils/utils"
	// "log"
)

type OrderBase struct {
	Buy               bool
	Quantity          int32
	OverriddenPremium float64
}

type Single struct {
	OrderBase
	OptionKey string
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
	StrategyName       string
	PayoffExpectedGain float64 `gorm:"index:ByDateAndGain,priority:2"`
	PayoffExpectedLoss float64
	PayoffGainProb     float64 `gorm:"index:ByDateAndGainProb,priority:2"`
	PayoffBPE          float64
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

func (o *Single) ToJson() utils.StringMap {
	out := o.OrderBase.ToJson()
	out["option"] = o.OptionKey
	return out
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

func (o *Order) ToJson() utils.StringMap {
	if o.IsSingle() {
		return o.Single().ToJson()
	} else {
		return o.Spread().ToJson()
	}
}

func (t *Trade) ToJson() utils.StringMap {
	out := make(utils.StringMap)
	out["id"] = t.TradeId
	out["sym"] = t.Symbol
	out["date"] = t.Date
	out["strategy"] = t.StrategyName
	out["metadata"] = t.Metadata
	out["payoff"] = utils.StringMap{
		"expected_gain": t.PayoffExpectedGain,
		"expected_loss": t.PayoffExpectedLoss,
		"gain_prob":     t.PayoffGainProb,
		"bpe":           t.PayoffBPE,
	}
	var orders []utils.StringMap
	for _, s := range t.Orders {
		orders = append(orders, s.ToJson())
	}
	out["orders"] = orders
	return out
}
