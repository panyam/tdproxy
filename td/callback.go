package td

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type CallbackHandler struct {
	Port     int
	TDClient *Client
	CertFile string `default:"./td/server.crt"`
	PKeyFile string `default:"./td/server.key"`
}

func NewCallbackHandler(TDClient *Client, port int, cert_file string, pkey_file string) *CallbackHandler {
	return &CallbackHandler{TDClient: TDClient, Port: port, CertFile: cert_file, PKeyFile: pkey_file}
}

func (c *CallbackHandler) Start() error {
	handler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		// log.Printf("URL: %s", req.URL)
		// log.Printf("Query: %s", req.URL.Query())
		code := req.URL.Query().Get("code")
		log.Printf("Code: %s", code)
		resp, err := c.TDClient.CompleteAuth(code)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err)
		} else {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(resp)
		}
	})
	http.Handle("/callback", handler)
	http.Handle("/callback/", handler)
	log.Printf("Callback Handler Certificate file: %s\n", c.CertFile)
	log.Printf("Callback Handler Private key file: %s\n", c.PKeyFile)
	log.Printf("Running callback handler on part %d", c.Port)
	return http.ListenAndServeTLS(fmt.Sprintf(":%d", c.Port), c.CertFile, c.PKeyFile, nil)
}
