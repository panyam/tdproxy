package filedb

import (
	"encoding/json"
	"fmt"
	"github.com/panyam/goutils/utils"
	"log"
	"os"
	"path"
	"tdproxy/models"
)

type AuthDB struct {
	RootDir string
	auths   map[string]*models.Auth
}

func NewAuthDB(rootdir string) *AuthDB {
	out := &AuthDB{RootDir: rootdir}
	fmt.Println("Client Root Dir: ", rootdir)
	os.MkdirAll(rootdir, 0777)
	if err := out.Reload(); err != nil {
		log.Panic("Error creating auth db: ", err)
	}
	return out
}

func (a *AuthDB) TokensFilePath() string {
	return path.Join(a.RootDir, "tokens")
}

/**
 * Persistes auth tokens to file so it can be used later on.
 */
func (a *AuthDB) SaveAuth(auth *models.Auth) (err error) {
	log.Println("Saving tokens...")
	defer log.Println("Finished Saved Tokens, err: ", err)
	auths := make(utils.StringMap)
	for key, value := range a.auths {
		auths[key] = value.ToJson()
	}
	var marshalled []byte
	marshalled, err = json.MarshalIndent(auths, "", "  ")
	if err != nil {
		log.Printf("Could not marshall token: %+v", a.auths)
		return err
	}
	err = os.WriteFile(a.TokensFilePath(), marshalled, 0777)
	return err
}

/**
 * Reloads the auth store contents.
 */
func (a *AuthDB) Reload() (err error) {
	a.auths = make(map[string]*models.Auth)
	contents, err := os.ReadFile(a.TokensFilePath())
	if err != nil {
		log.Println(err)
		return err
	}
	tokens, err := utils.JsonDecodeBytes(contents)
	if err != nil {
		log.Println(err)
		return err
	}
	entries := tokens.(utils.StringMap)
	for clientId, entry := range entries {
		clientInfo := entry.(utils.StringMap)
		callback_url := clientInfo["callback_url"]
		auth, err := a.EnsureAuth(clientId)
		if err != nil {
			return err
		}
		auth.CallbackUrl = callback_url.(string)
		auth.FromJson(clientInfo)
	}
	fmt.Println("Loaded auth tokens: ", entries)
	return nil
}

/**
 * Creates a new auth object and adds to the store.
 */
func (a *AuthDB) EnsureAuth(client_id string) (auth *models.Auth, err error) {
	var ok bool
	auth, ok = a.auths[client_id]
	if !ok || auth == nil {
		auth = &models.Auth{ClientId: client_id}
		if a.auths == nil {
			a.auths = make(map[string]*models.Auth)
		}
		a.auths[client_id] = auth
	}
	auth.ClientId = client_id
	return
}
