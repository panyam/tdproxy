package db

import (
	"fmt"
	"github.com/panyam/goutils/utils"
	"github.com/stretchr/testify/assert"
	"tdproxy/models"
	"testing"
)

func MakeTestOptions(symbol string, date string, is_call bool, start_price float64, end_price float64, price_incr float64, ask float64, open_interest int) []*models.Option {
	var out []*models.Option
	for curr := start_price; curr <= end_price; curr += price_incr {
		infostr := fmt.Sprintf(`{"ask": %f, "bid": %f, "openInterest": %d, "multiplier": 100, "delta": 0.5}`, ask, ask, open_interest)
		info, _ := utils.JsonDecodeStr(infostr)
		option := models.NewOption(
			symbol,
			date,
			utils.PriceString(curr),
			is_call,
			info.(map[string]interface{}),
		)
		out = append(out, option)
	}
	return out
}

func AssertTickersEqual(t *testing.T, t1 *models.Ticker, t2 *models.Ticker) {
	assert.Equal(t, t1.Symbol, t2.Symbol, "Symbols should be equal")
	assert.Equal(t, t1.LastRefreshedAt, t2.LastRefreshedAt, "RefreshedTime should be equal")
}

func AssertChainsEqual(t *testing.T, t1 *models.Chain, t2 *models.Chain) {
	assert.Equal(t, t1.Symbol, t2.Symbol, "Symbols should be equal")
	assert.Equal(t, t1.DateString, t2.DateString, "DateStrings should be equal")
	assert.Equal(t, t1.IsCall, t2.IsCall, "IsCalls should be equal")
	assert.Equal(t, t1.LastRefreshedAt, t2.LastRefreshedAt, "RefreshedTime should be equal")
}
