package filedb

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"log"
	"os"
	dbp "tdproxy/db"
	"tdproxy/models"
	"testing"
	"time"
)

func CreateTestChainDB(t *testing.T) *ChainDB {
	dir, err := ioutil.TempDir("/tmp", "tdproxydb")
	if err != nil {
		log.Panic("Could not create tempdir: ", err)
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
	assert.Equal(t, err, nil, "Error should have been nil")
	assert.Equal(t, path, fmt.Sprintf("%s/tickers/SYM/chains/%s", db.DataRoot, TEST_DATE))

	// Save a ticker here
	now := time.Now().UTC()
	chain := models.NewChain(
		"SYM",
		TEST_DATE,
		true,
		dbp.MakeTestOptions("SYM", TEST_DATE, true, 10, 100, 10, 50, 10),
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
