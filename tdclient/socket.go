package tdclient

import (
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/panyam/goutils/utils"
	"log"
	"net/url"
	"sync"
	"time"
)

type Socket struct {
	TDClient   *Client
	wsConn     *websocket.Conn
	waitGroup  sync.WaitGroup
	is_running bool

	// Time allowed to write a message to the peer.
	writeWaitTime time.Duration

	// Time allowed to read the next pong message from the peer.
	readWaitTime time.Duration

	// Keeps track of request IDs
	currRequestId int64
	requestMap    map[int64]interface{}

	// Requests channel is where we a sent a request that expects a response
	requestsChannel chan utils.StringMap
	// Channel for general control messages to the server
	controlChannel chan utils.StringMap

	readerChannel chan utils.StringMap
}

/**
 * Create a new connection object to the TD streaming server.
 */
func NewSocket(TDClient *Client, readerChannel chan utils.StringMap) *Socket {
	if readerChannel == nil {
		readerChannel = make(chan utils.StringMap, 10)
	}
	return &Socket{
		TDClient:        TDClient,
		readerChannel:   readerChannel,
		writeWaitTime:   10 * time.Second,
		readWaitTime:    60 * time.Second,
		requestsChannel: make(chan utils.StringMap),
		controlChannel:  make(chan utils.StringMap),
		requestMap:      make(map[int64]interface{}),
	}
}

/**
 * Returns the socket's reader channel.
 */
func (s *Socket) ReaderChannel() chan utils.StringMap {
	return s.readerChannel
}

/**
 * Returns whether the socket reader/writer loops are running.
 */
func (s *Socket) IsRunning() bool {
	return s.is_running
}

/**
 * Sends a new request to the peer.  If the peer is not connected or
 * writers are not running, false is returned.
 */
func (s *Socket) SendRequest(req utils.StringMap) bool {
	if !s.is_running {
		log.Println("Server is not running")
		return false
	}
	s.requestsChannel <- req
	return true
}

func (s *Socket) UserPrincipals() (utils.StringMap, error) {
	if s.TDClient.Auth == nil {
		return nil, NotAuthenticated
	}
	return s.TDClient.Auth.EnsureUserPrincipals()
}

/**
 * Gets the credentials required to connect to the streaming server.
 * This method also checks that a valid login/auth is available and
 * will try to refresh any tokens if required before failing on errors.
 */
func (s *Socket) Credentials() (utils.StringMap, error) {
	if s.TDClient.Auth == nil {
		return nil, NotAuthenticated
	}
	return s.TDClient.Auth.StreamingCredentials()
}

/**
 * Returns the WS connection URL.
 */
func (s *Socket) WSUrl() (*url.URL, error) {
	if s.TDClient.Auth == nil {
		return nil, NotAuthenticated
	}
	return s.TDClient.Auth.WSUrl()
}

/**
 * Helper to create new requests.
 */
func (s *Socket) NewRequest(service string,
	command string,
	record bool,
	onParams func(params utils.StringMap)) (utils.StringMap, error) {

	userPrincipals, err := s.UserPrincipals()
	if err != nil {
		return nil, err
	}
	accounts := userPrincipals["accounts"].([]interface{})
	account := accounts[0].(utils.StringMap)
	streamerInfo := userPrincipals["streamerInfo"].(utils.StringMap)
	payload := fmt.Sprintf(`{"requests": [{
							"service": "%s",
							"command": "%s",
							"requestid": "%d",
							"account": "%s",
							"source": "%s",
							"parameters": {}
					}]}
	`, service, command, s.currRequestId,
		account["accountId"],
		streamerInfo["appId"],
	)
	val, _ := utils.JsonDecodeStr(payload)
	if record {
		s.requestMap[s.currRequestId] = val
	}
	s.currRequestId += 1
	out := val.(utils.StringMap)
	if onParams != nil {
		requests := out["requests"].([]interface{})
		req := requests[0].(utils.StringMap)
		reqparams := req["parameters"].(utils.StringMap)
		onParams(reqparams)
	}
	return out, nil
}

/**
 * This method is called to connect to the socket server.
 * If already connected then nothing is done and nil
 * is not already connected, a connection will first be established
 * (including auth and refreshing tokens) and then the reader and
 * writers are started.   SendRequest can be called to send requests
 * to the peer and the (user provided) readerChannel will be used to
 * handle messages from the server.
 */
