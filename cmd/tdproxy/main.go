package main

import (
	"flag"
	"fmt"
	badger "github.com/dgraph-io/badger/v3"
	"github.com/panyam/goutils/utils"
	"github.com/panyam/pslite/cli"
	pslconfig "github.com/panyam/pslite/config"
	"google.golang.org/grpc"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"net"
	"path"
	"strconv"
	"strings"
	"tdproxy/config"
	"tdproxy/db"
	"tdproxy/db/gormdb"
	"tdproxy/protos"
	svc "tdproxy/services"
	"tdproxy/tdclient"
)

const TEST_CALLBACK_URL = "https://localhost:8000/callback"
const TEST_CLIENT_ID = ""

var (
	port              = flag.Int("port", config.DefaultServerPort(), "Port on which gRPC server should listen TCP conn.")
	tdroot            = flag.String("tdroot", "~/.tdroot", "Root location of where TD data is downloaded too")
	client_id         = flag.String("client_id", TEST_CLIENT_ID, "TD Ameritrade Client ID")
	callback_port     = flag.Int("callback_port", config.DefaultCallbackPort(), "Port on which OAuth Callback handler listen on.")
	callback_url      = flag.String("callback_url", TEST_CALLBACK_URL, "TD Ameritrade Auth Callback URl")
	callback_cert     = flag.String("callback_cert", "./tdclient/server.crt", "Certificate file for SSL Callback handler")
	callback_pkey     = flag.String("callback_pkey", "./tdclient/server.key", "Private key file for SSL Callback handler")
	topic_endpoint    = flag.String("topic_endpoint", fmt.Sprintf("%d", pslconfig.DefaultServerPort()), "End point where topics can be published and subscribed to")
	topics_folder     = flag.String("topics_folder", "~/.tdroot/topics", "End point where topics can be published and subscribed to")
	db_endpoint       = flag.String("db_endpoint", "postgres://postgres:docker@localhost:5432/tdproxydb", "Endpoint of DB backing tdproxy shard targets.  Supported - sqlite eg (sqlite://~/.tdproxy/sqlite.db) or postgres eg (postgres://user:pass@localhost:5432/dbname)")
	trade_db_endpoint = flag.String("tradesdb_endpoint", "~/.tdroot/trades", "Endpoint of trades DB")
)

func createPubsubClient() *cli.PubSub {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	pubsub, err := cli.NewPubSub(*topic_endpoint)
	if err != nil {
		log.Fatal(err)
	}
	if strings.Index(*topic_endpoint, ":") < 0 {
		// Start us locally on the given port
		topicport, err := strconv.ParseInt(*topic_endpoint, 10, 64)
		if err != nil {
			log.Fatal(fmt.Errorf("Invalid topic_endpoint port: %s", *topic_endpoint))
		}
		log.Printf("Serving pubsub topics on %d - %s", topicport, *topics_folder)
		go cli.PSLServe(int(topicport), *topics_folder)
	}
	if err != nil {
		log.Fatal(err)
	}
	return pubsub
}

func getTradeDB() *db.TradeDB {
	tdbroot := path.Join(utils.ExpandUserPath(*tdroot), "tradedb")
	indexfile := path.Join(tdbroot, "index.db")
	tradesdir := path.Join(tdbroot, "trades")
	indexdb, err := gorm.Open(sqlite.Open(indexfile), &gorm.Config{})
	if err != nil {
		log.Println("Cannot open/create Trade DB.  Is the tradedb dir missing?")
		log.Panic(err)
	}
	opt := badger.DefaultOptions(tradesdir)
	tradedb, err := badger.Open(opt)
	if err != nil {
		log.Fatal("Could not open tradedb: ", err)
	}
	return db.NewTradeDB(tradedb, indexdb)
}

func main() {
	flag.Parse()
	authdb, tickerdb, chaindb := getDBs()
	tradedb := getTradeDB()
	// see if we need to start the pubsub endpoint locally
	auth_store := tdclient.NewAuthStore(authdb)
	if *client_id != "" && *callback_url != "" {
		auth_store.EnsureAuth(*client_id, *callback_url)
	}
	log.Println("Auth Last: ", auth_store.LastAuth())
	tdinfo := tdclient.NewClient(utils.ExpandUserPath(*tdroot), chaindb, tickerdb)
	tdinfo.Auth = auth_store.LastAuth()
	callbackHandler := tdclient.NewCallbackHandler(tdinfo,
		auth_store,
		*callback_port,
		*callback_cert,
		*callback_pkey)
	go callbackHandler.Start()

	grpcServer := grpc.NewServer()
	protos.RegisterTickerServiceServer(grpcServer, &svc.TickerService{TDClient: tdinfo, AuthStore: auth_store})
	protos.RegisterChainServiceServer(grpcServer, &svc.ChainService{TDClient: tdinfo, AuthStore: auth_store})
	protos.RegisterTradeServiceServer(grpcServer, &svc.TradeService{TradeDB: tradedb})

	auth_svc := &svc.AuthService{TDClient: tdinfo, AuthStore: auth_store}
	protos.RegisterAuthServiceServer(grpcServer, auth_svc)

	pubsub := createPubsubClient()
	streamer_svc := svc.NewStreamerService(tdinfo, pubsub)
	streamer_svc.TopicsFolder = *topics_folder
	protos.RegisterStreamerServiceServer(grpcServer, streamer_svc)
	log.Printf("Initializing TDProxy gRPC server on port %d", *port)
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer.Serve(lis)
	defer tradedb.Close()
}

func OpenDB(db_endpoint string) (db *gorm.DB, err error) {
	var dbpath string
	if strings.HasPrefix(db_endpoint, "sqlite://") {
		dbpath = utils.ExpandUserPath((db_endpoint)[len("sqlite://"):])
		db, err = gorm.Open(sqlite.Open(dbpath), &gorm.Config{})
	} else if strings.HasPrefix(db_endpoint, "postgres://") {
		db, err = gorm.Open(postgres.Open(db_endpoint), &gorm.Config{})
	}
	if err != nil {
		log.Printf("Cannot connect DB: %s", db_endpoint)
	}
	return
}

func getDBs() (authdb db.AuthDB, tickerdb db.TickerDB, chaindb db.ChainDB) {
	gdb, err := OpenDB(*db_endpoint)
	if err != nil {
		panic(err)
	}
	authdb = gormdb.NewAuthDB(gdb)
	chaindb = gormdb.NewChainDB(gdb)
	tickerdb = gormdb.NewTickerDB(gdb)
	return authdb, tickerdb, chaindb
}
