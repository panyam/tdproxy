package gormdb

import (
	"github.com/panyam/goutils/utils"
	"log"
	// "gorm.io/driver/sqlite"
	"errors"
	"gorm.io/gorm"
	// "gorm.io/gorm/clause"
	"tdproxy/models"
	"time"
)

type ChainDB struct {
	optiondb *OptionDB
	db       *gorm.DB
}

func NewChainDB(db *gorm.DB) *ChainDB {
	// db.AutoMigrate(&models.ChainInfo{})
	db.AutoMigrate(&models.Chain{})
	return &ChainDB{
		optiondb: NewOptionDB(db),
		db:       db,
	}
}

func (db *ChainDB) GetChainInfo(symbol string) (out *models.ChainInfo, err error) {
	var chains []models.Chain
	result := db.db.Order("date_string").Where("symbol = ?", symbol).Find(&chains)
	if result.Error != nil {
		return nil, result.Error
	}
	// Now go through all and collect the dates
	lowestRefreshedAt := time.Now().UTC()
	out = &models.ChainInfo{Symbol: symbol}
	last := ""
	today := time.Now().UTC()
	for _, chain := range chains {
		// use the refresh time of the oldest date as the chain's last refreshed date
		if chain.LastRefreshedAt.Sub(lowestRefreshedAt) < 0 {
			if utils.ParseDate(chain.DateString).Sub(today) >= 0 {
				lowestRefreshedAt = chain.LastRefreshedAt
			}
		}
		if last != chain.DateString {
			out.AvailableDates = append(out.AvailableDates, chain.DateString)
		}
		last = chain.DateString
	}
	if len(out.AvailableDates) == 0 {
		out = nil
	} else {
		out.LastRefreshedAt = lowestRefreshedAt
	}
	return
}

func (db *ChainDB) SaveChainInfo(symbol string, last_refreshed_at time.Time) error {
	return nil
	//return db.db.Save(option).Error
}

func (db *ChainDB) GetChain(symbol string, date string, is_call bool) (*models.Chain, error) {
	var out models.Chain
	err := db.db.First(&out, "is_call = ? AND symbol = ? AND date_string = ?", is_call, symbol, date).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		} else {
			return nil, err
		}
	}

	// now load options
	out.Options, err = db.optiondb.GetOptions(is_call, symbol, date)
	return &out, err
}

func (db *ChainDB) SaveChain(chain *models.Chain) error {
	result := db.db.Model(chain).Updates(chain)
	err := result.Error
	if err == nil && result.RowsAffected == 0 {
		err = db.db.Create(chain).Error
	}
	err = db.optiondb.SaveOptions(chain.Options)
	if err != nil {
		log.Println("Error Saving Options: ", err)
	}
	return err
}
