package tdclient

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type CallbackHandler struct {
	Port      int
	TDClient  *Client
	AuthStore *AuthStore
	CertFile  string `default:"./tdclient/server.crt"`
	PKeyFile  string `default:"./tdclient/server.key"`
}

func NewCallbackHandler(TDClient *Client, astore *AuthStore, port int, cert_file string, pkey_file string) *CallbackHandler {
	return &CallbackHandler{TDClient: TDClient, AuthStore: astore, Port: port, CertFile: cert_file, PKeyFile: pkey_file}
}

func (c *CallbackHandler) Start() (err error) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		log.Printf("URL: %s", req.URL)
		log.Printf("Query: %s", req.URL.Query())
		code := req.URL.Query().Get("code")
		log.Printf("Code: %s", code)
		err := c.TDClient.Auth.CompleteAuth(code)
		if err != nil {
			log.Println("CompleteAuthError: ", err)
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			c.AuthStore.SaveAuth(c.TDClient.Auth)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(c.TDClient.Auth.ToJson())
		}
	})
	http.Handle("/callback", handler)
	http.Handle("/callback/", handler)
	log.Printf("Running Callback Handler on Port: %d, Certificate file: %s, PKey File: %s\n", c.Port, c.CertFile, c.PKeyFile)
	if err = http.ListenAndServeTLS(fmt.Sprintf(":%d", c.Port), c.CertFile, c.PKeyFile, nil); err != nil {
		log.Fatal("Cannot start HTTPS callback handler: ", err)
	}
	return err
}
