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

func CreateTestTickerDB(t *testing.T) *TickerDB {
	dir, err := ioutil.TempDir("/tmp", "tdproxydb")
	if err != nil {
		log.Fatal(err)
	}

	dbroot := dir
	fmt.Println("DBRoot: ", dbroot)
	db := NewTickerDB(dbroot)
	_, err = os.Stat(db.TickersFolderPath)
	assert.Equal(t, err, nil, "Ticker folder path does not exist")
	return db
}

func TestNewTickerDB(t *testing.T) {
	db := CreateTestTickerDB(t)
	defer os.RemoveAll(db.DataRoot)
}

func TestTickerSaveAndGet(t *testing.T) {
	db := CreateTestTickerDB(t)
	defer os.RemoveAll(db.DataRoot)
	path, err := db.TickerPathForSymbol("SYM", false)
	assert.NotEqual(t, err, nil)
	assert.Equal(t, path, fmt.Sprintf("%s/tickers/SYM/QUOTES.json", db.DataRoot))

	// Here we will still fail as QUOTES.json doesnt exist
	path, err = db.TickerPathForSymbol("SYM", true)
	assert.NotEqual(t, err, nil)
	assert.Equal(t, path, fmt.Sprintf("%s/tickers/SYM/QUOTES.json", db.DataRoot))

	// Save a ticker here
	now := time.Now().UTC()
	info, err := utils.JsonDecodeStr(`{"a": 1, "b": 2}`)
	ticker := models.Ticker{
		Symbol:          "SYM",
		LastRefreshedAt: now,
		Info:            info.(map[string]interface{}),
	}
	db.SaveTicker(&ticker)

	loaded := db.GetTicker("SYM")
	assert.Equal(t, loaded, &ticker, "Saved ticker should be same due to caching")

	// Remove from cache and see what happens
	delete(db.tickerCache, "SYM")
	loaded = db.GetTicker("SYM")
	assert.Equal(t, loaded, &ticker, "Without caching Saved ticker should be same as?")
}
