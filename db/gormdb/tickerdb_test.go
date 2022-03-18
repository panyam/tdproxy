package gormdb

import (
	"fmt"
	"github.com/panyam/goutils/utils"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"io/ioutil"
	"log"
	"os"
	"path"
	"tdproxy/models"
	"testing"
	"time"
)

func CreateTestTickerDB(t *testing.T) (*TickerDB, string) {
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
	return NewTickerDB(db), dir
}

func TestNewTickerDB(t *testing.T) {
	_, dbroot := CreateTestTickerDB(t)
	defer os.RemoveAll(dbroot)
}

func TestTickerSaveAndGet(t *testing.T) {
	db, dbroot := CreateTestTickerDB(t)
	defer os.RemoveAll(dbroot)

	// Save a ticker here
	now := time.Now().UTC()
	info, err := utils.JsonDecodeStr(`{"a": 1, "b": 2}`)
	ticker := models.Ticker{
		Symbol:          "SYM",
		LastRefreshedAt: now,
		Info:            info.(map[string]interface{}),
	}
	db.SaveTicker(&ticker)

	loaded, err := db.GetTicker("SYM")
	assert.Equal(t, err, nil, "Should be able to load ticker")
	assert.Equal(t, loaded, &ticker, "Saved ticker should be same due to caching")
}
