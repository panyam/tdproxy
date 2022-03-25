package gormdb

import (
	"log"
	// "gorm.io/driver/sqlite"
	"errors"
	"gorm.io/gorm"
	// "gorm.io/gorm/clause"
	"tdproxy/models"
	"time"
)

type AuthDB struct {
	db *gorm.DB
}

func NewAuthDB(db *gorm.DB) *AuthDB {
	db.AutoMigrate(&models.Auth{})
	db.AutoMigrate(&models.AuthTokenJsonField{})
	db.AutoMigrate(&models.UserPrincipalsJsonField{})
	return &AuthDB{
		db: db,
	}
}

func (a *AuthDB) LastAuth() *models.Auth {
	var out models.Auth
	err := a.db.Preload("UserPrincipals").Preload("AuthToken").First(&out).Error
	if err != nil {
		return nil
	}
	return &out
}

func (authdb *AuthDB) EnsureAuth(client_id string) (auth *models.Auth, err error) {
	auth, err = authdb.GetAuth(client_id)
	log.Println("Err, out: ", err, auth)
	if err == nil && auth == nil {
		// Does not exist so create
		auth = &models.Auth{ClientId: client_id}
		err = authdb.SaveAuth(auth)
	}
	return
}

func (db *AuthDB) GetAuth(client_id string) (*models.Auth, error) {
	var out models.Auth
	err := db.db.Preload("UserPrincipals").Preload("AuthToken").First(&out, "client_id = ?", client_id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		} else {
			return nil, err
		}
	}
	out.AuthToken.Value()
	out.UserPrincipals.Value()
	return &out, err
}

func (authdb *AuthDB) SaveAuth(auth *models.Auth) (err error) {
	result := authdb.db.Model(auth).Updates(auth)
	err = result.Error
	if err == nil && result.RowsAffected == 0 {
		err = authdb.db.Create(auth).Error
	}
	return
}

type OptionDB struct {
	db *gorm.DB
}

func NewOptionDB(db *gorm.DB) *OptionDB {
	db.AutoMigrate(&models.Option{})
	db.AutoMigrate(&models.OptionJsonField{})
	return &OptionDB{db: db}
}

func (db *OptionDB) GetOption(symbol string, date string, is_call bool, price string) (*models.Option, error) {
	var option models.Option
	err := db.db.First(&option, "symbol = ? AND date_string = ? AND is_call = ? AND price_string = ?", symbol, date, is_call, price).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		} else {
			return nil, err
		}
	}
	return &option, err
}

/**
 * Get options for a particular symbol on a given date for either calls or puts.
 */
func (db *OptionDB) GetOptions(is_call bool, symbol string, date string) ([]*models.Option, error) {
	var options []*models.Option
	err := db.db.Where("symbol = ? AND date_string = ? AND is_call = ?", symbol, date, is_call).Find(&options).Error
	if err != nil {
		return nil, err
	}
	return options, err
}

/**
 * Save a particular option.
 */
func (db *OptionDB) SaveOption(option *models.Option) (err error) {
	result := db.db.Model(option).Updates(option)
	err = result.Error
	if err == nil && result.RowsAffected == 0 {
		err = db.db.Create(option).Error
	}
	return
}

/**
 * Delete options on a particular date for a given symbol
 */
func (db *OptionDB) DeleteOptions(symbol string, date string, is_call bool) error {
	return db.db.Where("symbol = ? AND date_string = ? AND is_call = ?", symbol, date, is_call).Delete(&models.Option{}).Error
}

func (db *OptionDB) SaveOptions(options []*models.Option) error {
	return db.db.Transaction(func(tx *gorm.DB) error {
		for _, option := range options {
			if err := db.SaveOption(option); err != nil {
				return err
			}
		}
		return nil
	})
}

type ChainDB struct {
	optiondb *OptionDB
	db       *gorm.DB
}

func NewChainDB(db *gorm.DB) *ChainDB {
	db.AutoMigrate(&models.ChainInfo{})
	db.AutoMigrate(&models.Chain{})
	return &ChainDB{
		optiondb: NewOptionDB(db),
		db:       db,
	}
}

func (db *ChainDB) GetChainInfo(symbol string) (*models.ChainInfo, error) {
	var out models.ChainInfo
	err := db.db.First(&out, "symbol = ?", symbol).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		} else {
			return nil, err
		}
	}
	return &out, err
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

type TickerDB struct {
	db *gorm.DB
}

func NewTickerDB(db *gorm.DB) *TickerDB {
	db.AutoMigrate(&models.Ticker{})
	db.AutoMigrate(&models.TickerJsonField{})
	return &TickerDB{db: db}
}

func (db *TickerDB) GetTicker(symbol string) (*models.Ticker, error) {
	var out models.Ticker
	err := db.db.First(&out, "symbol = ?", symbol).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		} else {
			return nil, err
		}
	}
	return &out, err
}

func (db *TickerDB) SaveTicker(ticker *models.Ticker) (err error) {
	result := db.db.Model(ticker).Updates(ticker)
	err = result.Error
	if err == nil && result.RowsAffected == 0 {
		err = db.db.Create(ticker).Error
	}
	return
}
