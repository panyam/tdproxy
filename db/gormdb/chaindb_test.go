package gormdb

import (
	"fmt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"path"
	// "github.com/panyam/goutils/utils"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"log"
	"os"
	dbp "tdproxy/db"
	"tdproxy/models"
	"testing"
	"time"
)

func CreateTestChainDB(t *testing.T) (*ChainDB, string) {
	dir, err := ioutil.TempDir("/tmp", "tdproxydb")
	if err != nil {
		log.Fatal(err)
	}

	filepath := path.Join(dir, "test.db")
	db, err := gorm.Open(sqlite.Open(filepath), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	fmt.Println("DBRoot: ", filepath)
	return NewChainDB(db), dir
}

func TestNewChainDB(t *testing.T) {
	_, dbroot := CreateTestChainDB(t)
	defer os.RemoveAll(dbroot)
}

func TestChainSaveAndGet(t *testing.T) {
	db, dbroot := CreateTestChainDB(t)
	defer os.RemoveAll(dbroot)
	TEST_DATE := "2022_01_02"

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

	loaded, err := db.GetChain("SYM", TEST_DATE, true)
	assert.Equal(t, err, nil, "Error should be nil")
	assert.Equal(t, loaded, chain, "Saved chain should be same due to caching")
}
