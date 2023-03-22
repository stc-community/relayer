package main

import (
	"encoding/json"
	"github.com/fiatjaf/relayer"
	"github.com/fiatjaf/relayer/storage/elasticsearch"
	"github.com/fiatjaf/relayer/storage/postgresql"
	"github.com/fiatjaf/relayer/storage/sqlite3"
	"github.com/kelseyhightower/envconfig"
	_ "github.com/lib/pq"
	"github.com/nbd-wtf/go-nostr"
	"log"
	"net/http"
)

const (
	Elasticsearch = "elasticsearch"
	Postgresql    = "postgresql"
	Sqlite3       = "sqlite3"
)

type Relay struct {
	StorageType     string `envconfig:"STORAGE_TYPE" required:"true" default:"sqlite3"` // elasticsearch, postgresql, sqlite3
	StorageDB       string `envconfig:"STORAGE_DB" default:"./storage.db"`
	CLNNodeID       string `envconfig:"CLN_NODE_ID"`
	CLNHost         string `envconfig:"CLN_HOST"`
	CLNRune         string `envconfig:"CLN_RUNE"`
	TicketPriceSats int64  `envconfig:"TICKET_PRICE_SATS"`
	storage         relayer.Storage
}

var r = &Relay{}

func (r *Relay) Name() string {
	return "ExpensiveRelay"
}

func (r *Relay) Storage() relayer.Storage {
	return r.storage
}

func (r *Relay) Init() error {
	// every hour, delete all very old events
	go r.storage.Clean()
	return nil
}

func (r *Relay) OnInitialized(s *relayer.Server) {
	// special handlers
	s.Router().Path("/").HandlerFunc(handleWebpage)
	s.Router().Path("/invoice").HandlerFunc(func(w http.ResponseWriter, rq *http.Request) {
		handleInvoice(w, rq, r)
	})
}

func (r *Relay) AcceptEvent(evt *nostr.Event) bool {
	// only accept they have a good preimage for a paid invoice for their public key
	if !checkInvoicePaidOk(evt.PubKey) {
		return false
	}
	// block events that are too large
	jsonb, _ := json.Marshal(evt)
	return len(jsonb) <= 100000
}

func main() {
	r := Relay{}
	if err := envconfig.Process("", &r); err != nil {
		log.Fatalf("failed to read from env: %v", err)
		return
	}
	r.storage = driver(r.StorageType, r.StorageDB)
	if err := relayer.Start(&r); err != nil {
		log.Fatalf("server terminated: %v", err)
	}
}

func driver(storageType, storageDB string) relayer.Storage {
	switch storageType {
	case Elasticsearch:
		return &elasticsearch.Elasticsearch{}
	case Postgresql:
		return &postgresql.Postgres{DatabaseURL: storageDB}
	case Sqlite3:
		return &sqlite3.SQLite3{DatabaseURL: storageDB}
	default:
		return &sqlite3.SQLite3{DatabaseURL: storageDB}
	}
}
