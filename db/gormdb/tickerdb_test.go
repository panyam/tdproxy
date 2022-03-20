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
	dbtu "tdproxy/db"
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
	ticker := models.NewTicker(
		"SYM",
		now,
		info.(map[string]interface{}),
	)
	err = db.SaveTicker(ticker)
	assert.Equal(t, err, nil, "SaveTicker Error should be nil")

	loaded, err := db.GetTicker("SYM")
	assert.Equal(t, err, nil, "Should be able to load ticker")
	dbtu.AssertTickersEqual(t, loaded, ticker)
}