func (s *Socket) Disconnect() error {
	if !s.IsRunning() {
		// already running do nothing
		return nil
	}
	s.controlChannel <- nil
	return nil
}

/**
 * Starts the connection to the server allowing sending and receiving of messages
 * to/from the server.
 *
 * This method will start the reader and writer go-routines and return immediately
 * if no errors are encountered.
 *
 * It is the user's responsibility to call the WaitForFinish method to ensure
 * no premature exits.
 */
func (s *Socket) StartConnection() error {
	if s.IsRunning() {
		// already running do nothing
		return nil
	}
	wsurl, err := s.WSUrl()
	if err != nil {
		return err
	}
	c, _, err := websocket.DefaultDialer.Dial(wsurl.String(), nil)
	if err != nil {
		log.Println("Error dialling: ", err, wsurl)
		return err
	}
	s.wsConn = c
	s.is_running = true
	s.waitGroup.Add(2)
	// Start the writer and readers
	go s.start_reader()
	go s.start_writer()
	return nil
}

/**
 * Waits until the socket connection is disconnected or manually stopped.
 */
func (s *Socket) WaitForFinish() {
	s.waitGroup.Wait()
	s.is_running = false
	s.wsConn = nil
}

func (s *Socket) start_reader() error {
	// Start reader goroutine
	// defer c.Close()
	s.wsConn.SetReadDeadline(time.Now().Add(s.readWaitTime))
	s.wsConn.SetPongHandler(func(string) error {
		s.wsConn.SetReadDeadline(time.Now().Add(s.readWaitTime))
		return nil
	})
	defer func() {
		if s.wsConn != nil {
			s.wsConn.Close()
		}
		s.waitGroup.Done()
	}()

	cmd, err := s.LoginCommand()
	if err != nil {
		log.Println("Cannot create login command: ", err)
		return err
	}
	s.SendRequest(cmd)
	log.Println("Starting reader: ")
	for {
		var newMessage interface{}
		err := s.wsConn.ReadJSON(&newMessage)
		if err != nil {
			log.Println("Error reading message: ", err)
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("Unexpected Close Error: %v", err)
			}
			// Closing, Send error too?
			s.readerChannel <- nil
			break
		} else {
			s.readerChannel <- newMessage.(utils.StringMap)
		}
	}
	return nil
}

// Start writer goroutine
func (s *Socket) start_writer() error {
	ticker := time.NewTicker((s.readWaitTime * 9) / 10)
	defer func() {
		ticker.Stop()
		if s.wsConn != nil {
			s.wsConn.Close()
		}
		s.waitGroup.Done()
	}()

	log.Println("Starting writer: ")
	for {
		select {
		case newRequest := <-s.requestsChannel:
			s.wsConn.SetWriteDeadline(time.Now().Add(s.writeWaitTime))
			// Here we send a request to the server
			err := s.wsConn.WriteJSON(newRequest)
			if err != nil {
				log.Println("Error sending request: ", newRequest, err)
				return err
			}
			log.Println("Successfully Sent Request: ", newRequest)
			break
		case controlRequest := <-s.controlChannel:
			// For now only a "kill" can be sent here
			log.Println("Received kill signal.  Quitting Reader.", controlRequest)
			newReq, err := s.NewRequest("ADMIN", "LOGOUT", false, nil)
			if err == nil {
				err := s.wsConn.WriteJSON(newReq)
				if err == nil {
					err = s.wsConn.WriteMessage(websocket.CloseMessage, []byte{})
				}
			}
			return err
		case <-ticker.C:
			s.wsConn.SetWriteDeadline(time.Now().Add(s.writeWaitTime))
			if err := s.wsConn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return err
			}
		}
	}
}

/**
 * Get the login command payload.
 */
func (s *Socket) LoginCommand() (utils.StringMap, error) {
	userPrincipals, err := s.TDClient.Auth.EnsureUserPrincipals()
	if err != nil {
		return nil, err
	}
	streamerInfo := userPrincipals["streamerInfo"].(utils.StringMap)
	return s.NewRequest("ADMIN", "LOGIN", true, func(params utils.StringMap) {
		creds, _ := s.Credentials()
		params["credential"] = utils.JsonToQueryString(creds)
		params["token"] = streamerInfo["token"]
		params["version"] = "1.0"
	})
}
