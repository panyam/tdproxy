package filedb

import (
	"fmt"
	"github.com/panyam/goutils/utils"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"legfinder/tdproxy/models"
	"log"
	"os"
	"testing"
	"time"
)

// type OptionsByPrice =

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
	chain := models.ChainFromDict("TEST", "2022_01_02", true, options_by_price.(map[string]interface{}), now)
	assert.Equal(t, len(chain.Options), 2, "Expected 2 prices")
	fmt.Printf("%+v\n", chain.Options[0])
	assert.Equal(t, chain.Options[0].PriceString, "230", "Option price does not match")
	assert.Equal(t, chain.Options[0].AskPrice(), 400.0, "Ask price doesnt match")
	assert.Equal(t, chain.Options[0].BidPrice(), 390.0, "Bid price doesnt match")
	assert.Equal(t, chain.Options[0].Mark(), 396.67, "Mark price doesnt match")
	assert.Equal(t, chain.Options[0].Delta(), 1.003, "Delta doesnt match")
	assert.Equal(t, chain.Options[0].Multiplier(), 100.0, "Multiplier doesnt match")
	assert.Equal(t, chain.Options[0].OpenInterest(), 7, "Open Interest doesnt match")

	assert.Equal(t, chain.Options[1].PriceString, "240", "Option price does not match")
	assert.Equal(t, chain.Options[1].AskPrice(), 390.0, "Ask price doesnt match")
	assert.Equal(t, chain.Options[1].BidPrice(), 380.0, "Bid price doesnt match")
	assert.Equal(t, chain.Options[1].Mark(), 387.0, "Mark price doesnt match")
	assert.Equal(t, chain.Options[1].Delta(), 1.001, "Delta doesnt match")
	assert.Equal(t, chain.Options[1].Multiplier(), 100.0, "Multiplier doesnt match")
	assert.Equal(t, chain.Options[1].OpenInterest(), 258, "Open Interest doesnt match")
}

func CreateTestChainDB(t *testing.T) *ChainDB {
	dir, err := ioutil.TempDir("/tmp", "tdproxydb")
	if err != nil {
		log.Fatal(err)
	}

	dbroot := dir
	fmt.Println("DBRoot: ", dbroot)
	db := NewChainDB(dbroot)
	_, err = os.Stat(db.TickersFolderPath)
	assert.Equal(t, err, nil, "Ticker folder path does not exist")
	return db
}

func TestNewChainDB(t *testing.T) {
	db := CreateTestChainDB(t)
	defer os.RemoveAll(db.DataRoot)
}

func TestChainSaveAndGet(t *testing.T) {
	db := CreateTestChainDB(t)
	fmt.Println("DB Root: ", db.DataRoot)
	// defer os.RemoveAll(db.DataRoot)
	TEST_DATE := "2022_01_02"
	path, err := db.ChainPathForSymbol("SYM", TEST_DATE, false)
	assert.NotEqual(t, err, nil, fmt.Sprintf("Error should not be nil but found '%+v'", err))
	assert.Equal(t, path, fmt.Sprintf("%s/tickers/SYM/chains/%s", db.DataRoot, TEST_DATE))

	// Here we will still fail as QUOTES.json doesnt exist
	path, err = db.ChainPathForSymbol("SYM", "2022_01_02", true)
	assert.Equal(t, err, nil)
	assert.Equal(t, path, fmt.Sprintf("%s/tickers/SYM/chains/%s", db.DataRoot, TEST_DATE))

	// Save a ticker here
	now := time.Now().UTC()
	chain := models.NewChain(
		"SYM",
		TEST_DATE,
		true,
		MakeTestOptions("SYM", TEST_DATE, true, 10, 100, 10, 50, 10),
	)
	chain.LastRefreshedAt = now
	db.SaveChain(chain)

	loaded := db.GetChain("SYM", TEST_DATE, true)
	assert.Equal(t, loaded, chain, "Saved chain should be same due to caching")

	// Remove from cache and see what happens
	chainKey := db.ChainKeyFor("SYM", TEST_DATE, true)
	delete(db.chainCache, chainKey)
	loaded = db.GetChain("SYM", TEST_DATE, true)
	assert.Equal(t, loaded, chain, "Without caching Saved ticker should be same as?")
}
