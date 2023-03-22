package main

import (
	"encoding/json"
	"fmt"
	"github.com/fiatjaf/relayer"
	"github.com/fiatjaf/relayer/storage/elasticsearch"
	"github.com/fiatjaf/relayer/storage/postgresql"
	"github.com/fiatjaf/relayer/storage/sqlite3"
	"github.com/kelseyhightower/envconfig"
	"github.com/nbd-wtf/go-nostr"
	"log"
)

const (
	Elasticsearch = "elasticsearch"
	Postgresql    = "postgresql"
	Sqlite3       = "sqlite3"
)

type Relay struct {
	StorageType string   `envconfig:"STORAGE_TYPE" required:"true" default:"sqlite3"` // elasticsearch, postgresql, sqlite3
	StorageDB   string   `envconfig:"STORAGE_DB" default:"./storage.db"`
	Whitelist   []string `envconfig:"WHITELIST"`
	storage     relayer.Storage
}

func (r *Relay) Name() string {
	return "BasicRelay"
}

func (r *Relay) Storage() relayer.Storage {
	return r.storage
}

func (r *Relay) OnInitialized(*relayer.Server) {}

func (r *Relay) Init() error {
	err := envconfig.Process("", r)
	if err != nil {
		return fmt.Errorf("couldn't process envconfig: %w", err)
	}
	go r.Storage().Clean()
	return nil
}

func (r *Relay) AcceptEvent(evt *nostr.Event) bool {
	// disallow anything from non-authorized pubkeys
	for _, pubkey := range r.Whitelist {
		if pubkey == evt.PubKey {
			jsonb, _ := json.Marshal(evt)
			return len(jsonb) <= 100000
		}
	}
	return false
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
