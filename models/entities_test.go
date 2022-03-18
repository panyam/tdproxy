package models

import (
	"fmt"
	"github.com/panyam/goutils/utils"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestChainFromDict(t *testing.T) {
	optjson := `{
		"230": [{
        "putCall": "CALL",
        "bid": 390.0,
        "ask": 400.0,
        "last": 361.0,
        "mark": 396.67,
        "delta": 1.003,
        "gamma": 0.0,
        "openInterest": 7,
        "multiplier": 100.0
      }],
		"240": [{
        "putCall": "CALL",
        "bid": 380.0,
        "ask": 390.0,
        "last": 377.0,
        "mark": 387.0,
        "delta": 1.001,
        "openInterest": 258,
        "multiplier": 100.0
      }]
	}`

	options_by_price, _ := utils.JsonDecodeStr(optjson)
	now := time.Now()
	chain := ChainFromDict("TEST", "2022_01_02", true, options_by_price.(map[string]interface{}), now)
	assert.Equal(t, len(chain.Options), 2, "Expected 2 prices")
	fmt.Printf("%+v\n", chain.Options[0])
	assert.Equal(t, chain.Options[0].PriceString, "230", "Option price does not match")
	assert.Equal(t, chain.Options[0].AskPrice, 400.0, "Ask price doesnt match")
	assert.Equal(t, chain.Options[0].BidPrice, 390.0, "Bid price doesnt match")
	assert.Equal(t, chain.Options[0].MarkPrice, 396.67, "Mark price doesnt match")
	assert.Equal(t, chain.Options[0].Delta, 1.003, "Delta doesnt match")
	assert.Equal(t, chain.Options[0].Multiplier, 100.0, "Multiplier doesnt match")
	assert.Equal(t, chain.Options[0].OpenInterest, 7, "Open Interest doesnt match")

	assert.Equal(t, chain.Options[1].PriceString, "240", "Option price does not match")
	assert.Equal(t, chain.Options[1].AskPrice, 390.0, "Ask price doesnt match")
	assert.Equal(t, chain.Options[1].BidPrice, 380.0, "Bid price doesnt match")
	assert.Equal(t, chain.Options[1].MarkPrice, 387.0, "Mark price doesnt match")
	assert.Equal(t, chain.Options[1].Delta, 1.001, "Delta doesnt match")
	assert.Equal(t, chain.Options[1].Multiplier, 100.0, "Multiplier doesnt match")
	assert.Equal(t, chain.Options[1].OpenInterest, 258, "Open Interest doesnt match")
}
