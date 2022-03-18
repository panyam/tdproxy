package filedb

import (
	"fmt"
	"github.com/panyam/goutils/utils"
	"legfinder/tdproxy/models"
)

func MakeTestOptions(symbol string, date string, is_call bool, start_price float64, end_price float64, price_incr float64, ask float64, open_interest int) []*models.Option {
	var out []*models.Option
	for curr := start_price; curr <= end_price; curr += price_incr {
		infostr := fmt.Sprintf(`{"ask": %f, "bid": %f, "openInterest": %d, "multiplier": 100, "delta": 0.5}`, ask, ask, open_interest)
		info, _ := utils.JsonDecodeStr(infostr)
		option := models.Option{
			Symbol:      symbol,
			DateString:  date,
			PriceString: utils.PriceString(curr),
			IsCall:      is_call,
			Info:        info.(map[string]interface{}),
		}
		out = append(out, &option)
	}
	return out
}
