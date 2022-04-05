package services

import (
	structpb "google.golang.org/protobuf/types/known/structpb"
	"log"
	"tdproxy/models"
	"tdproxy/protos"
)

func OptionToProto(option *models.Option) (out *protos.Option, err error) {
	out = &protos.Option{
		Symbol:       option.Symbol,
		DateString:   option.DateString,
		PriceString:  option.PriceString,
		StrikePrice:  option.StrikePrice,
		IsCall:       option.IsCall,
		AskPrice:     option.AskPrice,
		BidPrice:     option.BidPrice,
		MarkPrice:    option.MarkPrice,
		Multiplier:   option.Multiplier,
		Delta:        option.Delta,
		OpenInterest: option.OpenInterest,
	}
	val := option.Info
	if val != nil {
		out.Info, err = structpb.NewStruct(val)
	}
	if err != nil {
		log.Println("Error parsing info: ", err)
	}
	return out, nil
}

func TradeFromProto(trade *protos.Trade) (out *models.Trade) {
	out = &models.Trade{
		TradeId:            trade.TradeId,
		Symbol:             trade.Symbol,
		Date:               trade.Date,
		PayoffExpectedGain: trade.PayoffExpectedGain,
		PayoffExpectedLoss: trade.PayoffExpectedLoss,
		PayoffGainProb:     trade.PayoffGainProb,
		PayoffBPE:          trade.PayoffBpe,
		Orders:             OrdersFromProto(trade.Orders),
	}
	return
}

func TradeToProto(trade *models.Trade) (out *protos.Trade) {
	out = &protos.Trade{
		TradeId:            trade.TradeId,
		Symbol:             trade.Symbol,
		Date:               trade.Date,
		PayoffExpectedGain: trade.PayoffExpectedGain,
		PayoffExpectedLoss: trade.PayoffExpectedLoss,
		PayoffGainProb:     trade.PayoffGainProb,
		PayoffBpe:          trade.PayoffBPE,
		Orders:             OrdersToProto(trade.Orders),
	}
	return
}

func OrderBaseFromProto(item *protos.OrderBase) *models.OrderBase {
	return &models.OrderBase{
		Buy:               item.Buy,
		Quantity:          item.Quantity,
		OverriddenPremium: item.OverriddenPremium,
	}
}

func OrderBaseToProto(item *models.OrderBase) *protos.OrderBase {
	return &protos.OrderBase{
		Buy:               item.Buy,
		Quantity:          item.Quantity,
		OverriddenPremium: item.OverriddenPremium,
	}
}

func SingleFromProto(item *protos.Single) (out *models.Single) {
	out = &models.Single{
		OrderBase: *OrderBaseFromProto(item.OrderBase),
		// OptionKey: item.OptionKey,
		Symbol:      item.Symbol,
		IsCall:      item.IsCall,
		PriceString: item.PriceString,
		DateString:  item.DateString,
	}
	return
}

func SingleToProto(item *models.Single) (out *protos.Single) {
	out = &protos.Single{
		OrderBase: OrderBaseToProto(&item.OrderBase),
		// OptionKey: item.OptionKey,
		Symbol:      item.Symbol,
		IsCall:      item.IsCall,
		PriceString: item.PriceString,
		DateString:  item.DateString,
	}
	return
}

func SpreadFromProto(item *protos.Spread) (out *models.Spread) {
	out = &models.Spread{
		OrderBase: *OrderBaseFromProto(item.OrderBase),
		Name:      item.Name,
		Singles:   SinglesFromProto(item.Singles),
	}
	return
}

func SpreadToProto(item *models.Spread) (out *protos.Spread) {
	out = &protos.Spread{
		OrderBase: OrderBaseToProto(&item.OrderBase),
		Name:      item.Name,
		Singles:   SinglesToProto(item.Singles),
	}
	return
}

func OrderFromProto(item *protos.Order) (out *models.Order) {
	switch item.Details.(type) {
	case *protos.Order_Single:
		return models.SingleOrder(SingleFromProto(item.GetSingle()))
	case *protos.Order_Spread:
		return models.SpreadOrder(SpreadFromProto(item.GetSpread()))
	}
	return
}

func OrderToProto(item *models.Order) (out *protos.Order) {
	if item.IsSingle() {
		return &protos.Order{
			Details: &protos.Order_Single{Single: SingleToProto(item.Single())},
		}
	}
	if item.IsSpread() {
		return &protos.Order{
			Details: &protos.Order_Spread{Spread: SpreadToProto(item.Spread())},
		}
	}
	return
}

// Convert lists into and from protos
func SinglesFromProto(items []*protos.Single) (out []*models.Single) {
	for _, item := range items {
		out = append(out, SingleFromProto(item))
	}
	return
}

func SinglesToProto(items []*models.Single) (out []*protos.Single) {
	for _, item := range items {
		out = append(out, SingleToProto(item))
	}
	return
}

func SpreadsFromProto(items []*protos.Spread) (out []*models.Spread) {
	for _, item := range items {
		out = append(out, SpreadFromProto(item))
	}
	return
}

func SpreadsToProto(items []*models.Spread) (out []*protos.Spread) {
	for _, item := range items {
		out = append(out, SpreadToProto(item))
	}
	return
}

func TradesFromProto(items []*protos.Trade) (out []*models.Trade) {
	for _, item := range items {
		out = append(out, TradeFromProto(item))
	}
	return
}

func TradesToProto(items []*models.Trade) (out []*protos.Trade) {
	for _, item := range items {
		out = append(out, TradeToProto(item))
	}
	return
}

func OrdersFromProto(items []*protos.Order) (out []*models.Order) {
	for _, item := range items {
		out = append(out, OrderFromProto(item))
	}
	return
}

func OrdersToProto(items []*models.Order) (out []*protos.Order) {
	for _, item := range items {
		out = append(out, OrderToProto(item))
	}
	return
}

func ProbDistFromProto(item *protos.ProbDist) (out *models.ProbDist) {
	out = &models.ProbDist{
		Distribution: ProbRangesFromProto(item.Distribution),
	}
	return
}

func ProbDistToProto(item *models.ProbDist) (out *protos.ProbDist) {
	out = &protos.ProbDist{
		Distribution: ProbRangesToProto(item.Distribution),
	}
	return
}

func ProbRangeFromProto(item *protos.ProbRange) (out *models.ProbRange) {
	out = &models.ProbRange{
		X1:        item.X1,
		X2:        item.X2,
		LtProb:    item.LtProb,
		RangeProb: item.RangeProb,
	}
	return
}

func ProbRangeToProto(item *models.ProbRange) (out *protos.ProbRange) {
	out = &protos.ProbRange{
		X1:        item.X1,
		X2:        item.X2,
		LtProb:    item.LtProb,
		RangeProb: item.RangeProb,
	}
	return
}

func ProbRangesFromProto(items []*protos.ProbRange) (out []*models.ProbRange) {
	for _, item := range items {
		out = append(out, ProbRangeFromProto(item))
	}
	return
}

func ProbRangesToProto(items []*models.ProbRange) (out []*protos.ProbRange) {
	for _, item := range items {
		out = append(out, ProbRangeToProto(item))
	}
	return
}
