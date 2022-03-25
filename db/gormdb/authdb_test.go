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
	// "tdproxy/models"
	"testing"
	// "time"
)

func CreateTestAuthDB(t *testing.T, filepath string) (*AuthDB, string) {
	dir := ""
	if filepath == "" {
		dir, err := ioutil.TempDir("/tmp", "tdproxydb")
		if err != nil {
			log.Fatal(err)
		}

		filepath = path.Join(dir, "test.db")
	}
	db, err := gorm.Open(sqlite.Open(filepath), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	fmt.Println("DBRoot: ", filepath)
	return NewAuthDB(db), dir
}

func TestNewAuthDB(t *testing.T) {
	_, dbroot := CreateTestAuthDB(t, "")
	defer os.RemoveAll(dbroot)
}

func TestAuthEnsureAuth(t *testing.T) {
	db, dbroot := CreateTestAuthDB(t, "")
	defer os.RemoveAll(dbroot)

	auth, err := db.EnsureAuth("testclient1")
	assert.Equal(t, err, nil, "Error should be nil")
	assert.Equal(t, auth.ClientId, "testclient1", "ClientId Mismatch")
}

func TestSaveAuth(t *testing.T) {
	db, dbroot := CreateTestAuthDB(t, "/tmp/sq.db")
	if dbroot != "" {
		defer os.RemoveAll(dbroot)
	}

	auth, err := db.EnsureAuth("testclient1")
	assert.Equal(t, err, nil, "Error should be nil")
	assert.Equal(t, auth.ClientId, "testclient1", "ClientId Mismatch")
	assert.Equal(t, auth.AuthToken.HasValue(), false, "")
	assert.Equal(t, auth.UserPrincipals.HasValue(), false, "")

	auth.CallbackUrl = "http://hello.world.com"
	a1 := map[string]interface{}{
		"a":                        float64(11),
		"b":                        "hello",
		"c":                        true,
		"expires_in":               60.0,
		"refresh_token_expires_in": 120.0,
	}
	auth.SetAuthToken(a1)

	b1 := map[string]interface{}{
		"x": float64(42),
		"y": "world",
		"z": false,
	}
	auth.SetUserPrincipals(b1)

	err = db.SaveAuth(auth)
	assert.Equal(t, err, nil, "Save should succeed")

	fetched, err := db.GetAuth("testclient1")
	assert.Equal(t, err, nil, "GetAuth should succeed")
	assert.NotEqual(t, fetched, nil, "GetAuth should succeed")
	auth.AuthToken.LastUpdatedAt = fetched.AuthToken.LastUpdatedAt
	auth.UserPrincipals.LastUpdatedAt = fetched.UserPrincipals.LastUpdatedAt
	log.Println("Saved: ", auth)
	log.Println("Fetched: ", fetched)
	assert.Equal(t, auth, fetched, "Saved and Fetched auth should be equal")
}
