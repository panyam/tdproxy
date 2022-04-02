package gormdb

import (
	"log"
	// "gorm.io/driver/sqlite"
	"errors"
	"gorm.io/gorm"
	// "gorm.io/gorm/clause"
	"tdproxy/models"
)

type AuthDB struct {
	db *gorm.DB
}

func NewAuthDB(db *gorm.DB) *AuthDB {
	db.AutoMigrate(&models.Auth{})
	return &AuthDB{
		db: db,
	}
}

func (a *AuthDB) LastAuth() *models.Auth {
	var out models.Auth
	err := a.db.First(&out).Error
	if err != nil {
		return nil
	}
	return &out
}

func (authdb *AuthDB) EnsureAuth(client_id string) (auth *models.Auth, err error) {
	auth, err = authdb.GetAuth(client_id)
	if err != nil {
		log.Println("Err Ensuring Auth: ", err, auth)
	}
	if err == nil && auth == nil {
		// Does not exist so create
		auth = &models.Auth{ClientId: client_id}
		err = authdb.SaveAuth(auth)
	}
	return
}

func (db *AuthDB) GetAuth(client_id string) (*models.Auth, error) {
	var out models.Auth
	err := db.db.First(&out, "client_id = ?", client_id).Error
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
	result := authdb.db.Session(&gorm.Session{FullSaveAssociations: true}).Updates(auth)
	err = result.Error
	if err == nil && result.RowsAffected == 0 {
		result = authdb.db.Create(auth)
	}
	err = result.Error
	return
}
